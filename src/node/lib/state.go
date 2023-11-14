package lib

import (
	"cc_project/helpers"
	"cc_project/protocol"
)

type State struct {
	MyFiles      map[string]protocol.FileHash // paths to my files
	Peers         helpers.Set[protocol.DeviceIdentifier]
	KnownFiles    map[protocol.FileHash]protocol.FileMetaData
}

func NewState() *State {
	state := &State{}
	state.MyFiles = make(map[string]protocol.FileHash)
	state.Peers = *(helpers.NewSet[protocol.DeviceIdentifier]())
	state.KnownFiles = make(map[protocol.FileHash]protocol.FileMetaData)
	return state
}

func (state *State)GetKnownFiles() map[protocol.FileHash]protocol.FileMetaData {
	return state.KnownFiles
}