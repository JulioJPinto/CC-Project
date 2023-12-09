package lib

import (
	"cc_project/helpers"
	"cc_project/protocol"
	"cc_project/protocol/fstp"
	"cc_project/protocol/p2p"
	"sync"
)

type Node struct {
	FSTPclient *fstp.Client
	P2PConfig  p2p.Config
	sender     *helpers.Uploader
	Downloads  helpers.SyncList[Downloader]
	Chanels    *sync.Map // [protocol.FileHash, chan p2p.Message]
	NodeDir    string                       // path to the folder where files will be stored
	MyFiles    map[protocol.FileHash]string // paths to my files
	Peers      helpers.SyncMap[protocol.DeviceIdentifier, PeerStats]
	// PeerStats helpers.SyncMap[protocol.DeviceIdentifier]
	KnownFiles map[protocol.FileHash]protocol.FileMetaData
}

func NewNode(fstp_config fstp.Config, p2p_config p2p.Config) (*Node, error) {

	client := &Node{}
	client.MyFiles = make(map[protocol.FileHash]string)
	// client.Peers = *(helpers.NewSyncMap[protocol.DeviceIdentifier, PeerStats](protocol.HashDeviceIdentifier))
	client.Peers = *(helpers.NewSyncMap[protocol.DeviceIdentifier, PeerStats](protocol.HashDeviceIdentifier))
	client.KnownFiles = make(map[protocol.FileHash]protocol.FileMetaData)

	client.Chanels = &sync.Map{} //helpers.NewSyncMap[protocol.FileHash, chan p2p.Message](protocol.HashFileHash)
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
