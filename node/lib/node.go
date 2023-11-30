package lib

import (
	"cc_project/helpers"
	"cc_project/protocol/p2p"
	"cc_project/protocol"
	"cc_project/protocol/fstp"
)

type Node struct {
	FSTPclient *fstp.Client
	P2PConfig  p2p.Config
	sender     *helpers.Uploader
	Chanels    *helpers.SyncMap[protocol.FileHash, chan p2p.Message]

	MyFiles map[string]protocol.FileHash // paths to my files
	Peers   helpers.SyncMap[protocol.DeviceIdentifier, Stats]
	// PeerStats helpers.SyncMap[protocol.DeviceIdentifier]
	KnownFiles map[protocol.FileHash]protocol.FileMetaData
}

func NewNode(fstp_config fstp.Config, p2p_config p2p.Config) (*Node, error) {

	client := &Node{}
	client.MyFiles = make(map[string]protocol.FileHash)
	client.Peers = *(helpers.NewSyncMap[protocol.DeviceIdentifier, Stats]())
	client.KnownFiles = make(map[protocol.FileHash]protocol.FileMetaData)

	client.Chanels = helpers.NewSyncMap[protocol.FileHash, chan p2p.Message]()
	var fstp_client *fstp.Client
	var err error = nil
	if fstp_client, err = fstp.NewClient(fstp_config); err != nil {
		return nil, err
	}
	client.FSTPclient = fstp_client
	client.P2PConfig = p2p_config
	client.sender = helpers.NewUploader(5)
	return client, err
}
