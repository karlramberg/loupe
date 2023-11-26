# Loupe

A gnostic command line interface for organizing an archive of photographs

## Summary

Loupe is essentially a set of commands you can run to organize a folder of photographs. This
includes file renaming, basic conversions, resizing, sorting, and importing. Loupe requires no
external configuration files to run. Filenames themselves hold the metadata used to sort themselves.

Files use the basic format `identifier_group_version`. Each of the three parts can be further split
up. The identifier always needs to be formatted to be `date-number`. The group can optionally
include a class in order to group groups (`class-group`). The version can optionally include a 
subversion to better organize things (`version-subversion`). Photographs are sorted in the same
order they are named: class, group, version and subversion.

**Two example filenames**

`241201-007_granite_master.tif` would be the name for the master tif file of the 7th photograph
taken on December 1st, 2024, grouped into granite. It would reside in the folder
`photographs/granite/masters/`.

`270630-028_trip-berlin2027_print-8x10.tif` would be the name for the print-ready file, sized to
8x10" of the 28th photograph taken on June 30th, 2027 on a trip to Berlin. It would reside in the
folder `photographs/trips/berlin2023/prints/8x10/`.


**An example folder**

```
photographs/
	andesite/
		masters/
			231125-001_andesite_master.tif
			231125-002_andesite_master.tif
			231125-003_andesite_master.tif
		prints/
			231125-001_andesite_print.tif
			231125-002_andesite_print.tif
			231125-003_andesite_print.tif
	assignments/
		basalt/
			workprints/
				250402-001_assignment-basalt_workprint.jpg
				250402-002_assignment-basalt_workprint.jpg
				250402-003_assignment-basalt_workprint.jpg
			finals/
				250402-001_assignment-basalt_final.jpg
				250402-002_assignment-basalt_final.jpg
				250402-003_assignment-basalt_final.jpg
	trips/
		chalk/
			masters/
				260331-001_assignment-chalk_master.tif
				260331-002_assignment-chalk_master.tif
				260331-003_assignment-chalk_master.tif
			prints/
				8x10s/
					271102-001_assignment-chalk_print-8x10.tif
					271102-002_assignment-chalk_print-8x10.tif
					271102-003_assignment-chalk_print-8x10.tif
				16x20s/
					271102-001_assignment-chalk_print-16x20.tif
					271102-002_assignment-chalk_print-16x20.tif
					271102-003_assignment-chalk_print-16x20.tif
```

Note that every image must have a group, but that the group does not need to be in a class. Also
note that files must have a version, but not always a subversion.

**Why is Loupe gnostic?**

Loupe offers only one way of organizing photographs - by group and version. The first reason for
this is because offering "sane defaults" is a lot more helpful than offering a host of complex,
configurable ways of organizing files. I do not want Loupe to trigger analysis paralysis in it's
user. I simply want it to do a few things well and get the hell out of the way when it's not in
use.

Second, I did want Loupe to use external files to work. My biggest problem with software like Adobe
Lightroom is the need to launch a bloated piece of software to simply browse your pictures in the
way you organized them. File and folder systems will not go anywhere, will not change, and will be
readable on any system.

**Why is Loupe written in Go?**

I wrote this program out of personal necessity and I tend to write my personal projects in Go. Go is
a simple language, so it is easier to read and maintain code after long periods without development.
I did not want to come back to add a feature in a few years time and have to deal with complex
language features (Rust) or shitty cross-compilation memory safety things (C).

**Why is Loupe only for the command line?**

Loupe is written for the command line because I have failed to find a good, lightweight, stable, and
cross-platform UI library for Go. I want Loupe to be easy to maintain by myself, so a dependency on
a flavor-of-the-month framework doesn't appeal to me. I want Loupe to work when I'm 50 as well as it
does today.

When planning Loupe, I also found that all the things I wanted it to do are very procedural. This
means it lends it's self to a very basic "ask the user questions one at a time and do things based
on the answers", rather than a complex UI where it can be easy for a user to not check a box or be
overwhelmed by options.

## Installation

## Operations

### `audit`

### `sort`

### `rename`

### `resize`

## Additional operations

### `list`

### `help`

### `quit`

