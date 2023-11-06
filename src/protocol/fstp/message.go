package fstp

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"unsafe"
)

const (
	IHaveReq     = 0b0001
	IHaveFileReq = 0b0010
	WhoHasReq    = 0b0011
	AllFilesReq  = 0b0100

	OKResp  = 0b1000
	ErrResp = 0b1001

	AllFilesResp = 0b1010
	WhoHasResp   = 0b1011
)

func HeaderType(flags int) string {
	var m = map[int]string{
		IHaveReq:     "IHaveReq",
		IHaveFileReq: "IHaveFileReq",
		WhoHasReq:    "WhoHasReq",
		AllFilesReq:  "AllFilesReq",
		OKResp:       "OKResp",
		ErrResp:      "ErrResp",
		AllFilesResp: "AllFilesResp",
		WhoHasResp:   "WhoHasResp",
	}
	ret, ok := m[flags]
	if !ok {
		return "INVALID HEADER"
	}
	return ret
}

type FSTPHeader struct {
	Flags uint8
}

// FSTPmessage implements
type FSTPmessage struct {
	Header  FSTPHeader
	Payload any //JSON Serializable
}

type FSTPRequest FSTPmessage
type FSTPresponse FSTPmessage

type ErrorResponse struct {
	Err string `json:"error"`
}

func NewErrorResponse(err error) FSTPresponse {
	return FSTPresponse{FSTPHeader{ErrResp}, &ErrorResponse{err.Error()}}
}

func NewOkResponse() FSTPresponse {
	return FSTPresponse{FSTPHeader{OKResp}, nil}
}

func NewAllFilesResponse(files map[FileHash]FileMetaData) FSTPresponse {
	return FSTPresponse{FSTPHeader{Flags: AllFilesResp}, files}
}

const FSTPHEaderSize = 5 // 5 bytes

type IHaveProps struct {
	Files []FileInfo `json:"Files"`
}

type IHaveFileReqProps FileMetaData

type WhoHasReqProps struct {
	Files []FileHash `json:"Files"`
}

type WhoHasRespProps map[FileHash]DeviceIdentifier

type AllFilesRespProps struct {
	Files map[FileHash]FileMetaData `json:"Files"`
}

func MessageType(byteArray []byte) byte {
	return byteArray[0]
}

func (message *FSTPmessage) Serialize() ([]byte, error) {
	tag := message.Header.Flags
	if message.Payload == nil {
		return []byte{tag, 0, 0, 0, 0}, nil
	}
	payload, _ := json.Marshal(message.Payload) // WARN!!! ignoring serialization errors like a chad
	payload_size := uint32(len(payload))        // Mudar o tamanho do int se necessário, improvável af tho (*)
	serialized_payload_size := make([]byte, unsafe.Sizeof(payload_size))
	binary.LittleEndian.PutUint32(serialized_payload_size, payload_size) // (*) aqui também
	ret := append(append([]byte{tag}, serialized_payload_size...), payload...)
	return ret, nil

}

var tag_struct_map = map[int]any{
	IHaveReq:     &IHaveProps{},
	IHaveFileReq: &IHaveFileReqProps{},
	WhoHasReq:    &WhoHasReqProps{},
	WhoHasResp:   &WhoHasRespProps{},
	AllFilesReq:  nil,
	AllFilesResp: &AllFilesRespProps{},
	OKResp:       nil,
	ErrResp:      &ErrorResponse{},
}

func (message *FSTPmessage) Deserialize(byteArray []byte) error {
	message.Header.Flags = MessageType(byteArray)
	var err error = nil
	var payload any // json serializable
	payload, ok := tag_struct_map[int(message.Header.Flags)]
	if !ok {
		return fmt.Errorf("invalid header type: %v", message.Header.Flags)
	}
	// err = payload.Deserialize(byteArray[FSTPHEaderSize:])
	err = json.Unmarshal(byteArray[FSTPHEaderSize:], payload)
	message.Payload = payload
	return err
}

func PayloadSize(serializedHeader []byte) uint32 {
	return binary.LittleEndian.Uint32(serializedHeader[1:FSTPHEaderSize])
}
