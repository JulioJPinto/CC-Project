package db

var db Database = nil

// Database interface
type Database interface {
	Connect() error
	Close() error
	RegisterDevice(DeviceData) error
	GetAllDevices() []DeviceData
}
