package main

import (
	"cc_project/protocol/fstp"
	"fmt"
	// "fmt"
	// "net"
)

// const buffer_limit = 8

type handler struct {
}

// handleRequest(FSTPrequest) FSTPresponse
func (s *handler) HandleRequest(req fstp.FSTPrequest) fstp.FSTPresponse {
	fmt.Println("bué da louco")
	resp := fstp.FSTPmessage{Payload: req.Payload}
	resp.Header = fstp.FSTPHeader{Flags: fstp.IHave}
	return fstp.FSTPresponse(resp)
}

func main() {
	my_handler := handler{}
	config := fstp.FSTPConfig{Host: "localhost", Port: "8080"}
	server := fstp.New(&config, &my_handler)
	server.Run()
}
