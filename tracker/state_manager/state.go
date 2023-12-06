package state_manager

import (
	"cc_project/helpers"
	"cc_project/protocol"
	"encoding/json"
)

type State struct {
	Registered_nodes *helpers.Set[protocol.Device]                        `json:"registered_nodes"`
	Registered_files map[protocol.FileHash]protocol.FileMetaData          `json:"registered_files"` // mapeia a hash do ficheiro para os dados
	Nodes_segments   map[protocol.DeviceIdentifier][]protocol.FileSegment `json:"nodes_segments"`   // mapeia o
}

func newState() *State {
	s := &State{}
	s.Registered_nodes = helpers.NewSet[protocol.Device]()
	s.Registered_files = make(map[protocol.FileHash]protocol.FileMetaData)
	s.Nodes_segments = make(map[protocol.DeviceIdentifier][]protocol.FileSegment)
	return s
}

func (s *State) SegmentsNodes() map[protocol.FileSegment][]protocol.DeviceIdentifier {
	invertedMap := make(map[protocol.FileSegment][]protocol.DeviceIdentifier)

	for deviceID, fileSegments := range s.Nodes_segments {
		for _, fileSegment := range fileSegments {
			invertedMap[fileSegment] = append(invertedMap[fileSegment], deviceID)
		}
	}

	return invertedMap
}

func (s *State) Serialize() ([]byte, error) {
	return json.Marshal(*s)
}

func (s *State) Deserialize(data []byte) error {
	return json.Unmarshal(data, s)
}

var (
	ErrBadSchema             = helpers.WrapError{Msg: "bad schema"}
	ErrFileDoesNotExist      = helpers.WrapError{Msg: "file does not exist"}
	ErrFileAlreadyRegistered = helpers.WrapError{Msg: "file already registered"}
	ErrInvalidParameters     = helpers.WrapError{Msg: "invalid parameters"}
	ErrInvalidSegmentHash    = helpers.WrapError{Msg: "invalid segment hash"}
	ErrInvalidHeader         = helpers.WrapError{Msg: "invalid header"}
	ErrInvalidPayload        = helpers.WrapError{Msg: "invalid payload"}
	ErrNodeNotRegistered     = helpers.WrapError{Msg: "node not yet registered"}
)
