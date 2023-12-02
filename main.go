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

type Photograph struct {
	date       string
	letter     string
	number     string
	class      string
	group      string
	version    string
	subversion string
	extension  string
}

func (p *Photograph) init(name string) error {
	p.extension = filepath.Ext(name)
	name = strings.TrimSuffix(name, p.extension)

	tokens := strings.Split(name, "_")
	if len(tokens) != 3 {
		return errors.New("filname invalid. Use format identifier_group_version")
	}

	indentifier := strings.Split(tokens[0], "-")
	if len(indentifier) != 2 {
		return errors.New("filename invalid. Use format date-number for the identifier")
	}

	validDate, err := validDate(indentifier[0])
	if !validDate {
		return err
	}
	p.date = indentifier[0]

	validNumber, err := regexp.MatchString("^([A-Z]*[0-9]+)$", indentifier[1])
	if !validNumber || err != nil {
		return errors.New("filname invalid. Number should only capital letters and numbers")
	}
	p.number = indentifier[1]

	groups := strings.Split(tokens[1], "-")

	if len(groups) > 2 {
		return errors.New("filname invalid. Use format class-group or just group")
	}

	if len(groups) == 2 {
		validClass, err1 := validWord(groups[0])
		validGroup, err2 := validWord(groups[1])
		err = errors.Join(err1, err2)
		if !validClass || !validGroup || err != nil {
			return errors.New("filename invalid. User format class-group or just group")
		}
		p.class = groups[0]
		p.group = groups[1]
	} else {
		validGroup, err := validWord(groups[0])
		if !validGroup || err != nil {
			return errors.New("filename invalid. Group can only contain alphanumeric characters")
		}
		p.class = "none"
		p.group = groups[0]
	}

	versions := strings.Split(tokens[2], "-")

	if len(versions) > 2 {
		return errors.New("filname invalid. Use format version-subversion or just version")
	}

	if len(versions) == 2 {
		validVersion, err1 := validWord(versions[0])
		validSubversion, err2 := validWord(versions[1])
		err = errors.Join(err1, err2)
		if !validVersion || !validSubversion || err != nil {
			return errors.New("filename invalid. User format version-subversion or just version")
		}
		p.version = versions[0]
		p.subversion = versions[1]
	} else {
		validVersion, err := validWord(versions[0])
		if !validVersion || err != nil {
			return errors.New("filename invalid. Use format version-subversion or just version")
		}
		p.version = versions[0]
		p.subversion = "none"
	}

	return nil
}

func (p *Photograph) filename() (name string) {
	name += p.date
	name += "-"
	if p.letter != "none" {
		name += p.letter
	}
	name += p.number

	name += "_"

	if p.class != "none" {
		name += p.class + "-"
	}
	name += p.group

	name += "_"

	name += p.version
	if p.subversion != "none" {
		name += "-" + p.subversion
	}

	name += p.extension

	return
}

func (p *Photograph) directory() (dir string) {
	if p.class != "none" {
		dir = filepath.Join(dir, p.class+"s")
	}

	dir = filepath.Join(dir, p.group)
	dir = filepath.Join(dir, p.version+"s")

	if p.subversion != "none" {
		dir = filepath.Join(dir, p.subversion)
	}

	return
}

var rawExts = []string{".3fr", ".ari", ".arw", ".srf", "srf2", ".bay", ".braw", ".crw", ".cr2", ".cr3,", ".cap", ".iiq", ".eip", ".dcs", ".dcr", ".drf", ".k25", ".kdc", ".dng", ".erf", ".fff", ".gpr", ".jxs", ".mef", ".mdc", ".mos", ".mrw", ".nef", ".nrw", ".orf", ".pef", ".ptx", ".pxn", ".r3d", ".raf", ".raw", ".rw2", ".rwl", ".rwz", ".srw", ".tco", ".x3f"}
var imageExts = []string{".jpg", ".jpeg", ".jxl", ".jp2", ".png", ".gif", ".webp", ".heic", ".heif", ".avif", ".psd", ".tif", ".tiff", ".mov", ".mp4", ".ico", ".xcf", ".bmp"}

func getImageFiles(dir string) (files []string, err error) {
	err = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		// Ignore any folder that starts with an underscore
		if d.IsDir() && d.Name()[0] == '_' {
			return filepath.SkipDir
		}

		// Only save files that are an image
		ext := strings.ToLower(filepath.Ext(path))
		if slices.Contains(imageExts, ext) || slices.Contains(rawExts, ext) {
			files = append(files, path)
		}

		return nil
	})
	if err != nil {
		return nil, errors.New("there was trouble reading files from the directory")
	}

	return
}

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

func promptSelection(scanner *bufio.Scanner, length int) ([]int, error) {
	input, err := promptInput(scanner, "Select files", "all")
	if err != nil {
		return nil, errors.New("something failed while scanning for input")
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
			selection = append(selection, index-1)
		} else if len(digits) == 2 {
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
		if err != nil {
			return "", "", errors.New("something went wrong parsing the input")
		}
		if !valid {
			return "", "", errors.New("invalid roll letter. Only use alphabetic characters")
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
		return "", "", errors.New("invalid start number. Only use a whole number")
	}

	// %0*s pads input with 0s so the has length of pad
	number := fmt.Sprintf("%0*s", padding, input)

	return letter, number, nil
}

var monthLen = []int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}

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

// Validates a given string is a proper date in YYYYMMDD format
func validDate(date string) (bool, error) {
	if len(date) != 8 {
		return false, errors.New("invalid date. Use format YYYYMMDD")
	}

	// Check that all three parts are two digit integers
	year, err1 := strconv.Atoi(date[0:4])
	month, err2 := strconv.Atoi(date[4:6])
	day, err3 := strconv.Atoi(date[6:8])
	err := errors.Join(err1, err2, err3)
	if err != nil {
		return false, errors.New("invalid date. Only use digits")
	}

	// Check if the month is valid
	if month < 1 || month > 12 {
		return false, errors.New("invalid date. Month should be between 01 and 12")
	}

	// Check if the day is valid given the month
	if month != 2 {
		if day < 1 || day > monthLen[month-1] {
			return false, errors.New("invalid date. Check how many days are in the month")
		}
	} else {
		leapYear := (year%4 == 0) && (!(year%100 == 0) || (year%400 == 0))
		if (day > 28 && !leapYear) || day > 29 {
			return false, errors.New("invalid date. I also get confused by leap years")
		}
	}

	return true, nil
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

// Validates a given string is a lowercase alphanumeric word
func validWord(word string) (bool, error) {
	valid, err := regexp.MatchString("^([a-z0-9]+)$", word)
	if !valid || err != nil {
		return false, errors.New("invalid word. Only use alphanumeric characters")
	}
	return true, nil
}

// Prompts for a basic confirmation. True if the first character entered was a y, otherwise false
func promptConfimation(scanner *bufio.Scanner) (bool, error) {
	input, err := promptInput(scanner, "Do these changes look okay?", "no")
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

	sortCmd := flag.NewFlagSet("sort", flag.ExitOnError)
	sortDir := sortCmd.String("a", "", "Archive directory")

	fmt.Print("Loupe v0.1.0")

	if len(os.Args) < 2 {
		fmt.Println("")
		fmt.Println("Error: no subcommand provided")
		return
	}

	switch os.Args[1] {

	// Name images in Loupe's format from scratch, ignoring any previous filenames
	case "name":
		fmt.Println(" - Name")
		nameCmd.Parse(os.Args[2:])

		// Directory checking
		if *nameDir == "" {
			fmt.Println("Provide a working directory using the -w flag")
			return
		}

		stat, err := os.Stat(*nameDir)
		if os.IsNotExist(err) || !stat.IsDir() {
			fmt.Println("Error: invalid working directory")
			return
		}

		// Get and print image files
		files, err := getImageFiles(*nameDir)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		if len(files) == 0 {
			fmt.Println("No image files found in working directory")
			return
		}
		output := getFileTable(files)
		fmt.Println(output)

		// User input
		scanner := bufio.NewScanner(os.Stdin)

		selections, err := promptSelection(scanner, len(files))
		for err != nil {
			fmt.Println("Error:", err)
			selections, err = promptSelection(scanner, len(files))
		}

		// Disable autodating if any non-raw files are selected
		// TODO some sort of override because camera generated jpgs are a thing
		defaultDate := "auto"
		for _, selection := range selections {
			ext := strings.ToLower(filepath.Ext(files[selection]))
			if !slices.Contains(rawExts, ext) {
				defaultDate = "none"
				break
			}
		}

		// Create a template photo, getting values from the user
		var template Photograph

		template.date, err = promptDate(scanner, defaultDate)
		for err != nil {
			fmt.Println("Error:", err)
			template.date, err = promptDate(scanner, defaultDate)
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

		template.version, err = promptWord(scanner, "Enter version", "lolidk") // TODO
		for err != nil {
			fmt.Println("Error:", err)
			template.version, err = promptWord(scanner, "Enter version", "lolidk")
		}

		template.subversion, err = promptWord(scanner, "Enter subversion", "none")
		for err != nil {
			fmt.Println("Error:", err)
			template.subversion, err = promptWord(scanner, "Enter subversion", "none")
		}

		// Rename selected files using the new template photograph
		var newFiles []string
		dateCounter := make(map[string]int)
		for _, selection := range selections {
			photo := template

			if photo.date == "auto" {
				fileinfo, err := os.Stat(files[selection])
				if err != nil {
					fmt.Println("Error: there was as problem getting an auto date for", files[selection])
					fmt.Println(err)
				}

				/*
					NOTE: Files are auto dated using their last modification time instead of a
					creation date. This is because creation dates are inconsistent across systems
					and auto dating is only useful for digital raw files, which should never be
					modified after they are made anyways
				*/
				time := fileinfo.ModTime()
				date := strconv.Itoa(time.Year())
				date += fmt.Sprintf("%02s", strconv.Itoa(int(time.Month())))
				date += fmt.Sprintf("%02s", strconv.Itoa(time.Day()))

				photo.date = date
			}

			number := strconv.Itoa(dateCounter[photo.date] + 1)
			padding := 3
			if photo.letter != "none" {
				padding = 2
			}
			photo.number = fmt.Sprintf("%0*s", padding, number)
			dateCounter[photo.date] += 1

			photo.extension = strings.ToLower(filepath.Ext(files[selection]))

			newFiles = append(newFiles, photo.filename())
		}

		output = ""
		for i, s := range selections {
			output += "Renaming " + files[s] + " to " + newFiles[i] + "\n"
		}
		fmt.Print(output)

		okay, err := promptConfimation(scanner)
		if err != nil {
			fmt.Println("Error: ", err)
		}

		if okay {
			fmt.Println("Okay!")

			for index, selection := range selections {
				oldpath, err1 := filepath.Abs(files[selection])
				newpath, err2 := filepath.Abs(files[selection])
				err = errors.Join(err1, err2)
				if err != nil {
					fmt.Println("Error: something went wrong finding final file paths")
					fmt.Println(err)
				}

				newpath = filepath.Dir(newpath)
				newpath = filepath.Join(newpath, newFiles[index])

				err = os.Rename(oldpath, newpath)
				if err != nil {
					fmt.Println("Error: there was a problem renaming", files[selection])
					fmt.Println(err)
				} else {
					fmt.Println("Renamed", filepath.Base(oldpath), "to", filepath.Base(newpath))
				}
			}
		} else {
			fmt.Println("Abort! Abort! Abort!")
		}

	// Organize validly-named images based on their class, group, version and subversion
	// Invalidly-named images are put into the base folder
	case "sort":
		fmt.Println(" - Sort")
		sortCmd.Parse(os.Args[2:])

		if *sortDir == "" {
			fmt.Println("Provide an archive directory using the -a flag")
			return
		}

		stat, err := os.Stat(*sortDir)
		if os.IsNotExist(err) || !stat.IsDir() {
			fmt.Println("Error: invalid archive directory")
			return
		}

		files, err := getImageFiles(*sortDir)
		if err != nil {
			fmt.Println("Error:", err)
		}

		if len(files) == 0 {
			fmt.Println("No image files found in archive directory")
		}

		var validPhotos []Photograph
		var validFiles []string
		var invalidFiles []string

		for _, file := range files {
			var photo Photograph
			err := photo.init(filepath.Base(file))
			if err != nil {
				invalidFiles = append(invalidFiles, file)
			} else {
				validPhotos = append(validPhotos, photo)
				validFiles = append(validFiles, file)
			}
		}

		fmt.Println("Found", len(validPhotos), "valid photos")

		// Abort if the folder has less than 2/3rds validly named photos
		if float64(len(validPhotos)) < (0.66 * float64(len(files))) {
			fmt.Println("Less than 2/3rds of images in this directory are named correctly.")
			fmt.Println("Check that this is your archive!")
			return
		}

		// Create directories that don't exist
		var madeDirs []string
		for _, photo := range validPhotos {
			dir := filepath.Join(*sortDir, photo.directory())

			_, err := os.Stat(dir)
			if !slices.Contains(madeDirs, dir) && os.IsNotExist(err) {
				err := os.MkdirAll(dir, 0755)
				if err != nil {
					fmt.Println("Error: trouble while creating directory", dir)
					fmt.Println(err)
					return
				}
				madeDirs = append(madeDirs, dir)
				fmt.Println("Created folder", dir)
			}
		}

		// Move valid photos to their directories
		for index, oldpath := range validFiles {
			newpath := filepath.Join(*sortDir, validPhotos[index].directory(), validPhotos[index].filename())
			if oldpath != newpath {
				err = os.Rename(oldpath, newpath)
				if err != nil {
					fmt.Println("Error: something went wrong moving", oldpath)
					fmt.Println(err)
				} else {
					fmt.Println("Moved", filepath.Base(oldpath), "to", filepath.Dir(newpath))
				}
			}
		}

		// Move invalid files to the base folder
		for _, oldpath := range invalidFiles {
			newpath := filepath.Join(*sortDir, filepath.Base(oldpath))
			err = os.Rename(oldpath, newpath)
			if err != nil {
				fmt.Println("Error: something went wrong moving", oldpath)
				fmt.Println(err)
			} else {
				fmt.Println("Moved invalid photo", oldpath, "to", newpath)
			}
		}

		/*
			// Delete totally empty directories
			err = filepath.WalkDir(*sortDir, func(path string, d fs.DirEntry, err error) error {
				if d.IsDir() && path != *sortDir {
					contents, err := os.ReadDir(path)
					if err != nil {
						return err
					}

					if len(contents) == 0 {
						err = os.Remove(path)
						if err != nil {
							return err
						}
					}
				}
				return nil
			})
			if err != nil {
				fmt.Println("Error: something went wrong cleaning empty directories")
				fmt.Println(err)
			}*/

		// Put invalid photos in base directory
		fmt.Println("Found", len(invalidFiles), "invalid photos")

	default:
		fmt.Printf("Error: command \"%s\" not found\n", os.Args[1])
	}
}
