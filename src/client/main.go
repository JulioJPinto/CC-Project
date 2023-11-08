package main

import (
	"bufio"
	"cc_project/client/lib"
	"cc_project/helpers"
	"cc_project/protocol/fstp"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/fatih/color"
)

var commands = map[string]func(*lib.Client, []string) helpers.StatusMessage{
	"upload": lib.UploadFile,
	"files":  lib.ListFiles,
}

func main() {
	config := fstp.Config{Host: "localhost", Port: "8080"}
	client, err := lib.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}
	client.FSTPclient, err = fstp.NewClient(config)

	if err != nil {
		color.Red(err.Error())
		return
	}
	if len(os.Args) > 1 {
		lib.MakeDirectoryAvailable(client, os.Args[1])
	}

	lib.FetchFiles(client)

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
