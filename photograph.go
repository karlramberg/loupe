/*
	Karl Ramberg
	Loupe v0.1.0
	photograph.go
*/

package main

import (
	"errors"
	"path/filepath"
	"regexp"
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

// Initialize a new struct given the filename, validating each piece of data
func (p *Photograph) init(name string) error {
	// Snag the file extension off the end
	p.extension = filepath.Ext(name)
	name = strings.TrimSuffix(name, p.extension)

	// Split the name into it's main sections at the underscores
	sections := strings.Split(name, "_")
	if len(sections) != 3 {
		return errors.New("filename not formatted \"identifier_group_version\"")
	}

	// Split the identifier into it's two parts
	indentifier := strings.Split(sections[0], "-")
	if len(indentifier) != 2 {
		return errors.New("identifier not formatted \"date-number\"")
	}

	// Validate the identifier's date
	validDate, err := validDate(indentifier[0])
	if !validDate {
		return errors.Join(errors.New("date is incorrect"), err)
	}
	p.date = indentifier[0]

	// Validate the identifier's letter and/or number
	validNumber, err := regexp.MatchString("^([A-Z]*[0-9]+)$", indentifier[1])
	if !validNumber || err != nil {
		return errors.Join(errors.New("number should only use capital letters and numbers"))
	}
	p.number = indentifier[1]

	// Split the group section into it's one or more parts
	groups := strings.Split(sections[1], "-")

	if len(groups) > 2 {
		return errors.New("use format \"class-group\" or just \"group\"")
	}

	if len(groups) == 2 { // Class and group
		validClassAndGroup, err := validWord(groups[0] + groups[1])
		if !validClassAndGroup || err != nil {
			return errors.Join(errors.New("groups not formatted correctly"), err)
		}
		p.class = groups[0]
		p.group = groups[1]
	} else { // Just group
		validGroup, err := validWord(groups[0])
		if !validGroup || err != nil {
			return errors.Join(errors.New("group not formatted correctly"), err)
		}
		p.class = "none"
		p.group = groups[0]
	}

	// Split the version section into it's one or more parts
	versions := strings.Split(sections[2], "-")

	if len(versions) > 2 {
		return errors.New("use format \"version-subversion\" or just \"version\"")
	}

	// Version and subversion
	if len(versions) == 2 {
		validVersionAndSubversion, err := validWord(versions[0] + versions[1])
		if !validVersionAndSubversion || err != nil {
			return errors.Join(errors.New("versions not formatted correctly"), err)
		}
		p.version = versions[0]
		p.subversion = versions[1]
	} else {
		validVersion, err := validWord(versions[0])
		if !validVersion || err != nil {
			return errors.Join(errors.New("version not formatted correctly"), err)
		}
		p.version = versions[0]
		p.subversion = "none"
	}

	return nil
}

// Construct the photograph's filename from its data
func (p *Photograph) filename() (name string) {
	// Identifier
	name += p.date
	name += "-"
	if p.letter != "none" {
		name += p.letter
	}
	name += p.number

	name += "_"

	// Group(s)
	if p.class != "none" {
		name += p.class + "-"
	}
	name += p.group

	name += "_"

	// Version(s)
	name += p.version
	if p.subversion != "none" {
		name += "-" + p.subversion
	}

	// Extension
	name += p.extension

	return
}

// Construct the directory the photograph should live in based on its group(s) and version(s)
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

var monthLen = []int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}

// Validates a given string is a proper date in YYYYMMDD format
func validDate(date string) (bool, error) {
	if len(date) != 8 {
		return false, errors.New("use format YYYYMMDD")
	}

	// Check that all three parts are two digit integers
	year, err1 := strconv.Atoi(date[0:4])
	month, err2 := strconv.Atoi(date[4:6])
	day, err3 := strconv.Atoi(date[6:8])
	err := errors.Join(err1, err2, err3)
	if err != nil {
		return false, errors.Join(errors.New("only use digits"))
	}

	// Check if the month is valid
	if month < 1 || month > 12 {
		return false, errors.New("month should be between 01 and 12")
	}

	// Check if the day is valid given the month
	if month != 2 {
		if day < 1 || day > monthLen[month-1] {
			return false, errors.New("check how many days are in the month")
		}
	} else { // Leap year fancy math
		leapYear := (year%4 == 0) && (!(year%100 == 0) || (year%400 == 0))
		if (day > 28 && !leapYear) || day > 29 {
			return false, errors.New("leap years are confusing")
		}
	}

	return true, nil
}

// Validates a given string is a lowercase alphanumeric word
func validWord(word string) (bool, error) {
	valid, err := regexp.MatchString("^([a-z0-9]+)$", word)
	if !valid || err != nil {
		return false, errors.Join(errors.New("only use alphanumeric characters"), err)
	}
	return true, nil
}
