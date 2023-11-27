package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"
)

var imageExts = []string{".raf", ".nef", ".jpg", ".tif", ".mov"}

func getWorkingFiles(dir *string) (names []string) {
	// TODO err for whole walk
	filepath.WalkDir(*dir, func(path string, d fs.DirEntry, err error) error {
		ext := strings.ToLower(filepath.Ext(path))
		if slices.Contains(imageExts, ext) {
			names = append(names, path)
		}
		return nil // todo fix this shite
	})

	return names
}

func getInput(scanner *bufio.Scanner, prompt string, defaultInput string) string {
	fmt.Printf("\n%s (default: %s)\n>", prompt, defaultInput)

	scanner.Scan()

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error while scanning! %s", err)
	}

	if scanner.Text() == "\n" {
		return defaultInput
	}
	return scanner.Text()
}

func main() {
	start := time.Now()

	renameCmd := flag.NewFlagSet("rename", flag.ExitOnError)
	renameDir := renameCmd.String("w", ".", "Working directory")

	if len(os.Args) < 2 {
		fmt.Println("Provide a subcommand")
		return
	}

	switch os.Args[1] {
	case "rename":
		renameCmd.Parse(os.Args[2:])

		stat, err := os.Stat(*renameDir)
		if os.IsNotExist(err) || !stat.IsDir() {
			fmt.Println("Invalid working directory")
			return
		}

		files := getWorkingFiles(renameDir)

		if len(files) == 0 {
			fmt.Printf("No image files found in working directory\n")
			return
		}

		var output string
		for i, f := range files {
			output += "[" + fmt.Sprintf("%d", i+1) + "] " + f + "\n"
		}
		fmt.Print(output)

		s := bufio.NewScanner(os.Stdin)

		input := getInput(s, "Select files", "all")
		input = getInput(s, "Enter date", "auto")
		input = getInput(s, "Enter start number", "1")
		input = getInput(s, "Roll film?", "no")
		input = getInput(s, "Enter class", "none")
		input = getInput(s, "Enter group", "default")
		input = getInput(s, "Enter version", "fuck")
		input = getInput(s, "Enter subversion", "none")
		fmt.Print(input)

	default:
		fmt.Printf("Command \"%s\" not found\n", os.Args[1])
		return
	}

	fmt.Println("Time elapsed:", time.Since(start))
}
