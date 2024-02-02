/*
	Karl Ramberg
	Loupe v0.1.0
	main.go

	NOTE: I am leaning towards more verbose commenting here even if some of the things happening are
	quite basic. More than likely, development on this will as needed so I will probably not be
	overly familiar with the codebase when I do so. If you are reading this an are an expert on Go,
	I apologize.
*/

/*
TODO for v1
[ ] Commands for renaming attributes (refactor)
[ ] Commands for changing selected identifiers attributes (modify)
[ ] Help command
[ ] Really solid printing output
[ ] Robust error messages
[ ] Good comments for future Karl
[ ] "print" commands, prints a nice table of valid photos sorted by identifier (good to see a timeline of your work from
	start to present)
[x] Clean-up init() in particular
[x] File clean-up, possible split into multiple (cli and actual data)
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

const loupeVersion string = "v0.1.0"

var rawExtensions = []string{
	".3fr", ".ari", ".arw", ".srf", "srf2", ".bay", ".braw", ".crw", ".cr2", ".cr3,", ".cap",
	".iiq", ".eip", ".dcs", ".dcr", ".drf", ".k25", ".kdc", ".dng", ".erf", ".fff", ".gpr", ".jxs",
	".mef", ".mdc", ".mos", ".mrw", ".nef", ".nrw", ".orf", ".pef", ".ptx", ".pxn", ".r3d", ".raf",
	".raw", ".rw2", ".rwl", ".rwz", ".srw", ".tco", ".x3f",
}

var imageExtensions = []string{
	".jpg", ".jpeg", ".jxl", ".jp2", ".png", ".gif", ".webp", ".heic", ".heif", ".avif", ".psd",
	".tif", ".tiff", ".mov", ".mp4", ".ico", ".xcf", ".bmp",
}

// Walks through a directory, creating a list of image files
func getImageFiles(dir string) (files []string, err error) {
	err = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		// Ignore any folder that starts with an underscore
		if d.IsDir() && d.Name()[0] == '_' {
			return filepath.SkipDir
		}

		// Only save files that are an image
		ext := strings.ToLower(filepath.Ext(path))
		if slices.Contains(imageExtensions, ext) || slices.Contains(rawExtensions, ext) {
			files = append(files, path)
		}

		return nil
	})
	if err != nil {
		return nil, errors.Join(errors.New("there was trouble reading files from the directory"), err)
	}

	return
}

// Constructs a nice table of numbered files
func getFileTable(files []string) (table string) {
	var width int
	for _, f := range files {
		if len(f) > width {
			width = len(f)
		}
	}

	for i, f := range files {
		table += fmt.Sprintf(" %3d. %-*s", i+1, width+1, f)
		if i%2 == 1 {
			table += "\n"
		}
	}
	return
}

// Helper function for the other prompt functions
func promptInput(scanner *bufio.Scanner, prompt, defaultInput string) (string, error) {
	fmt.Printf("%s (default: %s) ~ ", prompt, defaultInput)

	scanner.Scan()

	if err := scanner.Err(); err != nil {
		return "", errors.Join(errors.New("something failed while scanning for input"), err)
	}

	if scanner.Text() == "" {
		return defaultInput, nil
	}
	return scanner.Text(), nil
}

func promptSelection(scanner *bufio.Scanner, length int) ([]int, error) {
	input, err := promptInput(scanner, "Select files", "all")
	if err != nil {
		return nil, err
	}

	if input == "all" {
		return makeRange(0, length-1), nil
	}

	/*
		Check that the given phrase matches a valid selection expression
		A valid selection expression is 1 or more numbers or ranges seperated by commas.
		e.g. "1-5, 14, 22, 117-114" would be a selection of [1 2 3 4 5 14 22 114 115 116 117]
	*/
	valid, err := regexp.MatchString("^(([0-9]+([-][0-9]+)?)([,]([0-9]+([-][0-9]+)?))*)$", input)
	if !valid || err != nil {
		return nil, errors.Join(errors.New("invaild selection expression"), err)
	}

	// Build a slice of indices based on the validated expression
	// Atoi() errors can be ignored because we filtered through a regex earlier. Sue me.
	selection := []int{}
	tokens := strings.Split(input, ",")
	for _, token := range tokens {
		digits := strings.Split(token, "-")
		if len(digits) == 1 { // Single number
			index, _ := strconv.Atoi(digits[0])
			selection = append(selection, index-1)
		} else if len(digits) == 2 { // Range
			start, _ := strconv.Atoi(digits[0])
			end, _ := strconv.Atoi(digits[1])
			selection = append(selection, makeRange(start-1, end-1)...)
		}
	}

	selection = cleanSelection(selection, 0, length-1)

	if len(selection) == 0 {
		return nil, errors.New("somehow you selected no actual images")
	}

	return selection, nil
}

// Creates a slice of ints from start to end
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

// Removes duplicate indices, indices outside min or max, and sorts in ascending order
func cleanSelection(dirty []int, min, max int) (clean []int) {
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

// Prompts for a number or a letter/number (e.g. 007 or B07)
func promptNumber(scanner *bufio.Scanner) (string, string, error) {
	input, err := promptInput(scanner, "Enter roll letter", "none")
	if err != nil {
		return "", "", err
	}

	padding := 3
	if input != "none" {
		padding = 2

		// Check that the given letter is only capital letters
		input = strings.ToUpper(input)
		valid, err := regexp.MatchString("^([A-Z]+)$", input)
		if !valid || err != nil {
			return "", "", errors.Join(errors.New("invalid roll letter. Only use capital letter"), err)
		}
	}
	letter := input

	// Grab input for the starting number
	input, err = promptInput(scanner, "Enter start number", "1")
	if err != nil {
		return "", "", err
	}

	// Check that the number is only digits
	valid, err := regexp.MatchString("^([0-9]+)$", input)
	if !valid || err != nil {
		return "", "", errors.Join(errors.New("invalid start number. Only use a whole number"), err)
	}

	// %0*s pads input with zeros so the has length of pad
	number := fmt.Sprintf("%0*s", padding, input)

	return letter, number, nil
}

// Prompts for a date from the user
func promptDate(scanner *bufio.Scanner, defaultDate string) (string, error) {
	date, err := promptInput(scanner, "Enter date", defaultDate)
	if err != nil {
		return "", nil
	}

	if date == "auto" {
		return date, nil
	}

	valid, err := validDate(date)
	if !valid {
		return "", err
	}

	return date, nil
}

// Prompts for a lowercase alphanumeric word. Used for class, group, version and subversion
func promptWord(scanner *bufio.Scanner, prompt, defaultWord string) (string, error) {
	word, err := promptInput(scanner, prompt, defaultWord)
	if err != nil {
		return "", err
	}

	word = strings.ToLower(word)

	valid, err := validWord(word)
	if !valid || err != nil {
		return "", err
	}

	return word, nil
}

// Prompts for a basic confirmation. True if the first character entered was a y, otherwise false
func promptConfimation(scanner *bufio.Scanner, message string) (bool, error) {
	input, err := promptInput(scanner, message, "no")
	if err != nil {
		return false, err
	}

	if strings.ToLower(input)[0] == 'y' {
		return true, nil
	}

	return false, nil
}

func main() {
	nameCmd := flag.NewFlagSet("name", flag.ExitOnError)
	nameDir := nameCmd.String("w", "", "Working directory")

	refactorCmd := flag.NewFlagSet("refactor", flag.ExitOnError)
	refactorDir := refactorCmd.String("a", "", "Archive directory")
	refactorType := refactorCmd.String("t", "", "Group type")
	refactorOld := refactorCmd.String("o", "", "Old group name")
	refactorNew := refactorCmd.String("n", "", "New group name")

	sortCmd := flag.NewFlagSet("sort", flag.ExitOnError)
	sortDir := sortCmd.String("a", "", "Archive directory")

	if len(os.Args) < 2 {
		fmt.Println("Loupe", loupeVersion)
		fmt.Println("Error: no subcommand provided")
		return
	}

	switch os.Args[1] {

	// Name images in Loupe's format from scratch, ignoring any previous filenames
	case "name":
		nameCmd.Parse(os.Args[2:])
		err := name(*nameDir)
		if err != nil {
			fmt.Println("Error:", err)
		}

	// Change the name of a class, group, version or subversion
	case "refactor":
		refactorCmd.Parse(os.Args[2:])
		err := refactor(*refactorDir, *refactorType, *refactorOld, *refactorNew)
		if err != nil {
			fmt.Println("Error:", err)
		}

	// Organize validly-named images based on their class, group, version and subversion
	// Invalidly-named images are put into the base folder
	case "sort":
		sortCmd.Parse(os.Args[2:])
		err := sort(*sortDir)
		if err != nil {
			fmt.Println("Error:", err)
		}

	default:
		fmt.Println("")
		fmt.Printf("Error: command \"%s\" not found\n", os.Args[1])
	}
}
