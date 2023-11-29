# Loupe


A gnostic command line interface for organizing an archive of photographs

## Summary

Loupe is simple set of commands you can run to organize a folder of photographs. This includes file naming, refactoring, sorting, and importing. Loupe requires no external configuration files to run. Filenames hold the metadata used to sort themselves.

Files use the basic format `identifier_group_version`. Each of the three parts can be further split up. The identifier always needs to be split into `date-number`. The group can optionally include a class in order to group groups (`class-group`). The version can optionally include a subversion to better organize things (`version-subversion`). Photographs are sorted into folders using the same order they are named: class, group, version and subversion. Photographs are not sorted by their date or number, instead that identifier is used to tie together different versions and subversions of an image.

**Two example filenames**

`20241201-007_granite_master.tif` would be the name for the master tif file of the 7th photograph taken on December 1st, 2024, grouped into granite. It would reside in the folder `photographs/granite/masters/`.

`20270630-028_trip-berlin2027_print-8x10.tif` would be the name for the print-ready file, sized to 8x10" of the 28th photograph taken on June 30th, 2027 on a trip to Berlin. It would reside in the folder `photographs/trips/berlin2023/prints/8x10/`.

**An example folder**

```
photographs/
	andesite/
		masters/
			20231125-001_andesite_master.tif
			20231125-002_andesite_master.tif
			20231125-003_andesite_master.tif
		prints/
			20231125-001_andesite_print.tif
			20231125-002_andesite_print.tif
			20231125-003_andesite_print.tif
	assignments/
		basalt/
			workprints/
				20250402-001_assignment-basalt_workprint.jpg
				20250402-002_assignment-basalt_workprint.jpg
				250402-003_assignment-basalt_workprint.jpg
			finals/
				20250402-001_assignment-basalt_final.jpg
				20250402-002_assignment-basalt_final.jpg
				20250402-003_assignment-basalt_final.jpg
	trips/
		chalk/
			masters/
				20260331-001_assignment-chalk_master.tif
				20260331-002_assignment-chalk_master.tif
				20260331-003_assignment-chalk_master.tif
			prints/
				8x10s/
					20271102-001_assignment-chalk_print-8x10.tif
					20271102-002_assignment-chalk_print-8x10.tif
					20271102-003_assignment-chalk_print-8x10.tif
				16x20s/
					20271102-001_assignment-chalk_print-16x20.tif
					20271102-002_assignment-chalk_print-16x20.tif
					20271102-003_assignment-chalk_print-16x20.tif
```

Note that every image must have a group, but that the group does not need to be in a class. Also note that files must have a version, but not always a subversion.

**Why is Loupe gnostic?**

Loupe offers only one way of organizing photographs - by group and version. The first reason for this is because offering "sane defaults" is a lot more helpful than offering a host of complex and configurable ways of organizing files. I do not want Loupe to trigger analysis paralysis in it's user. I simply want it to do a few things well and get the hell out of the way when it's not in use.

Second, I did want Loupe to use external files to work. My biggest problem with software like AdobeLightroom is the need to launch a bloated piece of software to simply browse your pictures in theway you organized them. File and folder systems will not go anywhere, will not change, and will bereadable on any system.

**Why is Loupe written in Go?**

I wrote this program out of personal necessity and I tend to write my personal projects in Go. Go isa simple language, so it is easier to read and maintain code after long periods without development.I did not want to come back to add a feature in a few years time and have to deal with complexlanguage features (Rust) or shitty cross-compilation memory safety things (C).

**Why is Loupe only for the command line?**

Loupe is written for the command line because I have failed to find a good, lightweight, stable, andcross-platform UI library for Go. I want Loupe to be easy to maintain by myself, so a dependency ona flavor-of-the-month framework doesn't appeal to me. I want Loupe to work when I'm 50 as well as itdoes today.

When planning Loupe, I also found that all the things I wanted it to do are very procedural. Thismeans it lends it's self to a very basic "ask the user questions one at a time and do things basedon the answers", rather than a complex UI where it can be easy for a user to not check a box or beoverwhelmed by options.

## Installation

## Operations

### Notes on flags

`-w` is the flag to point an operation to a working directory.

`-a` is the flag to point an operation to an archive directory.

Every operation except `help` mandates the use of a `-w` or `-a` flag. This is by design to stop braindead command typing. The user is always forced to think if they are running Loupe in a working directory with a little temporary chaos or if they are running Loupe in their organized archive directory. When sensitive data is at risk, being explicit and moving a little slower is important. 
  
### `loupe name -w`

Name is the command used to create brand new filenames for new photographs. It asks the user for every part required and and every part optional in a filename. When renaming, it ignores any previous name a file had. Ideally this command should only be used once on each photograph when it is brought into the archive.

This is the only command that can change the identifier (date and number) of a photograph. An identifier is the most sensitive part of a filename, because it ties different file versions together and ties a digital file to a (most likely) physical object such as a film negative. For these reasons, it is recommended that you never use this command in your archive directory, only in a working directoy that you will ingest later. Do not point `-w` at your archive.

### `loupe sort -a`

### `loupe stats -a`

### `loupe help`

Help will print a shorter verson of this README and a link to the full one into your console.

