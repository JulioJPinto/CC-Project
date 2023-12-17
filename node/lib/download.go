package lib

import (
	"cc_project/helpers"
	"cc_project/protocol"
	"cc_project/protocol/fstp"
	"cc_project/protocol/p2p"
	"encoding/hex"
	"encoding/json"
	"path"

	"github.com/fatih/color"

	"cc_project/helpers/sync"
	"fmt"
	"net"
	"os"
	"time"
)

type Status int

const (
	Missing    = -1
	Downloaded = 0
	Pending    = 1 // pending for n iterations
)

type Downloader struct {
	node     *Node
	file     protocol.FileHash
	channel  chan p2p.Message
	segments *sync.Map[int, Status]
	done     *sync.Flag
	whoHas   map[protocol.DeviceIdentifier][]protocol.FileSegment
}

func (d *Downloader) ForwardMessage(msg p2p.Message) {
	d.channel <- msg
}

func NewDownloader(node *Node, file protocol.FileHash) *Downloader {
	channel := make(chan p2p.Message)
	done := sync.Flag{}
	done.Unset()
	return &Downloader{
		node:     node,
		file:     file,
		channel:  channel,
		segments: &sync.Map[int, Status]{},
		done:     &done,
		whoHas:   nil,
	}
}
func (d *Downloader) Start() error {
	channel := make(chan p2p.Message)
	d.node.Downloads.Store(d.file, d)
	color.Green("created channel for " + fmt.Sprintf("%d", d.file))
	file_meta_data, ok := d.node.KnownFiles[d.file] // file_meta_data,ok := c.KnownFiles.get(d.file)

	if !ok {
		d.node.FetchFiles(nil)
		file_meta_data, ok = d.node.KnownFiles[d.file]
		if !ok {
			return fmt.Errorf("files does not exist: %v", d.file)
		}
	}

	resp, err := d.node.FSTPclient.Request(fstp.NewWhoHasRequest(fstp.WhoHasReqProps{File: d.file}))
	if err != nil {
		return err
	}
	pay, ok := resp.Payload.(*fstp.WhoHasRespProps)
	if !ok {
		return fmt.Errorf("invalid payload")
	}

	p := map[protocol.DeviceIdentifier][]protocol.FileSegment(*pay)
	d.whoHas = p

	for k, v := range p {
		color.Cyan(string(k) + ": ")
		for s := range v {
			x, _ := json.Marshal(s)
			color.Cyan("\t" + string(x))
		}
	}
	if resp.Header.Flags == fstp.ErrResp {
		err_resp := resp.Payload.(*fstp.ErrorResponse)
		return fmt.Errorf(err_resp.Err)
	}

	for n := range file_meta_data.SegmentHashes {
		d.segments.Store(n, Missing)
	}

	path := d.node.NodeDir

	go d.await_segment_responses(file_meta_data, path)
	d.send_segment_requests()
	close(channel)
	return nil
}
func (node *Node) DownloadFile(file_hash protocol.FileHash) error {
	color.Green("DOWNLOADIN " + fmt.Sprintf("%d", file_hash))

	if _, ok := node.Downloads.Load(file_hash); ok {
		return fmt.Errorf("download already in progress")
	}
	downloader := NewDownloader(node, file_hash)
	err := downloader.Start()
	node.Downloads.Delete(file_hash)
	return err
}

func (d *Downloader) send_segment_requests() {

	color.Cyan("requesting ...  ")

	for {
		if d.done.IsSet() {
			break
		}
		for id, segments := range d.whoHas {
			for _, segment := range segments {
				// var status Status = Missing
				status, _ := d.segments.Load(int(segment.BlockOffset))
				if status == Missing {
					d.node.RequestSegment(id, segment)
				}
			}
			time.Sleep(100 * time.Millisecond)
		}
	}
}
func (d *Downloader) checkIfDone() {
	done := d.segments.Fold(true, func(a any, k int, v Status) any {
		return a == true && v == Downloaded
	})
	if done == true {
		d.done.Set()
	}

}
func (d *Downloader) await_segment_responses(file protocol.FileMetaData, path_ string) {
	color.Cyan("awayting ...  " + path_)
	// defer node.Chanels.Delete(file.Hash)
	store_path := path.Join(path_, file.Name)
	writef, err := os.Create(store_path)
	if err != nil {
		return
	}
	
	for msg := range d.channel {
		// show := fmt.Sprint("\nrecieved: ", string(msg.Payload))
		// color.Cyan(show)
		segmente_offset := msg.Header.SegmentOffset * protocol.SegmentMaxLength
		headerJSON, _ := json.MarshalIndent(msg.Header, "", "  ")
		fmt.Println(string(headerJSON))
		fmt.Println(hex.EncodeToString(msg.Payload))

		if file.SegmentHashes[msg.Header.SegmentOffset] == protocol.HashSegment(msg.Payload, int(msg.Length)) {
			color.Cyan("the hashin do be matchin")
			d.segments.Store(int(msg.Header.SegmentOffset), Downloaded)
			go d.checkIfDone()
			writef.Seek(int64(segmente_offset), 0)
			writef.Write([]byte(msg.Payload))
		} else {
			color.Red("the hashin do NOT be matchin")
			show := fmt.Sprintf("%d vs %d", file.SegmentHashes[msg.Header.SegmentOffset], protocol.HashSegment(msg.Payload, int(msg.Length)))
			color.Red(show)
			os.Exit(1)
			d.segments.Store(int(msg.Header.SegmentOffset), Downloaded)
			go d.checkIfDone()
			writef.Seek(int64(segmente_offset), 0)
			writef.Write([]byte(msg.Payload))

		}
		// writef.Seek(int64(segmente_offset), 0)
	}
}

func (node *Node) Distribute(device_segments map[protocol.DeviceIdentifier][]protocol.FileSegment, metadata protocol.FileMetaData) map[int64]protocol.DeviceIdentifier {
	ret := make(map[int64]protocol.DeviceIdentifier)
	for n := int64(0); n < int64(metadata.Length/protocol.SegmentMaxLength); n++ {
		for device, segments := range device_segments {
			for _, segment := range segments {
				if segment.BlockOffset == n {
					ret[n] = device
				}
			}
		}
	}
	return ret
}

func (node *Node) RequestSegment(peer protocol.DeviceIdentifier, segment protocol.FileSegment) {
	timestamp := helpers.TrunkI64(time.Now().UnixMilli())
	req := p2p.Message(p2p.GimmeFileSegmentRequest(segment, uint32(timestamp)))
	x := &req
	b, err := x.Serialize()
	if err != nil {
		return
	}
	addr, _ := net.ResolveUDPAddr("udp", string(peer))
	addr.Port = 9090
	node.sender.Send(*addr, b)
}
