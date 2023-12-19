package main

import (
	"bufio"
	"cc_project/helpers"
	"cc_project/node/lib"
	"cc_project/protocol/fstp"
	"cc_project/protocol/p2p"
	"fmt"
	"log"
	"os"

	"github.com/fatih/color"
)

var commands = map[string]func(*lib.Node, []string) helpers.StatusMessage{
	"upload":   func(g *lib.Node, a []string) helpers.StatusMessage { return g.UploadFiles(a) },
	"files":    func(g *lib.Node, a []string) helpers.StatusMessage { return g.ListFiles(a) },
	"fetch":    func(g *lib.Node, a []string) helpers.StatusMessage { return g.FetchFiles(a) },
	"who":      func(g *lib.Node, a []string) helpers.StatusMessage { return g.WhoHas(a) },
	"download": func(g *lib.Node, a []string) helpers.StatusMessage { return g.Download(a) },
	"d":        func(g *lib.Node, a []string) helpers.StatusMessage { return g.Download(a) },
	"abort":    func(g *lib.Node, a []string) helpers.StatusMessage { return g.AbortDownload(a) },
	"d_state":  func(g *lib.Node, a []string) helpers.StatusMessage { return g.DownloadState(a) },
	"ongoing":  func(g *lib.Node, a []string) helpers.StatusMessage { return g.OngoingDownloads(a) },

	"test":  func(g *lib.Node, a []string) helpers.StatusMessage { return g.Test(a) },
	"peers": func(g *lib.Node, a []string) helpers.StatusMessage { return g.Stats() },
	// "status":
	"leave": func(g *lib.Node, a []string) helpers.StatusMessage { os.Exit(0); return helpers.StatusMessage{} },
}

func main() {
	fstp_config := fstp.Config{Host: "localhost", Port: "8080"}
	p2p_config := p2p.Config{Host: "localhost", Port: "9090"}
	if len(os.Args) > 1 {
		var err error
		fstp_config, err = fstp.ParseConfig(os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
	}

	debugging := false
	for _, arg := range os.Args {
		if arg == "-d" || arg == "--debug" {
			debugging = true
		}
	}

	node, err := lib.NewNode(fstp_config, p2p_config, debugging)

	if err != nil {
		log.Fatal(err)
	}
	var path string
	if len(os.Args) > 2 {
		path = os.Args[2]
		err := node.MakeDirectoryAvailable(os.Args[2])
		if os.IsNotExist(err) {
			os.Mkdir(path, 0700)
			fmt.Println("created", path, "folder")
		}
	} else {
		path = "node_files"
		os.Mkdir(path, 0700)
		fmt.Println("created", path, "folder")
	}

	node.NodeDir = path
	go node.WatchFolder()
	status := node.FetchFiles(nil)
	color.Green(status.ShowMessages())
	if status.Error() != nil {
		color.Red(status.ShowErrors())
	}
	reader := bufio.NewReader(os.Stdin)

	go node.ListenOnUDP()
	helpers.TUI(reader, node, commands)

}
