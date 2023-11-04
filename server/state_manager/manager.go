package state_manager

import (
	"cc_project/protocol/fstp"
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
	return nil
}

func (m *StateManager) DeviceIsRegistered(deviceID fstp.DeviceIdentifier) bool {
	f := func(d fstp.Device) bool { return d.GetIdentifier() == deviceID }
	return m.State.Registered_nodes.AnyMatch(f)
}

func (m *StateManager) RegisterFile(device fstp.DeviceIdentifier, file_info fstp.FileMetaData) error {
	file_info.OriginatorIP = string(device)
	f := func(d fstp.Device) bool { return d.GetIdentifier() == device }
	if !m.State.Registered_nodes.AnyMatch(f) {
		return ErrNodeNotRegistered
	}
	m.State.Registered_files[file_info.Hash] = file_info
	return nil
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
