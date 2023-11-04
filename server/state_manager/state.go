package state_manager

import (
	"cc_project/helpers"
	"cc_project/protocol/fstp"
	"encoding/json"
)

type State struct {
	Registered_nodes *helpers.Set[fstp.Device]                    `json:"registered_nodes"`
	Registered_files map[fstp.Hash]fstp.FileMetaData              `json:"registered_files"` // mapeia a hash do ficheiro para os dados
	Nodes_segments   map[fstp.DeviceIdentifier][]fstp.FileSegment `json:"nodes_segments"`   // mapeia o
}

func newState() *State {
	s := &State{}
	s.Registered_nodes = helpers.NewSet[fstp.Device]()
	s.Registered_files = make(map[fstp.Hash]fstp.FileMetaData)
	s.Nodes_segments = make(map[fstp.DeviceIdentifier][]fstp.FileSegment)
	return s
}

func (s *State) Serialize() ([]byte, error) {
	return json.Marshal(*s)
}

func (s *State) Deserialize(data []byte) error {
	return json.Unmarshal(data, s)
}

var (
	ErrBadSchema          = helpers.WrapError{Msg: "bad schema"}
	ErrFileDoesNotExist   = helpers.WrapError{Msg: "file does not exist"}
	ErrInvalidParameters  = helpers.WrapError{Msg: "invalid parameters"}
	ErrInvalidSegmentHash = helpers.WrapError{Msg: "invalid segment hash"}
	ErrNodeNotRegistered  = helpers.WrapError{Msg: "node not yet registered"}
)
