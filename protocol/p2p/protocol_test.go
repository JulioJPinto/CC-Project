package p2p

import (
	"testing"
)

func TestSerializeAndDeserialize(t *testing.T) {
	h := P2PHeader{
		flags: struct {
			is_request  bool
			is_transfer bool
		}{true, false},
		file_id:  0xABCABC,
		fst_byte: 0x100,
		length:   0xA00,
	}

	serialized := Serialize(h)
	deserialized := Deserialize(serialized)

	if h != deserialized {
		t.Errorf("Serialization and deserialization failed.\nOriginal: %+v\nDeserialized: %+v", h, deserialized)
	}
}
