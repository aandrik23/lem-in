package internal

import (
	"fmt"
)

var (
	red    = "\x1b[31m"
	orange = "\033[38;5;208m"
	reset  = "\x1b[37m"
)

func Log(s string, errType string) {
	switch errType {
	case "error":
		fmt.Printf(red+"[ERROR]"+reset+" %s\n", s)
	case "debug":
		if debug {
			fmt.Printf(orange+"[DEBUG]"+reset+" %s\n", s)
		}
	case "help":
		fmt.Printf(help())
	}
}

func help() string {
	help := `
	[Usage] 
	go run main.go -h							     Provides help menu
	go run main.go <filename>					Runs lem-in in normal mode
	go run main.go <filename> -v/--visualize	 Runs lem-on with visualizer`

	return help
}
