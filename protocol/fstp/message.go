package fstp

import (
	"bytes"
	"cc_project/protocol"
	"encoding/binary"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"unsafe"
)

const (
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

type AllFilesRespProps []byte // map[protocol.FileHash]protocol.FileMetaData

func NewAllFilesResponse(files AllFilesRespProps) Response {
	x, _ := json.Marshal(files)
	fmt.Println("IMA BOUTA GIVE THIS MF", string(x))
	props := AllFilesRespProps(files)
	return Response{Header{Flags: AllFilesResp}, props}
}

const HeaderSize = 5 // 5 bytes

type IHaveSegmentsReqProps struct {
	Segments []protocol.FileSegment `json:"segments"`
}

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
	payload, _ := json.Marshal(message.Payload) // WARN!!! ignoring serialization errors like a chad
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(message.Payload); err != nil {
		return nil, err
	}
	payload = buffer.Bytes()

	payload_size := uint32(len(payload)) // Mudar o tamanho do int se necessário, improvável af tho (*)
	serialized_payload_size := make([]byte, unsafe.Sizeof(payload_size))
	binary.LittleEndian.PutUint32(serialized_payload_size, payload_size) // (*) aqui também
	ret := append(append([]byte{tag}, serialized_payload_size...), payload...)
	return ret, nil

}

var tag_struct_map = map[int]any{
	IHaveFileReq: &IHaveFileReqProps{},
	WhoHasReq:    &WhoHasReqProps{},
	WhoHasResp:   &WhoHasRespProps{},
	AllFilesReq:  nil,
	AllFilesResp: &AllFilesRespProps{},
	OKResp:       nil,
	ErrResp:      &ErrorResponse{},
}

func (message *Message) Deserialize(byteArray []byte) error {
	message.Header.Flags = MessageType(byteArray)
	var err error = nil
	payload, ok := tag_struct_map[int(message.Header.Flags)]
	if !ok {
		return fmt.Errorf("invalid header type: %v", message.Header.Flags)
	}
	// err = payload.Deserialize(byteArray[FSTPHEaderSize:])
	err = json.Unmarshal(byteArray[HeaderSize:], payload)
	buffer := bytes.NewBuffer(byteArray[HeaderSize:])
	decoder := gob.NewDecoder(buffer)
	decoder.Decode(payload)
	message.Payload = payload
	return err
}

func PayloadSize(serializedHeader []byte) uint32 {
	return binary.LittleEndian.Uint32(serializedHeader[1:HeaderSize])
}
