package p2p

import (
	"encoding/binary"
)

type P2PHeader struct {
	flags struct {
		is_request  bool
		is_transfer bool
	}
	file_id  uint64
	fst_byte uint64
	length   uint64
}

func Serialize(h P2PHeader) [25]byte {
	header := [25]byte{}

	binary.LittleEndian.PutUint64(header[0:8], h.file_id)
	binary.LittleEndian.PutUint64(header[8:16], h.fst_byte)
	binary.LittleEndian.PutUint64(header[16:24], h.length)

	header[24] = 0b1000_0000
	return header
}

func Deserialize(byteArray [25]byte) P2PHeader {
	var h P2PHeader

	h.file_id = binary.LittleEndian.Uint64(byteArray[0:8])
	h.fst_byte = binary.LittleEndian.Uint64(byteArray[8:16])
	h.length = binary.LittleEndian.Uint64(byteArray[16:24])

	flags := byteArray[24]
	h.flags.is_request = (flags & 0b1000_0000) != 0
	h.flags.is_transfer = (flags & 0b0100_0000) != 0
	// Check the last byte to determine the flags (e.g., 0b1000_0000)
	// You can implement logic here to extract any additional fields based on flags.

	return h
}

