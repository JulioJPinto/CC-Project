package lib

import (
	"cc_project/helpers"
	"cc_project/protocol"
	"cc_project/protocol/fstp"
	"cc_project/protocol/p2p"
	"fmt"
	"os"
	"net"
	"time"
)

func (node *Node) DownloadFile(file_hash protocol.FileHash) error {
	if _, ok := node.Chanels.Get(file_hash); ok {
		return fmt.Errorf("download already in progress")
	}

	file_meta_data, ok := node.KnownFiles[file_hash] // file_meta_data,ok := c.KnownFiles.get(file_hash)

	if !ok {
		node.FetchFiles(nil)
		file_meta_data, ok = node.KnownFiles[file_hash]
		if !ok {
			return fmt.Errorf("files does not exist: %v", file_hash)
		}
	}

	resp, err := node.FSTPclient.Request(fstp.NewWhoHasRequest(fstp.WhoHasReqProps{File: file_hash}))
	if err != nil {
		return err
	}
	pay, ok := resp.Payload.(*fstp.WhoHasRespProps)
	if !ok {
		return fmt.Errorf("invalid payload")
	}
	p := map[protocol.DeviceIdentifier][]protocol.FileSegment(*pay)
	if resp.Header.Flags == fstp.ErrResp {
		err_resp := resp.Payload.(*fstp.ErrorResponse)
		return fmt.Errorf(err_resp.Err)
	}

	path := "../client_files/downloaded/" + "test"

	go node.send_segment_requests(p)
	go node.await_segment_responses(file_meta_data, path)

	return nil
}

func (node *Node) send_segment_requests(m map[protocol.DeviceIdentifier][]protocol.FileSegment) {
	for id, segments := range m {
		for _, segment := range segments {
			node.RequestSegment(id, segment)
		}
	}
}

func (node *Node) await_segment_responses(file protocol.FileMetaData, path string) {
	writef, err := os.Create(path)
	if err != nil {
		return
	}
	ch, _ := node.Chanels.Get(file.Hash)

	for msg := range ch {
		segmente_offset := msg.Header.SegmentOffset * protocol.SegmentLength
		if file.SegmentHashes[msg.Header.SegmentOffset] == protocol.HashSegment(msg.Payload, len(msg.Payload)) {
			writef.Seek(int64(segmente_offset), 0)
		}
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
	node.sender.Send(*addr, b)
}
