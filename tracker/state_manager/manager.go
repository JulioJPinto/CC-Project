package state_manager

import (
	"cc_project/helpers"
	"cc_project/protocol"
	"fmt"
	"os"

	"github.com/fatih/color"
)

type StateManager struct {
	filepath string
	*State
}

func NewManager(filepath string) *StateManager {
	return &StateManager{filepath, newState()}
}

func (m *StateManager) Connect() error {
	return nil
}

func (m *StateManager) RegisterDevice(device protocol.Device) error {
	m.State.RegisteredNodes.Add(device)
	fmt.Println(m.State.RegisteredNodes.List())
	return nil
}

func (m *StateManager) DeviceIsRegistered(deviceID protocol.DeviceIdentifier) bool {
	f := func(d protocol.Device) bool { return d.GetIdentifier() == deviceID }
	return m.State.RegisteredNodes.AnyMatch(f)
}

func (m *StateManager) LeaveNetwork(device protocol.DeviceIdentifier) error {
	f := func(d protocol.Device) bool { return d.GetIdentifier() == device }
	m.State.RegisteredNodes.RemoveIf(f)
	m.State.NodesSegments.Delete(device)

	print(m.SegmentsNodesString())
	exsiting_segments := helpers.MapKeys(m.SegmentsNodes())
	println("EXISTING: ", exsiting_segments.String(), "\n\n")
	m.State.RegisteredFiles.Range(func(key protocol.FileHash, val protocol.FileMetaData) bool {
		segments := helpers.NewSetFromSlice[protocol.FileSegment](val.FileSegments())
		print("CHECKING: ", segments.String(), "\n\n")

		if !segments.IsSubset(exsiting_segments) {
			color.Red("DELETING: " + fmt.Sprint(key))
			m.State.RegisteredFiles.Delete(key)
		}
		return true
	})
	return nil
}

func (m *StateManager) RegisterFile(device protocol.DeviceIdentifier, file_info protocol.FileMetaData) error {
	if !m.FileIsRegistered(file_info.Hash) {
		file_info.OriginatorIP = string(device)
		f := func(d protocol.Device) bool { return d.GetIdentifier() == device }
		if !m.State.RegisteredNodes.AnyMatch(f) {
			return ErrNodeNotRegistered
		}
		m.State.RegisteredFiles.Store(file_info.Hash, file_info)
	}
	for i, s_hash := range file_info.SegmentHashes {
		s := protocol.FileSegment{BlockOffset: int64(i), FileHash: file_info.Hash, Hash: s_hash}
		p, ok := m.State.NodesSegments.Load(device)
		if !ok {
			p = make([]protocol.FileSegment, 1)
		}
		m.State.NodesSegments.Store(device, append(p, s))
	}
	return nil
}

func (m *StateManager) FileIsRegistered(hash protocol.FileHash) bool {
	_, ok := m.State.RegisteredFiles.Load(hash)
	return ok
}

func (m *StateManager) RegisterFileSegment(device protocol.DeviceIdentifier, file_segment protocol.FileSegment) error {
	x, ok := m.State.RegisteredFiles.Load(file_segment.FileHash)
	if !ok {
		return ErrFileDoesNotExist
	}
	offset := file_segment.BlockOffset / protocol.SegmentMaxLength
	if x.SegmentHashes[offset] != file_segment.Hash {
		return ErrInvalidSegmentHash
	}
	p, _ := m.State.NodesSegments.Load(device) // WARNING WE BE IGNORIN STUFFS
	m.State.NodesSegments.Store(device, append(p, file_segment))
	return nil
}

func (m *StateManager) BatchRegisterFileSegments(device protocol.DeviceIdentifier, segments []protocol.FileSegment) error {
	for _, segment := range segments {
		err := m.RegisterFileSegment(device, segment)
		if err != nil {
			continue
		}
	}
	return nil
}

func (m *StateManager) GetAllFiles() map[protocol.FileHash]protocol.FileMetaData {

	return m.State.RegisteredFiles.ToMap()
}

// type WhoHasRespProps map[FileHash](map[DeviceIdentifier]FileSegment)

func (m *StateManager) WhoHasFile(hash protocol.FileHash) map[protocol.DeviceIdentifier][]protocol.FileSegment {
	ret := make(map[protocol.DeviceIdentifier][]protocol.FileSegment)
	m.State.NodesSegments.Range(func(device protocol.DeviceIdentifier, segments []protocol.FileSegment) bool {
		for _, segment := range segments {
			if segment.FileHash == hash {
				_, ok := ret[device]
				if !ok {
					ret[device] = []protocol.FileSegment{segment}
				} else {
					ret[device] = append(ret[device], segment)
				}
			}
		}
		return true
	})
	return ret
}

func (m *StateManager) DumpToFile() error {
	// Serialize the state to bytes
	bytes, err := m.State.Serialize()
	if err != nil {
		return err
	}
	// Write the bytes to the file specified by m.filepath
	err = os.WriteFile(m.filepath, bytes, 0644)
	if err != nil {
		return err
	}

	return nil
}
