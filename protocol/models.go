package protocol

import (
	"cc_project/helpers"
)

type Device struct {
	IP string `json:"IP"`
}

type DeviceIdentifier string // ip

func (d *Device) GetIdentifier() DeviceIdentifier {
	return DeviceIdentifier(d.IP)
}

type FileInfo struct {
	Id uint64 `json:"Id"`
}

type Hash uint32
type FileHash uint32
type FileMetaData struct {
	Hash          FileHash `json:"Hash"` // Primary Key
	Name          string   `json:"Name"`
	Length        int32    `json:"Length"`
	OriginatorIP  string   `json:"originatorIP"` //IP
	SegmentHashes []Hash   `json:"SegmentHashes"`
}

// length of a file segment in bytes
const SegmentLength = 1024

type FileSegment struct {
	BlockOffset int64    `json:"BlockOffset"` //
	FileHash    FileHash `json:"FileID"`      // Foriegn Key refere um FileMetaData
	Hash        Hash     `json:"Hash"`
}

func (s FileSegment) LastByte() int64 {
	return (s.BlockOffset+1)*SegmentLength - 1
}

func HashDeviceIdentifier(id DeviceIdentifier) uint64 {
	return helpers.HashString(string(id))
}

func HashFileHash(id FileHash) uint64 {
	return helpers.HashUint32(uint32(id))
}
