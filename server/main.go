package main

import (
	"cc_project/protocol/fstp"
	"cc_project/server/state_manager"
	"fmt"
	"net"
)

type handler struct{}

var database state_manager.Database

// handleRequest(FSTPrequest) FSTPresponse
func (s *handler) HandleRequest(conn net.Conn, req fstp.FSTPrequest) fstp.FSTPresponse {
	fmt.Println("handler: ", &s, "a fazer cenas com ", req.Header, " & ", req.Payload, "de", conn.RemoteAddr())
	
	switch req.Header.Flags {
	case fstp.IHave:
		database.RegisterDevice(state_manager.DeviceData{Ip:net.IP(conn.RemoteAddr().Network())})
		
	}

	
	
	
	resp := fstp.FSTPmessage{Payload: req.Payload}
	resp.Header = fstp.FSTPHeader{Flags: fstp.IHave}
	return fstp.FSTPresponse(resp)
}
func main() {
	database = state_manager.NewJSONDatabase("db.json")
	database.Connect()
	my_handler := handler{}
	config := fstp.FSTPConfig{Host: "localhost", Port: "8080"}
	server := fstp.New(&config, &my_handler)
	server.Run()
}
