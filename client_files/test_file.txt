package state_manager

import (
	"cc_project/protocol/fstp"
	"net"
)

type StateManager struct {
	filepath string
	state    *State
}

func NewManager(filepath string) *StateManager {
	return &StateManager{filepath, newState()}
}

func (m *StateManager) Connect() error {
	return nil
}

func (m *StateManager) RegisterDevice(device fstp.Device) error {
	m.state.registered_nodes.Add(device)
	return nil
}

func (m *StateManager) RegisterFile(device fstp.DeviceIdentifier, file_info fstp.FileMetaData) error {
	file_info.OriginatorIP = net.IP(device).String()
	m.state.registered_files[file_info.Hash] = file_info
	return nil
}

func (m *StateManager) RegisterFileSegment(device fstp.DeviceIdentifier, file_segment fstp.FileSegment) error {
	x, ok := m.state.registered_files[file_segment.FileHash]
	if !ok {
		return ErrFileDoesNotExist
	}
	offset := file_segment.FirstByte/fstp.SegmentLength
	if x.SegmentHashes[offset] != file_segment.Hash {
		return ErrInvalidSegmentHash
	}
	m.state.nodes_segments[&device] = append(m.state.nodes_segments[&device], file_segment)
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
