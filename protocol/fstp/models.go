package fstp

import "net"

type Device struct {
	IP net.IP `json:"IP"`
}

type DeviceIdentifier net.IP

func (d *Device) GetIdentifier() DeviceIdentifier {
	return DeviceIdentifier(d.IP)
}

func DeviceEqual(d1 Device, d2 Device) bool {
	return d1.IP.String() == d2.IP.String()
}

type FileInfo struct {
	Id uint64 `json:"Id"`
}

type Hash uint32
type FileMetaData struct {
	Hash          Hash   `json:"Hash"` // Primary Key
	Name          string `json:"Name"`
	Length        int32  `json:"Length"`
	OriginatorIP  string `json:"origintorIP"` //IP
	SegmentHashes []Hash `json:"SegmentHashes"`
}

// length of a file segment in bytes
const SegmentLength = 128

type FileSegment struct {
	FirstByte int64 `json:"FirstByte"`
	FileHash  Hash  `json:"FileID"` // Foriegn Key refere um FileMetaData
	Hash      Hash  `json:"Hash"`
}

func (s FileSegment) LastByte() int64 {
	return s.FirstByte + SegmentLength - 1
}
