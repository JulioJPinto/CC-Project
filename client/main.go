package main

import (
	"cc_project/client/lib"
	"cc_project/protocol/fstp"
	"os"

	"github.com/fatih/color"
)

func main() {

	config := fstp.FSTPConfig{Host: "localhost", Port: "8080"}
	client, err := fstp.NewFSTPClient(config)
	if err != nil {
		color.Red(err.Error())
		return
	}
	if len(os.Args) > 1 {
		lib.MakeDirectoryAvailable(client, os.Args[1])
	}

	// body := fstp.IHaveProps{Files: []fstp.FileInfo{{Id: 1}}}
	// client.Request(fstp.IHaveRequest(body))
	fdata, _ := lib.HashFile("/home/cralos/Uni/3Ano/CC/CC-Project/client_files/test_file.txt")
	fdata.OriginatorIP = client.Conn.LocalAddr().String()
	fdata.Name = "test_file.txt"
	// client.Request(fstp.IHaveFileRequest(fstp.IHaveFileReqProps(*fdata)))
	client.Request(fstp.AllFilesRequest())
	// make_file_available("client_files/test_file.txt")
	// Send and receive data with the server
	for {
	}
}
