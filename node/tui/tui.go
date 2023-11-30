package tui

import (
	"bufio"
	"cc_project/helpers"
	"cc_project/node/lib"
	"fmt"
	"log"
	"strings"

	"github.com/fatih/color"
)


func TUI(reader *bufio.Reader, client *lib.Node, commands map[string]func(*lib.Node, []string) helpers.StatusMessage) {
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
