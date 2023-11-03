package fstp

import (
	"cc_project/helpers"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"unsafe"
)

const (
	IHaveReq     = 0b0001
	IHaveFileReq = 0b0011
	WhoHasReq    = 0b0010

	OKResp  = 0b1000
	ErrResp = 0b1001
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

type ErrorResponse struct {
	Err string `json:"error"`
}

func (e *ErrorResponse) Serialize() ([]byte, error) {
	return []byte(e.Err), nil
}

func (e *ErrorResponse) Deserialize(data []byte) error {
	e.Err = string(data)
	return nil
}

func NewErrorResponse(err error) FSTPresponse {
	return FSTPresponse{FSTPHeader{ErrResp}, &ErrorResponse{err.Error()}}
}

const FSTPHEaderSize = 5 // 5 bytes

type IHaveProps struct {
	Files []FileInfo `json:"Files"`
}

type IHaveFileProps FileMetaData

func (data *IHaveFileProps) Deserialize(bytes []byte) error {
	return json.Unmarshal(bytes, data)
}

func (data *IHaveFileProps) Serialize() ([]byte, error) {
	return json.Marshal(data)
}

func (data *IHaveProps) Deserialize(bytes []byte) error {
	return json.Unmarshal(bytes, data)
}

func (data *IHaveProps) Serialize() ([]byte, error) {
	return json.Marshal(data)
}

func (message *FSTPmessage) Serialize() ([]byte, error) {
	tag := message.Header.Flags
	if message.Payload == nil {
		return []byte{tag, 0, 0, 0, 0}, nil
	}
	payload, _ := message.Payload.Serialize() // WARN!!! ignoring serialization errors like a chad
	payload_size := uint32(len(payload))      // Mudar o tamanho do int se necessário, improvável af tho (*)
	serialized_payload_size := make([]byte, unsafe.Sizeof(payload_size))
	binary.LittleEndian.PutUint32(serialized_payload_size, payload_size) // (*) aqui também
	ret := append(append([]byte{tag}, serialized_payload_size...), payload...)
	return ret, nil

}

func MessageType(byteArray []byte) byte {
	return byteArray[0]
}

func (message *FSTPmessage) Deserialize(byteArray []byte) error {
	message.Header.Flags = MessageType(byteArray)
	var err error = nil
	var payload helpers.Serializable
	switch message.Header.Flags {
	case IHaveReq:
		payload = &IHaveProps{}
	case IHaveFileReq:
		payload = &IHaveFileProps{}
	case WhoHasReq:
		// Deserialize WhoHas request
		// var payload = WhoHasProps{}º
	default:
		return fmt.Errorf("invalid header type: %v", message.Header.Flags)
	}
	err = payload.Deserialize(byteArray[FSTPHEaderSize:])
	message.Payload = payload
	return err
}
