package main

import (
	"cc_project/protocol/fstp"
)

func main() {
	config := fstp.FSTPConfig{Host: "localhost", Port: "8080"}
	client, _ := fstp.NewFSTPClient(config)
	body := fstp.IHaveProps{Files: []fstp.FileInfo{{Id: 1}}}
	client.Request(fstp.IHaveRequest(body))
	// Send and receive data with the server
}
