package client_tracker

import "encoding/binary"

type FSTPHeader struct {
	flags struct {
		who_has bool
		i_have  bool
	}
	payload_size uint32
}

const FSTPHEaderSize = 5 // 5 bytes

type whoHasRequestBody struct {
	file_name   string
	byte_offset uint64
	byte_length uint64
}

func Deserialize(byteArray []byte) FSTPHeader {
	flags := byteArray[0]
	payload_size := binary.NativeEndian.Uint32(byteArray[1:5])
	return FSTPHeader{
		flags: struct {
			who_has bool
			i_have  bool
		}{
			(flags & 0b1000_0000) != 0,
			(flags & 0b0100_0000) != 0,
		},
		payload_size: payload_size,
	}
}

func Serialize(header FSTPHeader) []byte {
	b := byte(0)

	// Set the bits in the byte based on flag values
	if header.flags.who_has {
		b |= 0b10000000
	}
	if header.flags.i_have {
		b |= 0b01000000
	}

	// Return a byte slice containing the serialized header
	return []byte{b}
}
