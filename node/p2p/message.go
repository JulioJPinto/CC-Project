package p2p

import (
	"cc_project/protocol"
	"encoding/binary"
)

type Header struct {
	IsRequest     bool              `json:"IsRequest"`
	Load          uint8             `json:"Load"`
	FileId        protocol.FileHash `json:"FileId"`
	SegmentOffset uint32            `json:"SegmentOffset"`
	Length        uint16            `json:"Length"`
	TimeStamp     uint32            `json:"TimeStamp"`
}

func (h *Header) Serialize() ([]byte, error) {
	load127 := h.Load << 1
	header := []byte{}
	// if h.IsRequest {
	// 	header[0] = byte(0b0000_0000)
	// } else {
	// 	header[0] = byte(0b1111_1111)
	// }
	if !h.IsRequest {
		load127 += 1
	}

	header[0] = (load127)

	binary.LittleEndian.PutUint32(header[1:5], uint32(h.FileId))
	binary.LittleEndian.PutUint32(header[5:9], h.SegmentOffset)
	binary.LittleEndian.PutUint16(header[9:11], h.Length)
	binary.LittleEndian.PutUint32(header[11:15], h.TimeStamp)

	return header, nil
}

func (h *Header) Deserialize(byteArray []byte) error {
	load127 := uint8(byteArray[0])
	h.IsRequest = (load127 % 2) == 0
	h.Load = load127 >> 1
	h.FileId = protocol.FileHash(binary.LittleEndian.Uint32(byteArray[1:5]))
	h.SegmentOffset = binary.LittleEndian.Uint32(byteArray[5:9])
	h.Length = binary.LittleEndian.Uint16(byteArray[9:11])
	h.TimeStamp = binary.LittleEndian.Uint32(byteArray[11:15])

	return nil
}

type Request Message
type Response Message

type Message struct {
	Header  `json:"Header"`
	Payload []byte //JSON Serializable
}

const PacketSize = 2048
const HeaderSize = 15

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
