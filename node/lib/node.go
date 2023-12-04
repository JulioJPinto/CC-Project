package lib

import (
	"cc_project/helpers"
	"cc_project/node/lib/udp_uploader"
	"cc_project/node/p2p"
	"cc_project/protocol"
	"cc_project/protocol/fstp"
	"encoding/json"
	"fmt"
	"time"
)

type Node struct {
	FSTPclient *fstp.Client
	P2PConfig  p2p.Config
	sender     *udp_uploader.Uploader
	Chanels    *helpers.SyncMap[protocol.FileHash, chan p2p.Message]

	MyFiles map[string]protocol.FileHash // paths to my files
	Peers   helpers.SyncMap[protocol.DeviceIdentifier, PeerStats]
	// PeerStats helpers.SyncMap[protocol.DeviceIdentifier]
	KnownFiles map[protocol.FileHash]protocol.FileMetaData
}

func NewNode(fstp_config fstp.Config, p2p_config p2p.Config) (*Node, error) {

	client := &Node{}
	client.MyFiles = make(map[string]protocol.FileHash)
	client.Peers = *(helpers.NewSyncMap[protocol.DeviceIdentifier, PeerStats](protocol.HashDeviceIdentifier))
	client.KnownFiles = make(map[protocol.FileHash]protocol.FileMetaData)

	client.Chanels = helpers.NewSyncMap[protocol.FileHash, chan p2p.Message](protocol.HashFileHash)
	var fstp_client *fstp.Client
	var err error = nil
	if fstp_client, err = fstp.NewClient(fstp_config); err != nil {
		return nil, err
	}
	client.FSTPclient = fstp_client
	client.P2PConfig = p2p_config
	client.sender = udp_uploader.NewUploader(5)
	return client, err
}

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
	go node.send_segment_requests(p)
	go node.await_segment_responses(file_meta_data)

	return nil
}

func (node *Node) send_segment_requests(m map[protocol.DeviceIdentifier][]protocol.FileSegment) {
	for id, segments := range m {
		for _, segment := range segments {
			node.RequestSegment(id, segment)
		}
	}
}

func (node *Node) await_segment_responses(file protocol.FileMetaData) {
	ch, _ := node.Chanels.Get(file.Hash)
	for msg := range ch {
		println(json.Marshal(msg))
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
	node.sender.Send(string(peer), b)
}
