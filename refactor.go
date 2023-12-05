/*
	Karl Ramberg
	Loupe v0.1.0
	refactor.go
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
)

func refactor(dir string) error {
	fmt.Println("Loupe", loupeVersion, "-", "Refactor")

	if dir == "" {
		return errors.New("provide an archive directory using the -a flag")
	}

	stats, err := os.Stat(dir)
	if os.IsNotExist(err) || !stats.IsDir() {
		return errors.New("directory \"" + dir + "\" not found")
	}

	files, err := getImageFiles(dir)
	if err != nil {
		return errors.Join(errors.New("trouble getting images files from \""+dir+"\""), err)
	}

	if len(files) == 0 {
		return errors.New("no image files found in \"" + dir + "\"")
	}

	scanner := bufio.NewScanner(os.Stdin)
	attribute, err := promptWord(scanner, "Enter a class, group, version or subversion", "none")
	for err != nil && attribute != "none" {
		fmt.Println("Invalid:", err)
		attribute, err = promptWord(scanner, "Enter a class, group, version or subversion", "none")
	}
	fmt.Printf("Searching for \"%s\"...\n", attribute)

	var classes []string
	var groups []string
	var versions []string
	var subversions []string
	for _, f := range files {
		var p Photograph
		err := p.init(filepath.Base(f))
		if err == nil {
			if p.class == attribute && !slices.Contains(classes, p.classDir()) {
				classes = append(classes, p.classDir())
			}
			if p.group == attribute && !slices.Contains(groups, p.groupDir()) {
				groups = append(groups, p.groupDir())
			}
			if p.version == attribute && !slices.Contains(versions, p.versionDir()) {
				versions = append(versions, p.versionDir())
			}
			if p.subversion == attribute && !slices.Contains(subversions, p.subversionDir()) {
				subversions = append(subversions, p.subversionDir())
			}
		}
	}

	output := "Found attributes:\n"
	count := 1
	for _, c := range classes {
		output += strconv.Itoa(count) + ". " + c + " \t(class)\n"
		count++
	}
	for _, g := range groups {
		output += strconv.Itoa(count) + ". " + g + " \t(group)\n"
		count++
	}
	for _, v := range versions {
		output += strconv.Itoa(count) + ". " + v + " \t(version)\n"
		count++
	}
	for _, s := range subversions {
		output += strconv.Itoa(count) + ". " + s + " \t(subversion)\n"
		count++
	}
	fmt.Print(output)

	index, err := promptList(scanner)
	for err != nil {
		fmt.Println("Invalid:", err)
		index, err = promptList(scanner)
	}
	fmt.Println("Selected", index)

	// other shit

	return nil
}

func promptList(scanner *bufio.Scanner) (int, error) {
	input, err := promptInput(scanner, "Which did you mean?", "none") // TODO prompt needs work
	if err != nil {
		return 0, err
	}

	index, err := strconv.Atoi(input)
	if err != nil {
		return 0, errors.New("only use digits select an item")
	}

	return index, nil
}
