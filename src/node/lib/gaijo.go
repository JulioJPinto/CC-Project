package lib

import (
	"cc_project/helpers"
	"cc_project/node/p2p"
	"cc_project/protocol"
	"cc_project/protocol/fstp"
	"fmt"
	"net"
	"time"
)

type Gaijo struct {
	FSTPclient *fstp.Client
	P2PConfig  p2p.Config
	udp_conn   *net.UDPConn
	Chanels    *helpers.SyncMap[protocol.FileHash, chan p2p.Message]

	MyFiles    map[string]protocol.FileHash // paths to my files
	Peers      helpers.Set[protocol.DeviceIdentifier]
	KnownFiles map[protocol.FileHash]protocol.FileMetaData
}

func NewGaijo(fstp_config fstp.Config, p2p_config p2p.Config) (*Gaijo, error) {
	client := &Gaijo{}
	client.MyFiles = make(map[string]protocol.FileHash)
	client.Peers = *(helpers.NewSet[protocol.DeviceIdentifier]())
	client.KnownFiles = make(map[protocol.FileHash]protocol.FileMetaData)

	client.Chanels = helpers.NewSyncMap[protocol.FileHash, chan p2p.Message]()

	var err error
	client.FSTPclient, err = fstp.NewClient(fstp_config)
	client.P2PConfig = p2p_config
	return client, err
}

func (c *Gaijo) DownloadFile(file_hash protocol.FileHash) error {
	if _, ok := c.Chanels.Get(file_hash); ok {
		return fmt.Errorf("download already in progress")
	}

	// file_meta_data,ok := c.KnownFiles.get(file_hash)
	file_meta_data, ok := c.KnownFiles[file_hash]

	if !ok {
		c.FetchFiles(nil)
		file_meta_data, ok = c.KnownFiles[file_hash]
		if !ok{return fmt.Errorf("files does not exist: %v", file_hash)}
	}

	resp, err := c.FSTPclient.Request(fstp.NewWhoHasRequest(fstp.WhoHasReqProps{File: file_hash}))
	if err != nil {
		return err
	}



	if resp.Header.Flags == fstp.ErrResp {
		err_resp := resp.Payload.(*fstp.ErrorResponse)
		return fmt.Errorf(err_resp.Err)
	}
	go c.SendSegmentRequests(file_meta_data)
	go c.AwaitSegmentResponses(file_meta_data)
	// device_segments, _ := resp.Payload.(*fstp.WhoHasRespProps)
	// c.Distribute(*device_segments, file_meta_data)

	return nil
}

func (g *Gaijo) SendSegmentRequests(file protocol.FileMetaData) {
}

func (g *Gaijo) AwaitSegmentResponses(file protocol.FileMetaData) {

}

func (g *Gaijo) Distribute(device_segments map[protocol.DeviceIdentifier][]protocol.FileSegment, metadata protocol.FileMetaData) map[int64]protocol.DeviceIdentifier {
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

func (g *Gaijo) RequestSegment(peer *net.UDPAddr, segment protocol.FileSegment) {
	timestamp := helpers.TrunkI64(time.Now().UnixMilli())
	req := p2p.Message(p2p.GimmeFileSegmentRequest(segment,uint32(timestamp)))
	x := &req
	b, err := x.Serialize()
	if err != nil {
		return
	}
	g.udp_conn.Write(b)

}
