package main

import (
	"cc_project/protocol/fstp"
	"cc_project/server/state_manager"
	"fmt"
	"net"
)

type handler struct{}

var s_manager *state_manager.StateManager

// handleRequest(FSTPrequest) FSTPresponse
func (s *handler) HandleRequest(conn net.Conn, req fstp.FSTPrequest) fstp.FSTPresponse {
	fmt.Println("handler: ", &s, "a fazer cenas com ", req.Header, " & ", req.Payload, "de", conn.RemoteAddr())
	switch req.Header.Flags {
	case fstp.IHaveReq:
		// s_manager.RegisterFileSegment(fstp.DeviceIdentifier(conn.RemoteAddr().(*net.TCPAddr).IP), fstp.FileSegment{FirstByte: 0, FileId: 1, Hash: "aaaa"})
	case fstp.IHaveFileReq:
		x, ok := req.Payload.(*fstp.IHaveFileProps)
		if ok {
			err := s_manager.RegisterFile(fstp.DeviceIdentifier(conn.RemoteAddr().String()), fstp.FileMetaData(*x))
			if err != nil {
				return fstp.NewErrorResponse(err)
			}
			s_manager.DumpToFile()
		} else {
			return fstp.FSTPresponse(fstp.FSTPresponse{Header: fstp.FSTPHeader{Flags: fstp.ErrResp}, Payload: nil})
		}

	}
	// resp := fstp.FSTPmessage{Payload: req.Payload}
	// resp.Header = fstp.FSTPHeader{Flags: fstp.IHave}
	return fstp.FSTPresponse(fstp.FSTPresponse{Header: fstp.FSTPHeader{Flags: fstp.OKResp}})
}
func main() {
	s_manager = state_manager.NewManager("db.json")
	// s_manager.Load()
	my_handler := handler{}
	config := fstp.FSTPConfig{Host: "localhost", Port: "8080"}
	server := fstp.New(&config, &my_handler)
	server.Run()
}
