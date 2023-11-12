package main

import (
	"bufio"
	"cc_project/node/lib"
	"cc_project/helpers"
	"cc_project/protocol/fstp"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/fatih/color"
)

var commands = map[string]func(*lib.Client, []string) helpers.StatusMessage{
	"upload": func(c *lib.Client, a []string) helpers.StatusMessage { return c.UploadFiles(a) },
	"files":  func(c *lib.Client, a []string) helpers.StatusMessage { return c.ListFiles(a) },
	"fetch":  func(c *lib.Client, a []string) helpers.StatusMessage { return c.FetchFiles(a) },
	"who":    func(c *lib.Client, a []string) helpers.StatusMessage { return c.WhoHas(a) },
	// "download", lib.Download
}

func main() {
	config := &fstp.Config{Host: "localhost", Port: "8080"}
	if len(os.Args) > 1 {
		var err error
		config, err = fstp.ParseConfig(os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
	}

	client, err := lib.NewClient(*config)
	if err != nil {
		log.Fatal(err)
	}
	if len(os.Args) > 2 {
		client.MakeDirectoryAvailable(os.Args[2])
	}

	status := client.FetchFiles(nil)
	color.Green(status.ShowMessages())
	if status.Error() != nil {
		color.Red(status.ShowErrors())
	}
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
