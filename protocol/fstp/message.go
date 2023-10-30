package fstp

import (
	"cc_project/helpers"
	"encoding/binary"
	"encoding/json"
	"unsafe"
)

const (
	IHave  = 0b0001
	WhoHas = 0b0010
)

type FSTPHeader struct {
	Flags uint8
}

// FSTPmessage implements
type FSTPmessage struct {
	Header  FSTPHeader
	Payload helpers.Serializable
}

type FSTPrequest FSTPmessage
type FSTPresponse FSTPmessage

const FSTPHEaderSize = 5 // 5 bytes

type IHaveProps struct {
	Files []FileInfo `json:"Files"`
}

func (data *IHaveProps) Deserialize(bytes []byte) error {
	return json.Unmarshal(bytes, data)
}

func (data *IHaveProps) Serialize() ([]byte, error) {
	return json.Marshal(data)
}


func (message *FSTPmessage) Serialize() ([]byte, error) {
	tag := message.Header.Flags
	payload, _ := message.Payload.Serialize() // WARN!!! ignoring serialization errors like a chad
	payload_size := uint32(len(payload))      // Mudar o tamanho do int se necessário, improvável af tho (*)
	serialized_payload_size := make([]byte, unsafe.Sizeof(payload_size))
	binary.LittleEndian.PutUint32(serialized_payload_size, payload_size) // (*) aqui também
	ret := append(append([]byte{tag}, serialized_payload_size...), payload...)
	return ret, nil

}

func MessageType(byteArray []byte) uint8 {
	return byteArray[0]
}

func (message *FSTPmessage) Deserialize(byteArray []byte) error {
	message.Header.Flags = byteArray[0]
	var err error = nil
	switch message.Header.Flags {

	case IHave:
		var payload = IHaveProps{}
		json.Unmarshal(byteArray[FSTPHEaderSize:], &payload)
		message.Payload = &payload
	case WhoHas:
		// Deserialize WhoHas request

	}
	return err
}
