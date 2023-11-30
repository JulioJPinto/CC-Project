package main

import (
	"cc_project/protocol"
	"cc_project/protocol/fstp"
	"cc_project/tracker/state_manager"
	"fmt"
	"net"

	"github.com/fatih/color"
)

type handler struct{}

var s_manager *state_manager.StateManager

// handleRequest(FSTPrequest) FSTPresponse
func (s *handler) HandleRequest(conn net.Conn, req fstp.Request) fstp.Response {
	fmt.Println()
	color.Blue("Handling request")
	str := fmt.Sprint("\theader: ", fstp.HeaderType(int(req.Header.Flags)), "\n\tdeserialized payload: ", req.Payload, "\n\tfrom: ", conn.RemoteAddr())
	color.Blue(str)
	fmt.Println()

	device := protocol.DeviceIdentifier(conn.RemoteAddr().String())
	if !s_manager.DeviceIsRegistered(device) {
		fmt.Println("registering device: ", device)
		s_manager.RegisterDevice(protocol.Device{IP: string(device)})
	}

	switch req.Header.Flags {
	case fstp.IHaveFileReq:
		x, ok := req.Payload.(*fstp.IHaveFileReqProps)
		if ok {
			return s.HandleIHaveFileRequest(device, x)
		} else {
			s_manager.DumpToFile()
			return fstp.Response(fstp.Response{Header: fstp.Header{Flags: fstp.ErrResp}, Payload: nil})
		}
	case fstp.AllFilesReq:
		return fstp.NewAllFilesResponse(s_manager.GetAllFiles())
	case fstp.WhoHasReq:
		req, ok := req.Payload.(*fstp.WhoHasReqProps)
		var ret fstp.WhoHasRespProps = s_manager.WhoHasFile(req.File)
		if ok {
			return fstp.NewWhoHasResponse(ret)
		} else {
			s_manager.DumpToFile()
			return fstp.NewErrorResponse(state_manager.ErrInvalidPayload)
		}
	default:
		return fstp.NewErrorResponse(state_manager.ErrInvalidHeader)
	}
}

func (s *handler) HandleShutdown(conn net.Conn, err error) {
	device := protocol.DeviceIdentifier(conn.RemoteAddr().String())
	s_manager.LeaveNetwork(device)
	fmt.Println(device, "left the network")

}

func (s *handler) HandleIHaveFileRequest(device protocol.DeviceIdentifier, req *fstp.IHaveFileReqProps) fstp.Response {
	err := s_manager.RegisterFile(device, protocol.FileMetaData(*req))
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
	config := fstp.Config{Host: "localhost", Port: "8080"}
	server := fstp.NewServer(&config, &my_handler)
	server.Run()
}