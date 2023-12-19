package fstp

import (
	"cc_project/protocol"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"unsafe"
)

const (
	IHaveFileReq = 0b0010
	WhoHasReq    = 0b0011
	AllFilesReq  = 0b0100
	IHaveSegReq  = 0b0111

	OKResp  = 0b1000
	ErrResp = 0b1001

	AllFilesResp = 0b1010
	WhoHasResp   = 0b1011
)

func HeaderType(flags int) string {
	var m = map[int]string{
		IHaveFileReq: "IHaveFileReq",
		IHaveSegReq:  "IHaveSegReq",
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

type Header struct {
	Flags uint8
}

// Message implements
type Message struct {
	Header  Header
	Payload any //JSON Serializable
}

type Request Message
type Response Message

type ErrorResponse struct {
	Err string `json:"error"`
}

func NewErrorResponse(err error) Response {
	return Response{Header{ErrResp}, &ErrorResponse{err.Error()}}
}

func NewOkResponse() Response {
	return Response{Header{OKResp}, nil}
}

type AllFilesRespProps map[protocol.FileHash]protocol.FileMetaData

func NewAllFilesResponse(files AllFilesRespProps) Response {
	props := AllFilesRespProps(files)
	return Response{Header{Flags: AllFilesResp}, props}
}

const HeaderSize = 5 // 5 bytes

type IHaveSegmentsReqProps []protocol.FileSegment

type IHaveFileReqProps protocol.FileMetaData

type WhoHasReqProps struct {
	File protocol.FileHash `json:"File"`
}

type WhoHasRespProps map[protocol.DeviceIdentifier][]protocol.FileSegment

func NewWhoHasRequest(req WhoHasReqProps) Request {
	return Request{Header{Flags: WhoHasReq}, &req}
}

func NewWhoHasResponse(ret WhoHasRespProps) Response {
	return Response{Header{Flags: WhoHasResp}, ret}
}

func MessageType(byteArray []byte) byte { return byteArray[0] }

func (message *Message) Serialize() ([]byte, error) {
	tag := message.Header.Flags
	if message.Payload == nil {
		return []byte{tag, 0, 0, 0, 0}, nil
	}
	payload, _ := json.Marshal(message.Payload)
	payload_size := uint32(len(payload))
	serialized_payload_size := make([]byte, unsafe.Sizeof(payload_size))
	binary.LittleEndian.PutUint32(serialized_payload_size, payload_size)
	ret := append(append([]byte{tag}, serialized_payload_size...), payload...)
	return ret, nil

}

func empty_payload(f int) (any, bool) {
	switch f {

	case IHaveFileReq:
		{
			return &IHaveFileReqProps{}, true
		}
	case IHaveSegReq:
		{
			return &IHaveSegmentsReqProps{}, true
		}
	case WhoHasReq:
		{
			return &WhoHasReqProps{}, true
		}
	case WhoHasResp:
		{
			return &WhoHasRespProps{}, true
		}
	case AllFilesReq:
		{
			return &struct{}{}, true
		}
	case AllFilesResp:
		{
			return &AllFilesRespProps{}, true
		}
	case OKResp:
		{
			return &struct{}{}, true
		}
	case ErrResp:
		{
			return &ErrorResponse{}, true
		}
	}
	return nil, false
}

func (message *Message) Deserialize(byteArray []byte) error {
	message.Header.Flags = MessageType(byteArray)
	var err error = nil
	payload, ok := empty_payload(int(message.Header.Flags))
	if !ok {
		return fmt.Errorf("invalid header type: %v", message.Header.Flags)
	}
	err = json.Unmarshal(byteArray[HeaderSize:], payload)
	message.Payload = payload
	return err
}

func PayloadSize(serializedHeader []byte) uint32 {
	return binary.LittleEndian.Uint32(serializedHeader[1:HeaderSize])
}
