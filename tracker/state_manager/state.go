package state_manager

import (
	"cc_project/helpers"
	helpers_sync "cc_project/helpers/sync"
	"cc_project/protocol"
	"encoding/json"
	"fmt"
)

type State struct {
	RegisteredNodes *helpers_sync.Set[protocol.Device]                                   `json:"registered_nodes"`
	RegisteredFiles *helpers_sync.Map[protocol.FileHash, protocol.FileMetaData]          `json:"registered_files"` // mapeia a hash do ficheiro para os dados
	NodesSegments   *helpers_sync.Map[protocol.DeviceIdentifier, []protocol.FileSegment] `json:"nodes_segments"`   // mapeia o
}

func newState() *State {
	s := &State{}
	s.RegisteredNodes = &helpers_sync.Set[protocol.Device]{}
	s.RegisteredFiles = &helpers_sync.Map[protocol.FileHash, protocol.FileMetaData]{}
	s.NodesSegments = &helpers_sync.Map[protocol.DeviceIdentifier, []protocol.FileSegment]{}
	return s
}

func (s *State) SegmentsNodes() map[protocol.FileSegment][]protocol.DeviceIdentifier {
	invertedMap := make(map[protocol.FileSegment][]protocol.DeviceIdentifier)

	s.NodesSegments.Range(func(deviceID protocol.DeviceIdentifier, fileSegments []protocol.FileSegment) bool {
		for _, fileSegment := range fileSegments {
			invertedMap[fileSegment] = append(invertedMap[fileSegment], deviceID)
		}
		return true
	})

	return invertedMap
}

func (s *State) Serialize() ([]byte, error) {
	return json.Marshal(*s)
}

func (s *State) Deserialize(data []byte) error {
	return json.Unmarshal(data, s)
}

func (s *State) Print() {
	fmt.Println("Registered Nodes:")
	for _, node := range s.RegisteredNodes.List() {
		fmt.Printf(" - %s\n", node)
	}

	fmt.Println("\nRegistered Files:")
	s.RegisteredFiles.Range(func(fileHash protocol.FileHash, fileMetaData protocol.FileMetaData) bool {
		fmt.Printf(" - Hash: %d, Name: %s, Length: %d\n", fileHash, fileMetaData.Name, fileMetaData.Length)
		return true
	})

	fmt.Println("\nNodes Segments:")
	s.NodesSegments.Range(func(deviceID protocol.DeviceIdentifier, fileSegments []protocol.FileSegment) bool {
		fmt.Printf(" - Device ID: %s\n", deviceID)
		fmt.Println("   - Segments:")
		for _, segment := range fileSegments {
			fmt.Printf("     - BlockOffset: %d, FileHash: %d, SegmentHash: %d\n",
				segment.BlockOffset, segment.FileHash, segment.Hash)
		}
		return true
	})
}

func (s *State) String() string {
	var result string

	result += "Registered Nodes:\n"
	for _, node := range s.RegisteredNodes.List() {
		result += fmt.Sprintf(" - %s\n", node)
	}

	result += "\nRegistered Files:\n"
	s.RegisteredFiles.Range(func(fileHash protocol.FileHash, fileMetaData protocol.FileMetaData) bool {
		result += fmt.Sprintf(" - Hash: %d, Name: %s, Length: %d\n", fileHash, fileMetaData.Name, fileMetaData.Length)
		return true
	})

	result += "\nNodes Segments:\n"
	s.NodesSegments.Range(func(deviceID protocol.DeviceIdentifier, fileSegments []protocol.FileSegment) bool {
		result += fmt.Sprintf(" - Device ID: %s\n", deviceID)
		result += "   - Segments:\n"
		for _, segment := range fileSegments {
			result += fmt.Sprintf("     - BlockOffset: %d, FileHash: %d, SegmentHash: %d\n",
				segment.BlockOffset, segment.FileHash, segment.Hash)
		}
		return true
	})

	return result
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
