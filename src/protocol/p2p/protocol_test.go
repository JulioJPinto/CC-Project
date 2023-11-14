package p2p

// import (
// 	"testing"
// )

// func TestSerializeAndDeserialize(t *testing.T) {
// 	h := Header{
// 		flags: struct {
// 			is_request  bool
// 			is_transfer bool
// 		}{true, false},
// 		file_id:  0xABCABC,
// 		fst_byte: 0x100,
// 		length:   0xA00,
// 	}

// 	serialized, err := h.Serialize()
// 	if err != nil {
// 		t.Errorf("Serialization and deserialization failed.\nOriginal: %+v\nDeserialized: %+v", h, deserialized)
// 	}
// 	var deserialized *Message
// 	deserialized.Deserialize(serialized)

// 	if h != deserialized {
// 		t.Errorf("Serialization and deserialization failed.\nOriginal: %+v\nDeserialized: %+v", h, deserialized)
// 	}
// }
