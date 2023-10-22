package fstp

import (
	"cc_project/helpers"
	"encoding/binary"
	"unsafe"
)

type FSTPHeader struct {
	flags        uint8
	payload_size uint32
}

// FSTPrequest implements message
type FSTPrequest struct {
	header  FSTPHeader
	payload helpers.Serializable
}

type FSTPresponse struct {
	header  FSTPHeader
	payload helpers.Serializable
}

const FSTPHEaderSize = 5 // 5 bytes

type File_info struct {
	Id uint64
}

func (message *FSTPresponse) Deserialize(byteArray []byte) error {
	message.header.flags = byteArray[0]
	var payload = helpers.SerializableMap{}
	_ = payload.Deserialize(byteArray[FSTPHEaderSize:]) // WARN!!! ignoring deserialization errors like a chad
	// fmt.Println(byteArray[FSTPHEaderSize:], payload)
	message.payload = &payload
	return nil
}

func (message *FSTPresponse) Serialize() ([]byte, error) {
	tag := message.header.flags

	payload, _ := message.payload.Serialize() // WARN!!! ignoring serialization errors like a chad
	payload_size := uint16(len(payload))      // Mudar o tamanho do int se necessário, improvável af tho (*)
	serialized_payload_size := make([]byte, unsafe.Sizeof(payload_size))
	binary.NativeEndian.PutUint16(serialized_payload_size, payload_size) // (*) aqui também
	// Return a byte slice containing the serialized header
	return append(append([]byte{tag}, serialized_payload_size...), payload...), nil
}

func (message *FSTPrequest) Deserialize(byteArray []byte) error {
	message.header.flags = byteArray[0]
	var payload = helpers.SerializableMap{}
	_ = payload.Deserialize(byteArray[FSTPHEaderSize:]) // WARN!!! ignoring deserialization errors like a chad
	// fmt.Println(byteArray[FSTPHEaderSize:], payload)
	message.payload = &payload
	return nil
}

func (message *FSTPrequest) Serialize() ([]byte, error) {
	tag := message.header.flags

	payload, _ := message.payload.Serialize() // WARN!!! ignoring serialization errors like a chad
	payload_size := uint16(len(payload))      // Mudar o tamanho do int se necessário, improvável af tho (*)
	serialized_payload_size := make([]byte, unsafe.Sizeof(payload_size))
	binary.NativeEndian.PutUint16(serialized_payload_size, payload_size) // (*) aqui também
	// Return a byte slice containing the serialized header
	return append(append([]byte{tag}, serialized_payload_size...), payload...), nil
}
