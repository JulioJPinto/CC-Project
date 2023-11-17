package p2p

import (
	"cc_project/protocol"
	"encoding/binary"
)

type Header struct {
	IsRequest     bool
	FileId        protocol.FileHash
	SegmentOffset uint32
	Length        uint16
	TimeStamp     uint32
}

func (h *Header) Serialize() ([]byte, error) {
	header := []byte{}
	if h.IsRequest {
		header[0] = byte(0b0000_0000)
	} else {
		header[0] = byte(0b1111_1111)
	}

	binary.LittleEndian.PutUint32(header[1:5], uint32(h.FileId))
	binary.LittleEndian.PutUint32(header[5:9], h.SegmentOffset)
	binary.LittleEndian.PutUint16(header[9:11], h.Length)
	binary.LittleEndian.PutUint32(header[11:15], h.TimeStamp)
	

	return header, nil
}

func (h *Header) Deserialize(byteArray []byte) error {
	flags := byteArray[0]
	h.IsRequest = flags == 0

	h.FileId = protocol.FileHash(binary.LittleEndian.Uint32(byteArray[1:5]))
	h.SegmentOffset = binary.LittleEndian.Uint32(byteArray[5:9] )
	h.Length = binary.LittleEndian.Uint16(byteArray[9:11])
	h.TimeStamp = binary.LittleEndian.Uint32(byteArray[11:15])

	// Check the first byte to determine the flags (e.g., 0b1000_0000)
	// You can implement logic here to extract any additional fields based on flags.

	return nil
}

type Request Message
type Response Message

type Message struct {
	Header
	Payload []byte //JSON Serializable
}

const PacketSize = 2048
const HeaderSize = 25

func (m *Message) Serialize() ([]byte, error) {
	header_bytes, err := m.Header.Serialize()
	if err != nil {
		return nil, err
	}
	return append(header_bytes, m.Payload...), nil
}

func (m *Message) Deserialize(bytes []byte) error {
	if err := m.Header.Deserialize(bytes[:HeaderSize]); err != nil {
		return err
	}
	if len(bytes) > HeaderSize {
		m.Payload = bytes[HeaderSize:]
	} else {
		m.Payload = nil
	}
	return nil
}

func GimmeFileSegmentRequest(segment protocol.FileSegment, time_stamp uint32) Request {
	header := Header{
		IsRequest:     true,
		FileId:        segment.FileHash,
		TimeStamp:     time_stamp,
		SegmentOffset: uint32(segment.BlockOffset),
		Length:        1,
	}
	return Request{Header: header, Payload: nil}
}
