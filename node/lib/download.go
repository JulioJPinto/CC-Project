package lib

import (
	"cc_project/helpers"
	"cc_project/protocol"
	"cc_project/protocol/fstp"
	"cc_project/protocol/p2p"
	"encoding/json"
	"path"
	"sort"

	"github.com/fatih/color"

	"cc_project/helpers/sync"
	"fmt"
	"os"
	"time"
)

type Status int

const (
	Missing    = -1
	Downloaded = 0
	Pending    = 1 // pending for n iterations
)

type AddrMessage struct {
	msg  p2p.Message
	peer protocol.DeviceIdentifier
}

type Downloader struct {
	node      *Node
	file      protocol.FileHash
	channel   chan AddrMessage
	segments  *sync.Map[int, Status]
	done      *sync.Flag
	interrupt *sync.Flag
	whoHas    map[protocol.DeviceIdentifier][]protocol.FileSegment
}

func NewDownloader(node *Node, file protocol.FileHash) *Downloader {
	channel := make(chan AddrMessage)
	done := sync.Flag{}
	interrupt := sync.Flag{}
	interrupt.Unset()
	done.Unset()
	return &Downloader{
		node:      node,
		file:      file,
		channel:   channel,
		segments:  &sync.Map[int, Status]{},
		done:      &done,
		interrupt: &interrupt,
		whoHas:    nil,
	}
}

func (d *Downloader) stillNeed() map[protocol.DeviceIdentifier][]protocol.FileSegment {
	ret := map[protocol.DeviceIdentifier][]protocol.FileSegment{}
	for peer, segments := range d.whoHas {
		needed := []protocol.FileSegment{}
		for _, segment := range segments {
			if x, _ := d.segments.Load(int(segment.BlockOffset)); x == Missing {
				needed = append(needed, segment)
			}
		}
		if len(needed) > 0 {
			ret[peer] = needed
		}
	}
	return ret
}

func (d *Downloader) String() string {
	ret := ""
	ret += "Downloader State:\n"
	ret += fmt.Sprintf(" - File: %d\n", d.file)
	ret += fmt.Sprintf(" - Channel: %p\n", d.channel)
	ret += fmt.Sprintf(" - Done: %v\n", d.done.IsSet())

	fmt.Println(" - Segments:")
	d.segments.Range(func(key int, value Status) bool {
		segmentIndex := key

		status := value
		switch status {
		case Missing:
			{
				ret += fmt.Sprintf("   - Segment %d: Missing\n", segmentIndex)

			}
		case Downloaded:
			{
				ret += fmt.Sprintf("   - Segment %d: Downloaded\n", segmentIndex)
			}
		default:
			{
				ret += fmt.Sprintf("   - Segment %d: %d\n", segmentIndex, status)

			}
		}
		return true
	})

	ret += (" - WhoHas:")
	for deviceID, segments := range d.whoHas {
		ret += fmt.Sprintf("   - Device ID: %s\n", deviceID)
		ret += ("     - Segments:\n")
		for _, segment := range segments {
			ret += fmt.Sprintf("       - BlockOffset: %d, SegmentHash: %d\n",
				segment.BlockOffset, segment.Hash)
		}
	}
	return ret
}

func (d *Downloader) ForwardMessage(msg p2p.Message, peer protocol.DeviceIdentifier) {
	if !d.done.IsSet() {
		d.channel <- AddrMessage{msg: msg, peer: peer}
	}
}

func (d *Downloader) Start() error {
	d.node.Downloads.Store(d.file, d)
	file_meta_data, ok := d.node.KnownFiles.Load(d.file) // file_meta_data,ok := c.KnownFiles.get(d.file)

	if !ok {
		d.node.FetchFiles(nil)
		file_meta_data, ok = d.node.KnownFiles.Load(d.file)
		if !ok {
			return fmt.Errorf("files does not exist: %v", d.file)
		}
	}

	resp, err := d.node.FSTPclient.Request(fstp.NewWhoHasRequest(fstp.WhoHasReqProps{File: d.file}))
	if err != nil {
		return err
	}
	pay, ok := resp.Payload.(*fstp.WhoHasRespProps)
	x, _ := json.Marshal(pay)
	println(string(x))
	if !ok {
		return fmt.Errorf("invalid payload")
	}

	p := map[protocol.DeviceIdentifier][]protocol.FileSegment(*pay)
	d.whoHas = p
	if d.node.Debug {
		for k, v := range p {
			color.Cyan(string(k) + ": ")
			for s := range v {
				x, _ := json.Marshal(s)
				color.Cyan("\t" + string(x))
			}
		}
	}

	if resp.Header.Flags == fstp.ErrResp {
		err_resp := resp.Payload.(*fstp.ErrorResponse)
		return fmt.Errorf(err_resp.Err)
	}

	for n := range file_meta_data.SegmentHashes {
		d.segments.Store(n, Missing)
	}

	path_ := d.node.NodeDir
	go d.update_who_has(500)
	go d.await_segment_responses(file_meta_data, path_)
	d.send_segment_requests()

	go d.node.makeFileAvailable(path.Join(path_, file_meta_data.Name))
	return nil
}

func (d *Downloader) Abort() error {
	d.interrupt.Set()
	return nil
}

func (d *Downloader) update_who_has(sleep_for int) {
	for {
		if d.done.IsSet() {
			break
		}
		time.Sleep(time.Millisecond * time.Duration(sleep_for))
		resp, err := d.node.FSTPclient.Request(fstp.NewWhoHasRequest(fstp.WhoHasReqProps{File: d.file}))
		if err != nil {
			continue
		}
		pay, ok := resp.Payload.(*fstp.WhoHasRespProps)
		if !ok {
			continue
		}

		p := map[protocol.DeviceIdentifier][]protocol.FileSegment(*pay)
		d.whoHas = p
	}
}

func getSortedKeys(m map[protocol.DeviceIdentifier]float64) []protocol.DeviceIdentifier {
	// Create a slice to store the keys
	keys := make([]protocol.DeviceIdentifier, 0, len(m))

	// Populate the slice with keys from the map
	for key := range m {
		keys = append(keys, key)
	}

	// Sort the keys based on their corresponding values
	sort.Slice(keys, func(i, j int) bool {
		return m[keys[i]] < m[keys[j]]
	})

	return keys
}

func (d *Downloader) pickGaijos(set *helpers.Set[protocol.DeviceIdentifier], max int) map[protocol.DeviceIdentifier]float64 {
	ret_ := map[protocol.DeviceIdentifier]float64{}
	gajos := set.Slice()
	max_weight := 0.
	for _, gajo := range gajos {
		stats_gajo, ok := d.node.PeerStats.Load(gajo)
		if ok {
			weight := stats_gajo.Weight()
			if weight > max_weight {
				max_weight = weight
			}
			ret_[gajo] = weight
		} else {
			ret_[gajo] = 1

		}
	}
	sorted := getSortedKeys(ret_)
	ret := map[protocol.DeviceIdentifier]float64{}
	for i := 0; i < max && i < len(sorted); i++ {
		ret[sorted[i]] = max_weight / (ret_[sorted[i]] + 0.1)
	}
	return ret

}

func (d *Downloader) send_segment_requests() bool {
	for {
		if d.done.IsSet() || d.interrupt.IsSet() {
			close(d.channel)
			return true
		}

		n_gaijos := 4
		needed := d.stillNeed()
		set_gajos := helpers.MapKeys(needed)
		gajos := d.pickGaijos(set_gajos, n_gaijos)
		for gajo, peso := range gajos {
			needed_from_gajo := needed[gajo]
			for i := 0; i < int(peso)+1 && i < len(needed_from_gajo); i++ {
				segment := needed_from_gajo[i]
				status, _ := d.segments.Load(int(segment.BlockOffset))
				if status == Missing {
					d.node.RequestSegment(gajo, segment)
					d.segments.Store(int(segment.BlockOffset), Pending)
				}
			}
		}
		d.segments.Range(func(key int, status Status) bool {
			if status > 0 {
				status++
				d.segments.Store(key, status)
			}
			if status > 5 {
				d.segments.Store(key, Missing)
			}
			return true
		})
		time.Sleep(1000 * time.Millisecond)
	}
}

func (d *Downloader) checkIfDone() {
	done := true

	d.segments.Range(func(key int, status Status) bool {
		if status != Downloaded {
			done = false
			// Stop ranging since we found a segment that is not downloaded
			return false
		}

		return true
	})

	if done {
		d.done.Set()
	}
}

func (d *Downloader) await_segment_responses(file protocol.FileMetaData, path_ string) {
	fmt.Println(d.channel)

	fmt.Println(d.channel)

	store_path := path.Join(path_, file.Name)
	writef, err := os.Create(store_path)
	if err != nil {
		return
	}

	segments := []protocol.FileSegment{}

	for addrmsg := range d.channel {
		println("SEGMETNS:", segments)
		if len(segments) > 5 {
			println("\n\n\n\n\nAHOY CAPTAIN \n\n\n\n\n\n\n\n")
			d.node.FSTPclient.Request(fstp.IHaveSegmentsReq(segments))
			segments = []protocol.FileSegment{}
		}
		offset := addrmsg.msg.Header.SegmentOffset
		segmente_offset := offset * protocol.SegmentMaxLength
		// headerJSON, _ := json.MarshalIndent(addrmsg.msg.Header, "", "  ")
		// fmt.Println(string(headerJSON))

		segment_hash := protocol.HashSegment(addrmsg.msg.Payload, int(addrmsg.msg.Length))
		if file.SegmentHashes[addrmsg.msg.Header.SegmentOffset] == segment_hash {
			show1 := fmt.Sprintf("the hashin do be matchin in segment %d of file %s\n", addrmsg.msg.Header.SegmentOffset, file.Name)
			color.Green(show1)
			segments = append(segments, protocol.FileSegment{BlockOffset: int64(offset), FileHash: file.Hash, Hash: file.SegmentHashes[addrmsg.msg.Header.SegmentOffset]})
			d.segments.Store(int(addrmsg.msg.Header.SegmentOffset), Downloaded)
			d.checkIfDone()
			writef.Seek(int64(segmente_offset), 0)
			writef.Write([]byte(addrmsg.msg.Payload))
		} else {
			show1 := fmt.Sprintf("the hashin do NOT be matchin in segment %d of file %s\n", addrmsg.msg.Header.SegmentOffset, file.Name)
			color.Red(show1)
			show := fmt.Sprintf("%d vs %d", file.SegmentHashes[addrmsg.msg.Header.SegmentOffset], segment_hash)
			color.Red(show)
			peer_stats, _ := d.node.PeerStats.Load(addrmsg.peer)
			peer_stats.NDroppedPackets++
			d.node.PeerStats.Store(addrmsg.peer, peer_stats)
			d.checkIfDone()

		}
		if d.done.IsSet() {
			return
		}
	}
}
