package p2p

import (
	"cc_project/protocol"
	"encoding/binary"
)

type Header struct {
	is_request bool
	file_id    protocol.FileHash
	fst_byte   uint32
	length     uint16
	port       uint16
}

func (h *Header) Serialize() ([]byte, error) {
	header := []byte{}
	if h.is_request {
		header[0] = byte(0b0000_0000)
	} else {
		header[0] = byte(0b1111_1111)
	}

	binary.LittleEndian.PutUint32(header[1:5], uint32(h.file_id))
	binary.LittleEndian.PutUint32(header[9:17], h.fst_byte)
	binary.LittleEndian.PutUint16(header[17:25], h.length)

	return header, nil
}

func (h *Header) Deserialize(byteArray []byte) error {

	h.file_id = protocol.FileHash(binary.LittleEndian.Uint32(byteArray[1:5]))
	h.fst_byte = binary.LittleEndian.Uint32(byteArray[9:17])
	h.length = binary.LittleEndian.Uint16(byteArray[17:25])

	flags := byteArray[0]
	h.is_request = flags == 0
	// Check the first byte to determine the flags (e.g., 0b1000_0000)
	// You can implement logic here to extract any additional fields based on flags.

	return nil
}

type Request Message
type Response Message

type Message struct {
	Header  Header
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

func GimmeFileSegmentRequest(port uint16,segment protocol.FileSegment) Request {
	header := Header{
		is_request: true,
		file_id:    segment.FileHash,
		fst_byte:   uint32(segment.FirstByte),
		length:     1,
		port: port,
	}
	return Request{Header: header, Payload: nil}
}
