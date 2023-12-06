package state_manager

import (
	"cc_project/protocol"
	"fmt"
	"os"
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
	m.State.Registered_nodes.Add(device)
	fmt.Println(*m.State.Registered_nodes)
	return nil
}

func (m *StateManager) DeviceIsRegistered(deviceID protocol.DeviceIdentifier) bool {
	f := func(d protocol.Device) bool { return d.GetIdentifier() == deviceID }
	return m.State.Registered_nodes.AnyMatch(f)
}

func (m *StateManager) LeaveNetwork(device protocol.DeviceIdentifier) error {
	f := func(d protocol.Device) bool { return d.GetIdentifier() == device }
	m.State.Registered_nodes.RemoveIf(f)
	// segments_nodes := m.SegmentsNodes()
	delete(m.State.Nodes_segments, device)
	
	return nil
}


func (m *StateManager) RegisterFile(device protocol.DeviceIdentifier, file_info protocol.FileMetaData) error {
	if m.FileIsRegistered(file_info.Hash) {
		return ErrFileAlreadyRegistered
	}
	file_info.OriginatorIP = string(device)
	f := func(d protocol.Device) bool { return d.GetIdentifier() == device }
	if !m.State.Registered_nodes.AnyMatch(f) {
		return ErrNodeNotRegistered
	}
	m.State.Registered_files[file_info.Hash] = file_info
	for i, s_hash := range file_info.SegmentHashes {
		s := protocol.FileSegment{BlockOffset: int64(i), FileHash: file_info.Hash, Hash: s_hash}
		p, ok := m.State.Nodes_segments[device]
		if !ok {
			p = make([]protocol.FileSegment, 1)
		}
		m.State.Nodes_segments[device] = append(p, s)
	}
	return nil
}

func (m *StateManager) FileIsRegistered(hash protocol.FileHash) bool {
	_, ok := m.State.Registered_files[hash]
	return ok
}

func (m *StateManager) RegisterFileSegment(device protocol.DeviceIdentifier, file_segment protocol.FileSegment) error {
	x, ok := m.State.Registered_files[file_segment.FileHash]
	if !ok {
		return ErrFileDoesNotExist
	}
	offset := file_segment.BlockOffset / protocol.SegmentLength
	if x.SegmentHashes[offset] != file_segment.Hash {
		return ErrInvalidSegmentHash
	}
	m.State.Nodes_segments[device] = append(m.State.Nodes_segments[device], file_segment)
	return nil
}

func (m *StateManager) BatchRegisterFileSegments(device protocol.DeviceIdentifier, segments []protocol.FileSegment) error {
	for _, segment := range segments {
		err := m.RegisterFileSegment(device, segment)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *StateManager) GetAllFiles() map[protocol.FileHash]protocol.FileMetaData {
	return m.State.Registered_files
}

// type WhoHasRespProps map[FileHash](map[DeviceIdentifier]FileSegment)

func (m *StateManager) WhoHasFile(hash protocol.FileHash) map[protocol.DeviceIdentifier][]protocol.FileSegment {
	ret := make(map[protocol.DeviceIdentifier][]protocol.FileSegment)
	for device, segments := range m.State.Nodes_segments {
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
	}
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
