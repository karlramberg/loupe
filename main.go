/*
	Karl Ramberg
	Loupe v0.1.0
	main.go
*/

package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

var imageExts = []string{".raf", ".nef", ".jpg", ".tif", ".mov"}

func getWorkingFiles(dir *string) ([]string, error) {
	var names []string

	// Walk through the directory and it's subdirectories, creating one flat list of image files
	err := filepath.WalkDir(*dir, func(path string, d fs.DirEntry, err error) error {
		ext := strings.ToLower(filepath.Ext(path))
		if slices.Contains(imageExts, ext) {
			names = append(names, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return names, nil
}

func getSelections(input string, maxLength int) ([]int, error) {

	if input == "all" {
		return makeRange(1, maxLength), nil
	}

	/*
		Check that the given phrase matches a valid selection expression
		A valid selection expression is 1 or more numbers or ranges seperated by commas.
		e.g. "1-5, 14, 22, 117-114" would be a selection of [1 2 3 4 5 14 22 114 115 116 117]
	*/
	valid, err := regexp.MatchString("^(([0-9]+([-][0-9]+)?)([,]([0-9]+([-][0-9]+)?))*)$", input)
	if !valid || err != nil {
		return nil, errors.New("invaild selection phrase")
	}

	// Build a selection slice of indices based on the validated pharse
	var selection []int
	tokens := strings.Split(input, ",")
	for _, token := range tokens {
		digits := strings.Split(token, "-")
		if len(digits) == 1 {
			index, _ := strconv.Atoi(digits[0])
			selection = append(selection, index)
		} else if len(digits) == 2 {
			start, _ := strconv.Atoi(digits[0])
			end, _ := strconv.Atoi(digits[1])
			selection = append(selection, makeRange(start, end)...)
		}
	}

	selection = cleanSelection(selection, 1, maxLength)

	if len(selection) == 0 {
		return nil, errors.New("you somehow selected no actual images")
	}

	return selection, nil
}

func makeRange(start int, end int) []int {
	if start == end {
		return []int{start}
	}

	if start > end {
		start, end = end, start
	}

	r := make([]int, (end-start)+1)
	for i := range r {
		r[i] = i + start
	}
	return r
}

// Removes duplicate indices, extreme indices, and sorts the slice
func cleanSelection(dirty []int, min, max int) []int {
	var clean []int
	seen := make(map[int]bool)
	for _, num := range dirty {
		if !seen[num] {
			seen[num] = true
			if num >= min && num <= max {
				clean = append(clean, num)
			}
		}
	}
	slices.Sort(clean)
	return clean
}

func getInput(scanner *bufio.Scanner, prompt string, defaultInput string) string {
	fmt.Printf("\n%s (default: %s)\n>", prompt, defaultInput)

	scanner.Scan()

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error while scanning! %s", err)
	}

	if scanner.Text() == "" {
		return defaultInput
	}
	return scanner.Text()
}

func main() {
	renameCmd := flag.NewFlagSet("rename", flag.ExitOnError)
	renameDir := renameCmd.String("w", ".", "Working directory")

	if len(os.Args) < 2 {
		fmt.Println("Error: no subcommand provided")
		return
	}

	switch os.Args[1] {
	case "rename":
		renameCmd.Parse(os.Args[2:])

		stat, err := os.Stat(*renameDir)
		if os.IsNotExist(err) || !stat.IsDir() {
			fmt.Println("Error: invalid working directory")
			return
		}

		files, err := getWorkingFiles(renameDir)
		if err != nil {
			fmt.Println("Error: there was trouble reading files from the directory")
			fmt.Println(err)
			return
		}

		if len(files) == 0 {
			fmt.Println("Error: no image files found in working directory")
			return
		}

		var output string
		for i, f := range files {
			output += "[" + fmt.Sprintf("%d", i+1) + "] " + f + "\n"
		}
		fmt.Print(output)

		s := bufio.NewScanner(os.Stdin)

		input := getInput(s, "Select files", "all")
		selections, err := getSelections(input, len(files))
		for err != nil {
			fmt.Println("Error:", err)

			input = getInput(s, "Select files", "all")
			selections, err = getSelections(input, len(files))
		}
		fmt.Println(selections)

		preview := getInput(s, "Enter date", "auto")
		preview += "-"

		letter := getInput(s, "Enter roll letter", "none")

		padding := 3
		if letter != "none" {
			preview += letter
			padding = 2
		}
		preview += fmt.Sprintf("%0*s", padding, getInput(s, "Enter start number", "1"))

		preview += "_"

		class := getInput(s, "Enter class", "none")
		if class != "none" {
			preview += class + "-"
		}

		preview += getInput(s, "Enter group", "default")

		preview += "_"

		preview += getInput(s, "Enter version", "lolidk")

		subversion := getInput(s, "Enter subversion", "none")
		if subversion != "none" {
			preview += "-" + subversion
		}

		fmt.Println(preview)

	default:
		fmt.Printf("Error: command \"%s\" not found\n", os.Args[1])
	}
}
