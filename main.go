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
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// type Photograph struct {
// 	year       int
// 	month      int
// 	day        int
// 	letter     string
// 	number     string
// 	class      string
// 	group      string
// 	version    string
// 	subversion string
// }

// func (p *Photograph) create(f string) {
// 	return
// }

// func (p *Photograph) filename() (f string) {
// 	return ""
// }

// var folder string = "."

func changeFolder(input *[]string) {
	if len(*input) == 2 {
		if _, err := os.Stat((*input)[1]); !os.IsNotExist(err) {
			wd, _ := os.Getwd()
			os.Chdir(filepath.Join(wd, (*input)[1]))
		} else {
			fmt.Printf("Folder does not exist\n")
		}
	} else {
		fmt.Printf("Provide the name of a folder\n")
	}
}

var clear func()

func clearLinuxMacOS() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func clearWindows() {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func main() {

	scanner := bufio.NewScanner(os.Stdin)

	// Set up clear command
	switch runtime.GOOS {
	case "linux", "darwin":
		clear = clearLinuxMacOS
	case "windows":
		clear = clearWindows
	default:
		fmt.Printf("OS not supported!\n")
	}

	clear()
	fmt.Printf("Loupe v0.1.0\n")

	quit := false
	for !quit {

		// fmt.Printf("%s", folder) // TODO print folder name
		folder, _ := os.Getwd()
		fmt.Printf("%s", filepath.Base(folder))
		fmt.Printf("> ")
		scanner.Scan()

		err := scanner.Err()
		if err != nil {
			fmt.Printf(" Error! %s", err)
		}

		input := strings.Split(scanner.Text(), " ")

		switch input[0] {
		case "audit":
			fmt.Printf(" Command not implemented yet\n")
		case "rename":
			fmt.Printf(" Command not implemented yet\n")
		case "resize":
			fmt.Printf(" Command not implemented yet\n")
		case "sort":
			fmt.Printf(" Command not implemented yet\n")
		case "import":
			fmt.Printf(" Command not implemented yet\n")
		case "list":
			/*
				TODO use columns so scrolling is minimal. You should do fancy things with getting
				the terminal width so things aren't wrapping or shite.
			*/
			files, _ := os.ReadDir(folder)
			for i, f := range files {
				fmt.Printf(" (%d) %s\n", i, f.Name())
			}
		case "folder":
			changeFolder(&input)
		case "help":
			// TODO add more when other commands are done and first version is finalized
			fmt.Printf(" Available commands and their usage:\n")
		case "quit":
			fmt.Printf(" Quitting Loupe...\n")
			quit = true
		case "clear":
			cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
			cmd.Stdout = os.Stdout
			cmd.Run()
		default:
			fmt.Printf(" Command not recognized, run \"help\"\n")
		}
	}
}
