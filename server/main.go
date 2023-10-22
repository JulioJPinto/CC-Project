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
	return fstp.FSTPresponse{}
}

func main() {
	my_handler := handler{}
	config := fstp.FSTPConfig{Host: "localhost", Port: "8080"}
	server := fstp.New(&config, &my_handler)
	server.Run()
}
