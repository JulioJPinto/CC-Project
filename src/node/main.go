package main

import (
	"bufio"
	"cc_project/helpers"
	"cc_project/node/lib"
	"cc_project/protocol/p2p"
	"cc_project/protocol/fstp"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/fatih/color"
)

var commands = map[string]func(*lib.Gaijo, []string) helpers.StatusMessage{
	"upload": func(g *lib.Gaijo, a []string) helpers.StatusMessage { return g.UploadFiles(a) },
	"files":  func(g *lib.Gaijo, a []string) helpers.StatusMessage { return g.ListFiles(a) },
	"fetch":  func(g *lib.Gaijo, a []string) helpers.StatusMessage { return g.FetchFiles(a) },
	"who":    func(g *lib.Gaijo, a []string) helpers.StatusMessage { return g.WhoHas(a) },
	// "download", lib.Download
}

var State *lib.State

func main() {
	State = lib.NewState()
	fstp_config := fstp.Config{Host: "localhost", Port: "8080"}
	p2p_config := p2p.Config{Host: "localhost", Port: "9090"}
	if len(os.Args) > 1 {
		var err error
		fstp_config, err = fstp.ParseConfig(os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
	}

	gajo, err := lib.NewGaijo(fstp_config, p2p_config)
	if err != nil {
		log.Fatal(err)
	}
	if len(os.Args) > 2 {
		gajo.MakeDirectoryAvailable(os.Args[2])
	}

	go lib.ListenOnUDP(p2p_config)


	status := gajo.FetchFiles(nil)
	color.Green(status.ShowMessages())
	if status.Error() != nil {
		color.Red(status.ShowErrors())
	}
	reader := bufio.NewReader(os.Stdin)

	tui(reader, gajo)

}

func tui(reader *bufio.Reader, client *lib.Gaijo) {
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
