/*
	Karl Ramberg
	Loupe v0.1.0
	main.go
*/

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type photograph struct {
	year int
	month int
	day int
	number string
	class string
	group string
	version string
	subversion string
}

func (p *photograph) filename() string {
	return ""
}

func main() {
	fmt.Printf("Loupe\n")

	scanner := bufio.NewScanner(os.Stdin)

	quit := false
	for !quit {
		fmt.Printf("> ")
		scanner.Scan()

		err := scanner.Err()
		if err != nil {
			fmt.Printf("Error! %s", err)
		}

		input := strings.Split(scanner.Text(), " ")

		switch input[0] {
		case "audit":
			fmt.Printf("Command not implemented yet\n")
		case "rename":
			fmt.Printf("Command not implemented yet\n")
		case "resize":
			fmt.Printf("Command not implemented yet\n")
		case "sort":
			fmt.Printf("Command not implemented yet\n")
		case "import":
			fmt.Printf("Command not implemented yet\n")
		case "folder":
		case "help":
			fmt.Printf("Available commands and their usage:\n")
			// TODO add more when other commands are done and first version is finalized
		case "quit":
			fmt.Printf("Quitting Loupe...\n")
			quit = true
		default:
			fmt.Printf("Command not recognized, run \"help\"\n")
		}
	}
}
