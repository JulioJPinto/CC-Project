package db

import (
	"cc_project/helpers"
	"net"
)

var db Database = nil

// Database interface
type Database interface {
	Connect() error
	Close() error
	RegisterDevice(DeviceData) error
	RegisterFile(FileMetaData) error
	RegisterFileSegment(net.IP, FileSegment) error
	GetAllDevices() []DeviceData
}

var (
	ErrBadSchema         = helpers.WrapError{Msg: "bad schema"}
	ErrFileDoesNotExist  = helpers.WrapError{Msg: "file does not exist"}
	ErrInvalidParameters = helpers.WrapError{Msg: "invalid parameters"}
)
