package main

import (
	"cc_project/protocol/fstp"
	"cc_project/server/state_manager"
	"fmt"
	"net"

	"github.com/fatih/color"
)

type handler struct{}

var s_manager *state_manager.StateManager

// handleRequest(FSTPrequest) FSTPresponse
func (s *handler) HandleRequest(conn net.Conn, req fstp.FSTPRequest) fstp.FSTPresponse {
	fmt.Println()
	color.Blue("Handling request")
	str := fmt.Sprint("\theader: ", fstp.HeaderType(int(req.Header.Flags)), "\n\tdeserialized payload: ", req.Payload, "\n\tfrom: ", conn.RemoteAddr())
	color.Blue(str)
	fmt.Println()

	device := fstp.DeviceIdentifier(conn.RemoteAddr().String())
	if !s_manager.DeviceIsRegistered(device) {
		fmt.Println("registering device: ", device)

		s_manager.RegisterDevice(fstp.Device{IP: string(device)})
	}

	switch req.Header.Flags {
	case fstp.IHaveReq:
		// s_manager.RegisterFileSegment(fstp.DeviceIdentifier(conn.RemoteAddr().(*net.TCPAddr).IP), fstp.FileSegment{FirstByte: 0, FileId: 1, Hash: "aaaa"})
	case fstp.IHaveFileReq:
		x, ok := req.Payload.(*fstp.IHaveFileReqProps)
		if ok {
			return s.HandleIHaveFileRequest(device, x)
		} else {
			s_manager.DumpToFile()
			return fstp.FSTPresponse(fstp.FSTPresponse{Header: fstp.FSTPHeader{Flags: fstp.ErrResp}, Payload: nil})
		}
	case fstp.AllFilesReq:
		return fstp.NewAllFilesResponse(s_manager.GetAllFiles())
	default:
		return fstp.NewErrorResponse(state_manager.ErrInvalidHeader)
	}
	

	// resp := fstp.FSTPmessage{Payload: req.Payload}
	// resp.Header = fstp.FSTPHeader{Flags: fstp.IHave}
	return fstp.FSTPresponse(fstp.FSTPresponse{Header: fstp.FSTPHeader{Flags: fstp.OKResp}})
}

func (s *handler) HandleIHaveFileRequest(device fstp.DeviceIdentifier, req *fstp.IHaveFileReqProps) fstp.FSTPresponse {
	err := s_manager.RegisterFile(device, fstp.FileMetaData(*req))
	s_manager.DumpToFile()

	if err != nil {
		return fstp.NewErrorResponse(err)
	}
	return fstp.NewOkResponse()
}

func main() {
	s_manager = state_manager.NewManager("db.json")
	// s_manager.Load()
	my_handler := handler{}
	config := fstp.FSTPConfig{Host: "localhost", Port: "8080"}
	server := fstp.New(&config, &my_handler)
	server.Run()
}
