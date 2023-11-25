# Loupe

A gnostic command line interface for organizing an archive of photographs

## Summary

Essentially, Loupe is set of commands you can run to organize a directory of photographs. This
includes file renaming, basic conversions, resizing, sorting, and ingesting. Loupe requires no
external configuration files to run. Filenames themselves hold the most important metadata used to
sort themselves. File use the basic format `date_number_group_version`. Pictures are sorted into
folders first by their group and then their version.

**An Example**

```
photographs/
	alabaster/
		proofs/
			231125_001_alabaster_proof.tif
			231127_007_alabaster_proof.tif
			231201_002_alabaster_proof.tif
			...
		masters/
			231125_001_alabaster_master.tif
			231127_007_alabaster_master.tif
			231201_002_alabaster_master.tif
			...
		web/
			231125_001_alabaster_web.jpg
			231127_007_alabaster_web.jpg
			231201_002_alabaster_web.jpg
			...
	basalt/
		raws/
			241001_010_basalt_raw.RAF	
			241001_011_basalt_raw.RAF	
			241001_012_basalt_raw.RAF	
			...
		8x10prints/
			241001_010_basalt_8x10print.RAF	
			241001_011_basalt_8x10print.RAF	
			241001_012_basalt_8x10print.RAF	
			...
```

**Why is Loupe gnostic?**

Loupe offers only one way of organizing photographs - by project and then image version. The first
reason for this is because I believe offering "sane defaults" is a lot more helpful than offering
a host of complex, configurable ways of organizing files. I do not want Loupe to trigger analysis
paralysis in it's user. I simply want it to do a few things well and get the hell out of the way
when it's not in use.

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

