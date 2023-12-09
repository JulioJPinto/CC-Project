package lib

import (
	"cc_project/helpers"
	"cc_project/protocol"
	"cc_project/protocol/fstp"
	"cc_project/protocol/p2p"
	"encoding/json"
	"path"

	"github.com/fatih/color"

	"fmt"
	"net"
	"os"
	"time"
)

type Downloader struct {
	node    *Node
	file    protocol.FileHash
	channel chan p2p.Message
}

func (d *Downloader) ForwardMessage(msg p2p.Message) {
	d.channel <- msg
}

func NewDownloader(node *Node, file protocol.FileHash) *Downloader {
	channel := make(chan p2p.Message)
	return &Downloader{node: node, file: file, channel: channel}
}
func (d *Downloader) Start() error {
	channel := make(chan p2p.Message)
	d.node.Chanels.Store(d.file, channel)
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

	path := d.node.NodeDir

	go d.node.await_segment_responses(file_meta_data, path)
	d.node.send_segment_requests(p)

	return nil
}
func (node *Node) DownloadFile(file_hash protocol.FileHash) error {
	color.Green("DOWNLOADIN " + fmt.Sprintf("%d", file_hash))

	if _, ok := node.Chanels.Load(file_hash); ok {
		return fmt.Errorf("download already in progress")
	}
	downloader := NewDownloader(node, file_hash)
	return downloader.Start()
}

func (node *Node) send_segment_requests(m map[protocol.DeviceIdentifier][]protocol.FileSegment) {
	color.Cyan("requesting ...  ")

	for id, segments := range m {
		for _, segment := range segments {
			node.RequestSegment(id, segment)
		}
	}
}

func (node *Node) await_segment_responses(file protocol.FileMetaData, path_ string) {
	color.Cyan("awayting ...  " + path_)
	// defer node.Chanels.Delete(file.Hash)
	store_path := path.Join(path_, file.Name)
	writef, err := os.Create(store_path)
	if err != nil {
		return
	}
	ch_, _ := node.Chanels.Load(file.Hash)
	ch := ch_.(chan p2p.Message)

	for msg := range ch {
		// show := fmt.Sprint("\nrecieved: ", string(msg.Payload))
		// color.Cyan(show)
		segmente_offset := msg.Header.SegmentOffset * protocol.SegmentLength
		if file.SegmentHashes[msg.Header.SegmentOffset] == protocol.HashSegment(msg.Payload, len(msg.Payload)) {
			color.Cyan("the hashin do be matchin")
			writef.Seek(int64(segmente_offset), 0)
			writef.Write([]byte(msg.Payload))
		} else {
			color.Red("the hashin do NOT be matchin")
			show := fmt.Sprintf("%d vs %d", file.SegmentHashes[msg.Header.SegmentOffset], protocol.HashSegment(msg.Payload, len(msg.Payload)))
			color.Red(show)
			writef.Seek(int64(segmente_offset), 0)
			writef.Write([]byte(msg.Payload))

		}
		// writef.Seek(int64(segmente_offset), 0)
	}
}

func (node *Node) Distribute(device_segments map[protocol.DeviceIdentifier][]protocol.FileSegment, metadata protocol.FileMetaData) map[int64]protocol.DeviceIdentifier {
	ret := make(map[int64]protocol.DeviceIdentifier)
	for n := int64(0); n < int64(metadata.Length/protocol.SegmentLength); n++ {
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
