package lib

import (
	"cc_project/helpers"
	"cc_project/protocol/p2p"
	"cc_project/protocol"
	"cc_project/protocol/fstp"
	"fmt"
	"net"
)

type Gaijo struct {
	FSTPclient  *fstp.Client
	P2PListener *p2p.Listener
	State       State
	UdpCenas    struct {
		adress net.UDPAddr
	}
}

func NewGaijo(fstp_config fstp.Config, p2p_config p2p.Config) (*Gaijo, error) {
	client := &Gaijo{}
	client.State.MyFiles = make(map[string]protocol.FileHash)
	client.State.Peers = *(helpers.NewSet[protocol.DeviceIdentifier]())
	client.State.KnownFiles = make(map[protocol.FileHash]protocol.FileMetaData)

	var err error
	client.FSTPclient, err = fstp.NewClient(fstp_config)
	client.P2PListener = p2p.P2PConn(p2p_config)
	return client, err
}

func (c *Gaijo) Listen(){
	
}

func (c *Gaijo) DownloadFile(fileHash protocol.FileHash) error {
	resp, err := c.FSTPclient.Request(fstp.NewWhoHasRequest(fstp.WhoHasReqProps{File: fileHash}))
	if err != nil {
		return err
	}
	if resp.Header.Flags == fstp.ErrResp {
		// depois vê-se
		pay, _ := resp.Payload.(*fstp.ErrorResponse)
		return fmt.Errorf(pay.Err)
	}
	_, _ = resp.Payload.(*fstp.WhoHasRespProps)
	return nil
}

func DownloadSegment(segment protocol.FileSegment, peer *net.UDPAddr) {
	var port uint16 = 9091
	_ = p2p.GimmeFileSegmentRequest(port, segment)

}
