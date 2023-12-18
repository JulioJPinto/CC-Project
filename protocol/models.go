package protocol

import "fmt"

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

func (meta *FileMetaData) FileSegments() []FileSegment {
	var segments []FileSegment

	for i, segmentHash := range meta.SegmentHashes {
		segment := FileSegment{
			BlockOffset: int64(i),
			FileHash:    meta.Hash,
			Hash:        segmentHash,
		}

		segments = append(segments, segment)
	}

	return segments
}

// length of a file segment in bytes
const SegmentMaxLength = 1024

type FileSegment struct {
	BlockOffset int64    `json:"BlockOffset"` //
	FileHash    FileHash `json:"FileID"`      // Foriegn Key refere um FileMetaData
	Hash        Hash     `json:"Hash"`
	// Length      uint16   `json:"Length"` //
}

func (fs FileSegment) String() string {
	return fmt.Sprintf("BlockOffset: %d, FileHash: %d, Hash: %d",
		fs.BlockOffset, fs.FileHash, fs.Hash)
}

func (s FileSegment) LastByte() int64 {
	return (s.BlockOffset+1)*SegmentMaxLength - 1
}
