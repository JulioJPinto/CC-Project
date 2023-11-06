package main

import (
	"bufio"
	actions "cc_project/client/lib"
	"cc_project/protocol/fstp"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/fatih/color"
)

var commands = map[string]func(*fstp.FSTPclient, []string) actions.StatusMessage{
	"upload": actions.UploadFile,
}

func main() {

	config := fstp.FSTPConfig{Host: "localhost", Port: "8080"}
	client, err := fstp.NewFSTPClient(config)
	if err != nil {
		color.Red(err.Error())
		return
	}
	if len(os.Args) > 1 {
		actions.MakeDirectoryAvailable(client, os.Args[1])
	}

	// body := fstp.IHaveProps{Files: []fstp.FileInfo{{Id: 1}}}
	// client.Request(fstp.IHaveRequest(body))
	fdata, _ := fstp.HashFile("/home/cralos/Uni/3Ano/CC/CC-Project/client_files/test_file.txt")
	fdata.OriginatorIP = client.Conn.LocalAddr().String()
	fdata.Name = "test_file.txt"
	// client.Request(fstp.IHaveFileRequest(fstp.IHaveFileReqProps(*fdata)))
	client.Request(fstp.AllFilesRequest())
	// make_file_available("client_files/test_file.txt")
	// Send and receive data with the server
	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		line = strings.TrimRight(line, "\n")
		split_line := strings.Split(line, " ")
		command := split_line[0]
		f, ok := commands[command]
		if !ok {
			color.Red(fmt.Sprint("Invalid command: ", command))
			break
		}
		status := f(client, split_line[1:])
		color.Green(status.ShowMessages())
		color.Red(status.ShowErrors())
	}
}
