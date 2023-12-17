/*
	Karl Ramberg
	Loupe v0.1.0
	sort.go
*/

package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func sort(dir string) error {
	fmt.Println("Loupe", loupeVersion, "-", "Sort")

	// Check that the -a flag was used
	if dir == "" {
		return errors.New("provide an archive directory using the -a flag")
	}

	// Check that the given directory exists
	stats, err := os.Stat(dir)
	if os.IsNotExist(err) || !stats.IsDir() {
		return errors.New("directory \"" + dir + "\" not found")
	}

	// Get a list of image files in the directory and its subdirectories
	files, err := getImageFiles(dir)
	if err != nil {
		return errors.Join(errors.New("trouble getting image files from \""+dir+"\""), err)
	}

	// Check that the directory actually has image files to sort
	if len(files) == 0 {
		return errors.New("no image files found in \"" + dir + "\"")
	}

	// Sort the files into valid and invalid slices. Instantiate a Photograph struct for the valids
	var validPhotos []Photograph
	var validFiles []string
	var invalidFiles []string
	for _, file := range files {
		var photo Photograph
		err := photo.init(filepath.Base(file))
		if err != nil {
			fmt.Printf("Found invalid file: \"%s\", %s", file, err)
			invalidFiles = append(invalidFiles, file)
		} else {
			validPhotos = append(validPhotos, photo)
			validFiles = append(validFiles, file)
		}
	}

	// Ask for a confimation if the folder has less than 2/3rds validly named photos
	if float64(len(validPhotos)) < (0.66 * float64(len(files))) {
		scanner := bufio.NewScanner(os.Stdin)
		okay, err := promptConfimation(scanner,
			"Less than 2/3rds of images in this directory are named correctly, do you wish to proceed?")
		if err != nil {
			return err
		}

		if !okay {
			fmt.Println("Aborting!")
			return nil
		}
	}

	var duplicateCount int

	// Move valids to their directories, creating them if they don't exist
	for index, photo := range validPhotos {
		// Check that the new directory exists, creating it if it doesn't
		newdir := filepath.Join(dir, photo.directory())
		_, err := os.Stat(newdir)
		if os.IsNotExist(err) {
			err := os.MkdirAll(newdir, 0755)
			if err != nil {
				return errors.Join(errors.New("trouble while creating directory \""+newdir+"\""), err)
			}
			fmt.Println("Created folder", newdir)
		}

		// Move the file with the Rename function
		oldpath := validFiles[index]
		newpath := filepath.Join(newdir, photo.filename())
		_, err = os.Stat(newpath)
		if os.IsNotExist(err) && oldpath != newpath {
			err = os.Rename(oldpath, newpath)
			if err != nil {
				return errors.Join(errors.New("trouble while moving \""+oldpath+"\""), err)
			} else {
				fmt.Println("Moved", filepath.Base(oldpath), "to", filepath.Dir(newpath))
			}
		} else if oldpath != newpath {
			invalidFiles = append(invalidFiles, validFiles[index])
			duplicateCount++
			fmt.Println("Left", filepath.Base(oldpath), "alone, file already exists at the destination")
		}
	}

	// Move invalids to the base folder
	for _, oldpath := range invalidFiles {
		newpath := filepath.Join(dir, filepath.Base(oldpath))
		if oldpath != newpath {
			err = os.Rename(oldpath, newpath)
			if err != nil {
				return errors.Join(errors.New("trouble while moving \""+oldpath+"\""), err)
			} else {
				fmt.Println("Moved invalid photo", filepath.Base(oldpath), "to", dir)
			}
		}
	}

	// Clean empty directories
	_, err = cleanEmptyDirs(dir)
	if err != nil {
		return err
	}

	fmt.Println(len(validPhotos)-duplicateCount, "sorted photograph(s)")
	fmt.Println(len(invalidFiles), "photograph(s) to be fixed")

	return nil
}

// Traverses a directory, recursively removing any empty subdirectories
func cleanEmptyDirs(dir string) (bool, error) {
	// Get a list of contents in the directory
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false, errors.Join(errors.New("trouble while reading \""+dir+"\""), err)
	}

	// Go through the entries, recursing on any other directory
	removedEntries := 0
	for _, entry := range entries {
		if entry.IsDir() {
			removed, err := cleanEmptyDirs(filepath.Join(dir, entry.Name()))
			if err != nil {
				return false, err
			}
			if removed {
				removedEntries++
			}
		}
	}

	// If the directory is originally empty or all entries have been deleted
	if len(entries)-removedEntries <= 0 {
		err := os.Remove(dir)
		if err != nil {
			return false, errors.Join(errors.New("trouble while deleting \""+dir+"\""), err)
		}
		fmt.Println("Removed empty directory", dir)
		return true, nil
	}

	return false, nil
}
