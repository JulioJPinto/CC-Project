package state_manager

import (
	"cc_project/protocol/fstp"
	"fmt"
	"os"
)

type StateManager struct {
	filepath string
	State    *State
}

func NewManager(filepath string) *StateManager {
	return &StateManager{filepath, newState()}
}

func (m *StateManager) Connect() error {
	return nil
}

func (m *StateManager) RegisterDevice(device fstp.Device) error {
	m.State.Registered_nodes.Add(device)
	fmt.Println(*m.State.Registered_nodes)
	return nil
}

func (m *StateManager) DeviceIsRegistered(deviceID fstp.DeviceIdentifier) bool {
	f := func(d fstp.Device) bool { return d.GetIdentifier() == deviceID }
	return m.State.Registered_nodes.AnyMatch(f)
}

func (m *StateManager) LeaveNetwork(device fstp.DeviceIdentifier) error {
	f := func(d fstp.Device) bool { return d.GetIdentifier() == device }
	m.State.Registered_nodes.RemoveIf(f)
	delete(m.State.Nodes_segments, device)

	return nil
}

func (m *StateManager) RegisterFile(device fstp.DeviceIdentifier, file_info fstp.FileMetaData) error {
	if m.FileIsRegistered(file_info.Hash) {
		return ErrFileAlreadyRegistered
	}
	file_info.OriginatorIP = string(device)
	f := func(d fstp.Device) bool { return d.GetIdentifier() == device }
	if !m.State.Registered_nodes.AnyMatch(f) {
		return ErrNodeNotRegistered
	}
	m.State.Registered_files[file_info.Hash] = file_info
	for i, s_hash := range file_info.SegmentHashes {
		s := fstp.FileSegment{FirstByte: int64(i * fstp.SegmentLength), FileHash: file_info.Hash, Hash: s_hash}
		p, ok := m.State.Nodes_segments[device]
		if !ok {
			p = make([]fstp.FileSegment, 1)
		}
		m.State.Nodes_segments[device] = append(p, s)
	}
	return nil
}

func (m *StateManager) FileIsRegistered(hash fstp.FileHash) bool {
	_, ok := m.State.Registered_files[hash]
	return ok
}

func (m *StateManager) RegisterFileSegment(device fstp.DeviceIdentifier, file_segment fstp.FileSegment) error {
	x, ok := m.State.Registered_files[file_segment.FileHash]
	if !ok {
		return ErrFileDoesNotExist
	}
	offset := file_segment.FirstByte / fstp.SegmentLength
	if x.SegmentHashes[offset] != file_segment.Hash {
		return ErrInvalidSegmentHash
	}
	m.State.Nodes_segments[device] = append(m.State.Nodes_segments[device], file_segment)
	return nil
}

func (m *StateManager) BatchRegisterFileSegments(device fstp.DeviceIdentifier, segments []fstp.FileSegment) error {
	for _, segment := range segments {
		err := m.RegisterFileSegment(device, segment)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *StateManager) GetAllFiles() map[fstp.FileHash]fstp.FileMetaData {
	return m.State.Registered_files
}

// type WhoHasRespProps map[FileHash](map[DeviceIdentifier]FileSegment)

func (m *StateManager) WhoHasFile(hash fstp.FileHash) map[fstp.DeviceIdentifier][]fstp.FileSegment {
	ret := make(map[fstp.DeviceIdentifier][]fstp.FileSegment)
	for device, segments := range m.State.Nodes_segments {
		for _, segment := range segments {
			if segment.FileHash == hash {
				_, ok := ret[device]
				if !ok {
					ret[device] = []fstp.FileSegment{segment}
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
