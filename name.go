/*
	Karl Ramberg
	Loupe v0.1.0
	name.go
*/

package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
)

func name(dir string) error {
	fmt.Println("Loupe", loupeVersion, "-", "Name")

	// Check that the -w flag was used
	if dir == "" {
		return errors.New("provide a working directory using the -w flag")
	}

	// Check that the given directory exists
	stat, err := os.Stat(dir)
	if os.IsNotExist(err) || !stat.IsDir() {
		return errors.New("directory \"" + dir + "\" not found")
	}

	// Get a list of image files in the directory and its subdirectories
	files, err := getImageFiles(dir)
	if err != nil {
		return err
	}

	// Check that directory actually has images we can name
	if len(files) == 0 {
		return errors.New("no image files found in \"" + dir + "\"")
	}

	// Print a table of the image files
	fmt.Println(getFileTable(files))

	// Setup a scanner to standard input for the user to give values
	scanner := bufio.NewScanner(os.Stdin)

	// Ask the user for a selection of all or some of the files
	selections, err := promptSelection(scanner, len(files))
	for err != nil {
		fmt.Println("Error:", err)
		selections, err = promptSelection(scanner, len(files))
	}

	// Disable autodating by default if any non-raw files are selected
	defaultDate := "auto"
	for _, selection := range selections {
		extension := strings.ToLower(filepath.Ext(files[selection]))
		if !slices.Contains(rawExtensions, extension) {
			defaultDate = "none"
			break
		}
	}

	// Create a template photo to store the values we get from the use
	var template Photograph

	// Ask the user for a date string, format YYYYMMDD
	template.date, err = promptDate(scanner, defaultDate)
	for err != nil {
		fmt.Println("Invalid:", err)
		template.date, err = promptDate(scanner, defaultDate)
	}

	// Ask the use for an optional roll letter, "none" is the default value
	template.letter, template.number, err = promptNumber(scanner)
	for err != nil {
		fmt.Println("Invalid:", err)
		template.letter, template.number, err = promptNumber(scanner)
	}

	// Ask the user for an optional class, "none" is the default value
	template.class, err = promptWord(scanner, "Enter class", "none")
	for err != nil {
		fmt.Println("Invalid:", err)
		template.class, err = promptWord(scanner, "Enter class", "none")
	}

	// Ask the user for a group, the "default" group is the default value
	template.group, err = promptWord(scanner, "Enter group", "default")
	for err != nil {
		fmt.Println("Invalid:", err)
		template.group, err = promptWord(scanner, "Enter group", "default")
	}

	// Ask the user for a version, "lolidk" is the default version
	template.version, err = promptWord(scanner, "Enter version", "lolidk") // TODO
	for err != nil {
		fmt.Println("Invalid:", err)
		template.version, err = promptWord(scanner, "Enter version", "lolidk")
	}

	// Ask the user for an optional subversion, "none" is the default value
	template.subversion, err = promptWord(scanner, "Enter subversion", "none")
	for err != nil {
		fmt.Println("Invalid:", err)
		template.subversion, err = promptWord(scanner, "Enter subversion", "none")
	}

	// Setup and save the new filenames starting with the values from the template photograph
	var newFilenames []string
	dateCounter := make(map[string]int)
	checklist := ""
	for _, selection := range selections {
		photo := template

		/*
			NOTE: Files are auto dated using their last modification time instead of a
			creation date. This is because creation dates are inconsistent across systems
			and auto dating is only useful for digital raw files, which should never be
			modified after they are made anyways
		*/
		if photo.date == "auto" {
			stats, err := os.Stat(files[selection])
			if err != nil {
				fmt.Println("Error: there was as problem getting an auto date for", files[selection])
				fmt.Println(err)

			}

			time := stats.ModTime()
			date := strconv.Itoa(time.Year())
			date += fmt.Sprintf("%02s", strconv.Itoa(int(time.Month())))
			date += fmt.Sprintf("%02s", strconv.Itoa(time.Day()))

			photo.date = date
		}

		// Add the number, based on the given start number and the number of times we've seen
		// photographs with the same date before in this loop
		number := strconv.Itoa(dateCounter[photo.date] + 1)
		padding := 3
		if photo.letter != "none" {
			padding = 2
		}
		photo.number = fmt.Sprintf("%0*s", padding, number)
		dateCounter[photo.date]++

		// Add the extension from the original filename
		photo.extension = strings.ToLower(filepath.Ext(files[selection]))

		// Covert the photograph struct to it's filename and add it to the list for later
		// Also add it to a checklist of changes to print directly after the loop
		filename := photo.filename()
		newFilenames = append(newFilenames, filename)
		checklist += "Renaming " + files[selection] + " to " + photo.filename() + "\n"
	}

	// Ask the user for a final confirmation of the changes
	fmt.Print(checklist)
	okay, err := promptConfimation(scanner, "Do these changes look okay?")
	if err != nil {
		return err
	}

	if okay {
		fmt.Println("Okay!")

		for index, selection := range selections {
			// Get the new path for the renamed file by replacing the filename in the old path
			oldpath := files[selection]
			newpath := filepath.Join(filepath.Dir(oldpath), newFilenames[index])

			// Rename the file!
			err = os.Rename(oldpath, newpath)
			if err != nil {
				return errors.Join(errors.New("there was a problem renaming \""+filepath.Base(oldpath)+"\""), err)
			} else {
				fmt.Println("Renamed", filepath.Base(oldpath), "to", filepath.Base(newpath))
			}
		}
	} else {
		fmt.Println("Aborting!")
	}

	return nil
}
