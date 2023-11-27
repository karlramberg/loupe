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
	"time"
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

const listCols int = 4

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
		// Every command (if not told otherwise) works on the current directory
		folder, _ := os.Getwd()

		// Print the prompt
		fmt.Printf("\nWorking in: %s/\n", filepath.Base(folder))
		fmt.Printf("> ")
		scanner.Scan()

		err := scanner.Err()
		if err != nil {
			fmt.Printf(" Error! %s", err)
		}

		input := strings.Split(scanner.Text(), " ")

		start := time.Now() // DEBUG

		switch input[0] {
		case "audit":
			fmt.Printf(" Command not implemented yet\n")
		case "rename":
			fmt.Printf(" Command not implemented yet\n")
		case "sort":
			fmt.Printf(" Command not implemented yet\n")

		/*
			import <external folder> <base folder>
			* <base folder>, when not specified it imports them into the working directory
		*/
		case "import":
			// I think this will essentially move files into the base directory and call sort
			fmt.Printf(" Command not implemented yet\n")

		/*
			list <folder>
			* <folder> argument is optional
		*/
		case "list", "ls":
			/*
				TODO
				1. format a nice table if there are *many* files, probably over 50? For now big
				folders are just printed in 4 columns

				2. Filter out non-photo files
			*/

			if len(input) <= 2 {

				// Append a folder to the path if it was passed
				if len(input) == 2 {
					stats, err := os.Stat(input[1])
					if !os.IsNotExist(err) && stats.IsDir() {
						folder = filepath.Join(folder, input[1])
					} else {
						fmt.Printf(" \"%s\" is not an available folder\n", input[1])
					}
				}

				files, _ := os.ReadDir(folder)
				if len(files) != 0 {
					var fileList string
					for i, f := range files {
						fileList += " [" + fmt.Sprintf("%3d", i+1) + "] " + f.Name() + "\t"
						if i%listCols == listCols-1 {
							fileList += "\n"
						}
					}
					fmt.Print(fileList)
				} else {
					fmt.Printf(" Folder is empty\n")
				}
			} else {
				fmt.Printf(" Too many arguments! Run \"help\"\n")
			}

		/*
			folder <folder>
		*/
		case "folder", "cd":
			if len(input) == 2 {
				stats, err := os.Stat(input[1])
				if !os.IsNotExist(err) && stats.IsDir() {
					wd, _ := os.Getwd()
					os.Chdir(filepath.Join(wd, input[1]))
				} else {
					fmt.Printf(" \"%s\" is not an available folder\n", input[1])
				}
			} else if len(input) < 2 {
				fmt.Printf(" Provide the name of a folder\n")
			} else {
				fmt.Printf(" Too many arguments! Run \"help\"\n")
			}

		/*
			help
		*/
		case "help":
			// TODO add more when other commands are done and first version is finalized
			fmt.Printf(" Available commands and their usage:\n")
		case "quit", "exit":
			fmt.Printf(" Quitting Loupe...\n")
			quit = true
		case "clear":
			clear()
		default:
			fmt.Printf(" Command not recognized, run \"help\"\n")
		}

		fmt.Println(time.Since(start)) // DEBUG
	}
}
