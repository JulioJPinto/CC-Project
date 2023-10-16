package db

import "cc_project/error_helpers"

var db Database = nil

// Database interface
type Database interface {
	Connect() error
	Close() error
	RegisterDevice(DeviceData) error
	RegisterFile(FileMetaData) error
	GetAllDevices() []DeviceData
}

var (
	ErrBadSchema  = error_helpers.WrapError{Msg: "bad schema"}
	ErrFileDoesNotExist = error_helpers.WrapError{Msg: "file does not exist"}
	ErrInvalidParameters = error_helpers.WrapError{Msg: "invalid parameters"}
)
