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

/*
A Photograph is a group of strings that is formed into or from a filename. Essentially this type is
a convienient way of changing small parts of a filename without janky string manipulatione every
time. Its two methods, filename() and init() provide this construction and deconstruction.
*/
type Photograph struct {
	date       string
	letter     string
	number     string
	class      string
	group      string
	version    string
	subversion string
}

func (p *Photograph) filename() (name string) {
	// Identifier
	name += p.date
	name += "-"
	if p.letter != "none" {
		name += p.letter
	}
	name += p.number

	name += "_"

	// Group
	if p.class != "none" {
		name += p.class + "-"
	}
	name += p.group

	name += "_"

	// Version
	name += p.version
	if p.subversion != "none" {
		name += "-" + p.subversion
	}

	return
}

// func (p *Photograph) init(filename string) {}

var imageExts = []string{".raf", ".nef", ".jpg", ".jpeg", ".tif", ".tiff", ".mov"}

// Walk through the directory and it's subdirectories, creating one flat list of images
func getWorkingFiles(dir *string) (names []string, err error) {
	err = filepath.WalkDir(*dir, func(path string, d fs.DirEntry, err error) error {
		ext := strings.ToLower(filepath.Ext(path))
		if slices.Contains(imageExts, ext) {
			names = append(names, path)
		}
		return nil
	})
	if err != nil {
		return nil, errors.New("there was trouble reading files from the directory")
	}

	return
}

// Creates a nice table for printing the contents of a file
func getFileTable(files []string) (table string) {

	maxLength := 0
	for _, f := range files {
		if len(f) > 0 {
			maxLength = len(f)
		}
	}

	for i, f := range files {
		table += fmt.Sprintf(" %3d. %-*s", i+1, maxLength+1, f)
		if i%2 == 1 {
			table += "\n"
		}
	}
	return
}

/*
A helper function to prompt the user for an input. This function should never be called alone as it
does not data validation on it's own, it simply returns the typed string or a provided default value
if only a newline was entered.
*/
func promptInput(scanner *bufio.Scanner, prompt, defaultInput string) (string, error) {
	fmt.Printf("%s (default: %s) ~ ", prompt, defaultInput)

	scanner.Scan()

	if err := scanner.Err(); err != nil {
		return "", errors.New("something failed while scanning for input")
	}

	if scanner.Text() == "" {
		return defaultInput, nil
	}
	return scanner.Text(), nil
}

// Returns a slice of ints from start to end in ascending order
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

/*
A helper function for parsing a selection expression. This function removes duplicate indices,
indices outside min or max, and sorts the slice in ascending order
*/
func cleanSelection(dirty []int, min, max int) []int {
	clean := []int{}
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

func promptSelection(scanner *bufio.Scanner, maxLength int) ([]int, error) {
	input, err := promptInput(scanner, "Select files", "all")
	if err != nil {
		return nil, errors.New("something failed while scanning for input")
	}

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
		return nil, errors.New("invaild selection expression")
	}

	// Build a slice of indices based on the validated expression
	// Atoi() errors can be ignored because we filtered through a regex earlier. Sue me.
	selection := []int{}
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

/*
Prompts for a number or a letter/number combo to go with the photograph's date. Together they form
the unique indentifier for each photograph.
*/
func promptNumber(scanner *bufio.Scanner) (letter, number string, err error) {
	input := ""

	input, err = promptInput(scanner, "Enter roll letter", "none")
	if err != nil {
		return
	}
	padding := 3
	if input != "none" {
		padding = 2

		// Check that the given letter is only capital letters
		input = strings.ToUpper(input)
		valid := false
		valid, err = regexp.MatchString("^([A-Z]+)$", input)
		if !valid || err != nil {
			err = errors.New("invalid roll letter. Only use alphabetic characters")
			return
		}
	}
	letter = input

	// Grab input for the starting number
	input, err = promptInput(scanner, "Enter start number", "1")
	if err != nil {
		return
	}

	// Check that the number is only digits
	valid := false
	valid, err = regexp.MatchString("^([0-9]+)$", input)
	if !valid || err != nil {
		err = errors.New("invalid start number. Only use a whole number")
		return
	}

	// %0*s pads input with 0s so the has length of pad
	number = fmt.Sprintf("%0*s", padding, input)

	return
}

var monthLengths = []int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}

/*
Prompts for a valid date to go with the photoraph's number or letter/number combo.
Together they form the unique indentifier for each photograph.
*/
func promptDate(scanner *bufio.Scanner) (date string, err error) {
	input := ""
	input, err = promptInput(scanner, "Enter date", "auto")
	if err != nil {
		return
	}

	if input == "auto" {
		date = "auto"
		return
	}

	if len(input) != 8 {
		err = errors.New("invalid date. Use format YYYYMMDD")
		return
	}

	// Check that all three parts are two digit integers
	year, err1 := strconv.Atoi(input[0:4])
	month, err2 := strconv.Atoi(input[4:6])
	day, err3 := strconv.Atoi(input[6:8])
	err = errors.Join(err1, err2, err3)
	if err != nil {
		err = errors.New("invalid date. Use format YYYYMMDD")
		return
	}

	// Check if the month is valid
	if month < 1 || month > 12 {
		err = errors.New("month need to be between 01-12")
		return
	}

	// Check if the day is valid given the month
	if month != 2 {
		if day < 1 || day > monthLengths[month-1] {
			err = errors.New("invalid day. Check how many days the month has")
		}
	} else {
		leapYear := (year%4 == 0) && (!(year%100 == 0) || (year%400 == 0))
		if day == 29 && !leapYear {
			err = errors.New("invalid day. February only has 28 days that year")
		}
	}

	date = input
	return
}

/*
Prompts the user for a single lowercase alphanumeric word
This word can be used in a class, group, version or subversion
*/
func promptWord(scanner *bufio.Scanner, prompt, defaultWord string) (word string, err error) {
	word, err = promptInput(scanner, prompt, defaultWord)
	if err != nil {
		return "", err
	}

	word = strings.ToLower(word)
	valid, err := regexp.MatchString("^([a-z0-9]+)$", word)
	if !valid || err != nil {
		return "", errors.New("invalid input. Only use a single lowercase word")
	}

	return
}

/*
Prompts the user for a basic confirmation. Return true if the first character entered was a y, otherwise return false
*/
func promptConfimation(scanner *bufio.Scanner) (okay bool, err error) {
	input := ""
	input, err = promptInput(scanner, "Do these changes look okay?", "no")
	if err != nil {
		return
	}

	if strings.ToLower(input)[0] == 'y' {
		okay = true
	}

	return
}

func main() {
	renameCmd := flag.NewFlagSet("rename", flag.ExitOnError)
	renameDir := renameCmd.String("w", ".", "Working directory")

	fmt.Print("Loupe v0.1.0")

	if len(os.Args) < 2 {
		fmt.Println("Error: no subcommand provided")
		return
	}

	switch os.Args[1] {

	/*
		name creates filenames in Loupe's gnostic style from scratch. It ignores any filename the
		files had previously. To rename individual parts of files that already use the style, use
		the refactor command
	*/
	case "name":
		fmt.Println(" - Name")
		renameCmd.Parse(os.Args[2:])

		stat, err := os.Stat(*renameDir)
		if os.IsNotExist(err) || !stat.IsDir() {
			fmt.Println("Error: invalid working directory")
			return
		}

		files, err := getWorkingFiles(renameDir)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		if len(files) == 0 {
			fmt.Println("No image files found in working directory")
			return
		}

		output := getFileTable(files)
		fmt.Print(output, "\n")

		scanner := bufio.NewScanner(os.Stdin)

		selections, err := promptSelection(scanner, len(files))
		for err != nil {
			fmt.Println("Error:", err)
			selections, err = promptSelection(scanner, len(files))
		}
		// fmt.Println(selections) // DEBUG

		template := Photograph{}

		template.date, err = promptDate(scanner)
		for err != nil {
			fmt.Println("Error:", err)
			template.date, err = promptDate(scanner)
		}

		template.letter, template.number, err = promptNumber(scanner)
		for err != nil {
			fmt.Println("Error:", err)
			template.letter, template.number, err = promptNumber(scanner)
		}

		template.class, err = promptWord(scanner, "Enter class", "none")
		for err != nil {
			fmt.Println("Error:", err)
			template.class, err = promptWord(scanner, "Enter class", "none")
		}

		template.group, err = promptWord(scanner, "Enter group", "default")
		for err != nil {
			fmt.Println("Error:", err)
			template.group, err = promptWord(scanner, "Enter group", "default")
		}

		template.version, err = promptWord(scanner, "Enter version", "lolidk")
		for err != nil {
			fmt.Println("Error:", err)
			template.version, err = promptWord(scanner, "Enter version", "lolidk")
		}

		template.subversion, err = promptWord(scanner, "Enter subversion", "none")
		for err != nil {
			fmt.Println("Error:", err)
			template.subversion, err = promptWord(scanner, "Enter subversion", "none")
		}

		// fmt.Println(template.filename()) // DEBUG

		dateCounter := make(map[string]int)
		newFilenames := []string{}
		for _, s := range selections {
			photograph := template

			/*
				NOTE: Auto dating is only useful for digital raw files. Raw files should never be
				modified after they are created, so using their mod time to date them is fine.
				Getting an actual creation time for a file on every operating system is stupidly
				inconsistent, so it's not used here. It would only provide information for digital
				raw files anyway, so it's not a huge loss. Any file that isn't a raw file (e.g.
				negative scan or edited tif or exported jpg) will have both a creation and
				modification date *after* the day the photograph was captured.

				TODO: Require the user to enter a date manually if there are any files in the
				selection that are not raw. Do this when getting input, not here. I just thought
				this place was most appropriate for this TODO.
			*/
			if photograph.date == "auto" {
				fileinfo, err := os.Stat(files[s-1])
				if err != nil {
					fmt.Println("Error: there was as problem getting an auto date for", files[s])
					fmt.Println(err)
				}

				time := fileinfo.ModTime()
				date := strconv.Itoa(time.Year())[2:4]
				date += fmt.Sprintf("%02s", strconv.Itoa(int(time.Month())))
				date += fmt.Sprintf("%02s", strconv.Itoa(time.Day()))

				photograph.date = date
			}

			number := strconv.Itoa(dateCounter[photograph.date] + 1)
			padding := 3
			if photograph.letter != "none" {
				padding = 2
			}
			photograph.number = fmt.Sprintf("%0*s", padding, number)
			dateCounter[photograph.date] += 1

			filename := photograph.filename()

			ext := filepath.Ext(files[s-1])
			if ext == ".jpeg" {
				ext = ".jpg"
			} else if ext == ".tiff" {
				ext = ".tif"
			}
			filename += strings.ToLower(ext)

			newFilenames = append(newFilenames, filename)
		}

		output = ""
		for i, s := range selections {
			output += "Renaming " + files[s-1] + " to " + newFilenames[i] + "\n"
		}
		fmt.Print(output)

		okay, err := promptConfimation(scanner)
		if err != nil {
			fmt.Println("Error: ", err)
		}

		if okay {
			fmt.Println("Okay!")

			for i, s := range selections {
				oldpath, err1 := filepath.Abs(files[s-1])
				newpath, err2 := filepath.Abs(files[s-1])
				err = errors.Join(err1, err2)
				if err != nil {
					fmt.Println("Error: something went wrong finding final file paths")
					fmt.Println(err)
				}

				newpath = filepath.Dir(newpath)
				newpath = filepath.Join(newpath, newFilenames[i])

				err = os.Rename(oldpath, newpath)
				if err != nil {
					fmt.Println("Error: there was a problem renaming", files[s-1])
					fmt.Println(err)
				} else {
					fmt.Println("Renamed", filepath.Base(oldpath), "to", filepath.Base(newpath))
				}
			}

		} else {
			fmt.Println("Abort! Abort! Abort!")
		}

	default:
		fmt.Printf("Error: command \"%s\" not found\n", os.Args[1])
	}
}
