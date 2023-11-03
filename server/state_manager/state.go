package state_manager

import (
	"cc_project/helpers"
	"cc_project/protocol/fstp"
)

type State struct {
	registered_nodes *helpers.Set[fstp.Device] 
	registered_files map[fstp.Hash]fstp.FileMetaData // mapeia a hash do ficheiro para os dados
	nodes_segments map[*fstp.DeviceIdentifier][]fstp.FileSegment
}

func newState() *State {
	s := &State{}
	s.registered_nodes = helpers.NewSet[fstp.Device](fstp.DeviceEqual)
	s.registered_files = make(map[fstp.Hash ]fstp.FileMetaData)
	s.nodes_segments = make(map[*fstp.DeviceIdentifier][]fstp.FileSegment)
	return s
}


var (
	ErrBadSchema         = helpers.WrapError{Msg: "bad schema"}
	ErrFileDoesNotExist  = helpers.WrapError{Msg: "file does not exist"}
	ErrInvalidParameters = helpers.WrapError{Msg: "invalid parameters"}
	ErrInvalidSegmentHash = helpers.WrapError{Msg: "invalid segment hash" }
)
