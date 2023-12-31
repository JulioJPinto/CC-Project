package main

import (
	"bufio"
	"cc_project/helpers"
	"cc_project/protocol"
	"cc_project/protocol/fstp"
	"cc_project/tracker/state_manager"
	"fmt"
	"net"
	"os"

	"github.com/fatih/color"
)

type handler struct{}

var s_manager *state_manager.StateManager
var debug bool = false

// handleRequest(FSTPrequest) FSTPresponse
func (s *handler) HandleRequest(conn net.Conn, req fstp.Request) fstp.Response {
	if debug {
		color.Blue("Handling request")
		str := fmt.Sprint("\theader: ", fstp.HeaderType(int(req.Header.Flags)), "\n\tdeserialized payload: ", req.Payload, "\n\tfrom: ", conn.RemoteAddr())
		color.Blue(str)
		fmt.Println()
	}

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
	case fstp.IHaveSegReq:
		x, ok := req.Payload.(*fstp.IHaveSegmentsReqProps)
		if !ok {
			return fstp.NewErrorResponse(fmt.Errorf("invalid payload type"))
		}
		err := s_manager.BatchRegisterFileSegments(device, *x)
		if err != nil {
			return fstp.NewErrorResponse(err)
		} else {
			return fstp.NewOkResponse()
		}
	case fstp.AllFilesReq:
		allf := (s_manager.GetAllFiles())

		return fstp.NewAllFilesResponse(allf)
	case fstp.WhoHasReq:
		if debug {
			color.Cyan("\n\nWho has request\n\n")
		}
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

var commands = map[string]func(*state_manager.StateManager, []string) helpers.StatusMessage{
	"shutdown": func(g *state_manager.StateManager, a []string) helpers.StatusMessage { return state_manager.Shutdown() },
	"state": func(g *state_manager.StateManager, a []string) helpers.StatusMessage {
		s := helpers.NewStatusMessage()
		s.AddMessage(nil, g.State.String())
		return s
	},
	"segnod": func(g *state_manager.StateManager, a []string) helpers.StatusMessage {
		s := helpers.NewStatusMessage()
		s.AddMessage(nil, g.State.SegmentsNodesString())
		return s
	},
	"files": func(g *state_manager.StateManager, a []string) helpers.StatusMessage { return g.Files() },
}

func main() {

	s_manager = state_manager.NewManager("db.json")
	// s_manager.Load()
	my_handler := handler{}
	config := fstp.Config{Host: "0.0.0.0", Port: "8080"}

	for _, arg := range os.Args {
		if arg == "-d" || arg == "--debug" {
			debug = true
		}
	}
	server := fstp.NewServer(&config, &my_handler, debug)
	go server.Run()
	reader := bufio.NewReader(os.Stdin)
	helpers.TUI[*state_manager.StateManager](reader, s_manager, commands)
}
