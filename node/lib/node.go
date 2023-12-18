package lib

import (
	"cc_project/helpers"
	helpers_sync "cc_project/helpers/sync"
	"cc_project/protocol"
	"cc_project/protocol/fstp"
	"cc_project/protocol/p2p"
)

type Node struct {
	FSTPclient *fstp.Client
	P2PConfig  p2p.Config
	sender     *helpers.Uploader
	Downloads  *helpers_sync.Map[protocol.FileHash, *Downloader]
	NodeDir    string                                                   // path to the folder where files will be stored
	MyFiles    map[protocol.FileHash]string                             // paths to my files
	PeerStats  *helpers_sync.Map[protocol.DeviceIdentifier, *PeerStats] //
	KnownFiles map[protocol.FileHash]protocol.FileMetaData
}

func NewNode(fstp_config fstp.Config, p2p_config p2p.Config) (*Node, error) {
	client := &Node{}
	client.MyFiles = make(map[protocol.FileHash]string)
	client.KnownFiles = make(map[protocol.FileHash]protocol.FileMetaData)
	client.Downloads = &helpers_sync.Map[protocol.FileHash, *Downloader]{}
	client.PeerStats = &helpers_sync.Map[protocol.DeviceIdentifier, *PeerStats]{}
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
