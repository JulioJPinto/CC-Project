package main

import (
	"bufio"
	"cc_project/client/lib"
	"cc_project/helpers"
	"cc_project/protocol/fstp"
	"fmt"
	"github.com/fatih/color"
	"log"
	"os"
	"strings"
)

var commands = map[string]func(*lib.Client, []string) helpers.StatusMessage{
	"upload": lib.UploadFile,
}


func main() {
	client := &lib.Client{}
	config := fstp.FSTPConfig{Host: "localhost", Port: "8080"}
	var err error

	client.FSTPclient, err = fstp.NewFSTPClient(config)

	if err != nil {
		color.Red(err.Error())
		return
	}
	if len(os.Args) > 1 {
		lib.MakeDirectoryAvailable(client, os.Args[1])
	}

	fdata, _ := fstp.HashFile("/home/cralos/Uni/3Ano/CC/CC-Project/client_files/test_file.txt")
	fdata.OriginatorIP = client.FSTPclient.Conn.LocalAddr().String()
	fdata.Name = "test_file.txt"
	client.FSTPclient.Request(fstp.AllFilesRequest())

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
		if status.Error() != nil {
			color.Red(status.ShowErrors())
		}
	}
}
