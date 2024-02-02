/*
	Karl Ramberg
	Loupe v0.1.0
	refactor.go
*/

package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func refactor(dir, typeStr, old, new string) error {
	fmt.Println("Loupe", loupeVersion, "-", "Rename")

	if dir == "" {
		return errors.New("provide an archive directory using the -a flag")
	}

	stats, err := os.Stat(dir)
	if os.IsNotExist(err) || !stats.IsDir() {
		return errors.New("directory \"" + dir + "\" not found")
	}

	validType, err := validType(typeStr)
	if !validType {
		return err
	}

	validOld, err := validWord(old)
	validNew, err2 := validWord(new)
	if !validOld || !validNew {
		return errors.Join(err, err2)
	}

	files, err := getImageFiles(dir)
	if err != nil {
		return errors.Join(errors.New("trouble getting images files from \""+dir+"\""), err)
	}

	if len(files) == 0 {
		return errors.New("no image files found in \"" + dir + "\"")
	}

	var renameCount int
	for _, file := range files {
		var photograph Photograph
		err := photograph.init(filepath.Base(file))
		if err != nil {
			continue
		}

		// Save the photograph's old location before it is changed
		oldDir := photograph.directory()

		if typeStr == "class" && photograph.class == old {
			photograph.class = new
		} else if typeStr == "group" && photograph.group == old {
			photograph.group = new
		} else if typeStr == "version" && photograph.version == old {
			photograph.version = new
		} else if typeStr == "subversion" && photograph.subversion == old {
			photograph.subversion = new
		}

		// Update only the filename with the new grouping
		oldpath := file
		newpath := filepath.Join(dir, filepath.Join(oldDir, photograph.filename()))
		_, err = os.Stat(newpath)
		if os.IsNotExist(err) && oldpath != newpath {
			err = os.Rename(oldpath, newpath)
			if err != nil {
				return errors.Join(errors.New("trouble while renaming \""+oldpath+"\""), err)
			} else {
				fmt.Println("Renamed", oldpath, "to", filepath.Base(newpath))
				renameCount++
			}
		}
	}

	fmt.Printf("%d files renamed\n", renameCount)
	fmt.Println()

	// Sort the files using their new grouping
	sort(dir)

	return nil
}
