# Loupe

Loupe is simple set of commands you can run to organize a folder of photographs. This includes file naming, refactoring, sorting, and importing. Loupe requires no external configuration files to run. Filenames hold the metadata used to organize themselves.

## Naming and folders

Loupe names a photograph with four to six pieces of information: the date shot, a number to sequence each date, an optional class, a group, a version, and an optional subversion. The date and number form an identifier that ties different versions of an image together. Each file is sorted into folders first by class (if present), then group, version, and finally subversion (if present). Dates are formatted YYYYMMDD, the number starts at 1 and is padded to 3 digits for readability.

If a file orginates from a roll of film, you can name it slightly differently. The first digit of the number is replaced with a roll letter (A-Z) and the number after then represents the frame on the roll.

### Two example filenames

`20241201-007_granite_master.tif` would be the name for the master tif file of the 7th photograph taken on December 1st, 2024, grouped into granite. It would reside in the folder `photographs/granite/masters/`.

`20270630-B28_trip-berlin_print-8x10.tif` would be the name for the print-ready file, sized to 8x10" of the 28th frame on the 2nd roll started on June 30th, 2027, on a trip to Berlin. It would reside in the folder `photographs/trips/berlin2023/prints/8x10/`.

### An example folder

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
				20250402-003_assignment-basalt_workprint.jpg
			finals/
				20250402-001_assignment-basalt_final.jpg
				20250402-002_assignment-basalt_final.jpg
				20250402-003_assignment-basalt_final.jpg
	trips/
		chert/
			masters/
				20260331-001_assignment-chalk_master.tif
				20260331-002_assignment-chalk_master.tif
				20260331-003_assignment-chalk_master.tif
			prints/
				8x10/
					20271102-001_assignment-chalk_print-8x10.tif
					20271102-002_assignment-chalk_print-8x10.tif
					20271102-003_assignment-chalk_print-8x10.tif
				16x20/
					20271102-001_assignment-chalk_print-16x20.tif
					20271102-002_assignment-chalk_print-16x20.tif
					20271102-003_assignment-chalk_print-16x20.tif
```

Loupe has no restrictions on how you name your classes, groups, versions or subversions - except that they only contain alphanumeric characters. Note that for classes and versions, Loupe will add an "s" to the folder name that the filenames don't have.

## Commands

### `loupe name -w`

Name is the command used to create shiny new filenames for new photographs. It asks the user for every part required and and every part optional in a filename. When renaming, it ignores any previous name a file had. Ideally this command should only be used once on each photograph when it is brought into the archive.

This is the only command that can change the identifier (date and number) of a photograph. An identifier is the most sensitive part of a filename, because it ties different file versions together and ties a digital file to a physical object such as a print or film negative. For these reasons, it is recommended that you never use this command in your archive directory, only in a working directoy that you will ingest later. A warning will appear if you point it to a directory with over 50 image files. Do not point `-w` at your archive.

Name has the option of dating your raw files automatically. Only raw files can be dated automatically, because their modification time should always reflect the day they were shot and should never change. Other image files do not have that guaruntee.

### `loupe sort -a`

Sort is the command used to organize the files you've spent time naming. Properly named files will be moved to their respective directories: first by their class if present, then group, version, and finally subversion if present. Files that aren't properly named will be put into the base directory for you to deal with.

Note that Sort (and any other command that calls sort after it's used) will ignore any folder that starts with an underscore. This is essential for keeping auxillary files next to your photographs if you are working on a larger, more complex project. For example, in longterm photography projects where I need to deepdive on locations, people, art or writings, I will make a _research/ folder to hold all of this. Loupe won't touch it.

Sort is meant to be run on folder with many many properlly named photographs. If a folder you attempt to sort is more than a third improperly named image files, a warning is given and the command stops. This is to avoid a mess in the base directory and protect against accidently running the command in the wrong folder. Do not point `-a` at your crusty chaotic working directory.

### `loupe stats -a`

### `loupe help`

Help will print an abridged verson of this README and a link to the full one into your console.

### Notes on flags

`-w` is the flag to point an operation to a working directory.

`-a` is the flag to point an operation to an archive directory.

Every operation except `help` mandates the use of a `-w` or `-a` flag. This is by design to stop braindead command typing. The user is always forced to think if they are running Loupe in a working directory with a little temporary chaos or if they are running Loupe in their organized archive. When sensitive data is at risk, being explicit and moving a little slower is important. 

## Installation

## Errata

### Why is Loupe so stubborn?

Loupe offers only one way of organizing photographs - by group and version. This is because offering "sane defaults" is sometimes a lot more useful than offering a host of complex and configurable ways of organizing files. I do not want Loupe to trigger analysis paralysis in it's user. I simply want it to do a few things well and get the hell out of the way when it's not in use. I didn't want Loupe to use external files to work. My biggest problem with apps like Adobe Lightroom is the need to launch a bloated piece of software to simply browse your pictures in the way you organized them. File and folder systems will not go anywhere, will not change, and will be readable on any system.

For all of these reasons, Loupe offers no way of configuring anything. Date are always formatted YYYYMMDD. Numbers are always whole numbers padded to three spaces, two if its a frame on a lettered roll of film. Filename dividers are always underscores and hyphens. Folders are always structured class/group/version/subversion, you cannot just keep giving a photograph groups or versions to put it into deeper and deeper folders.

> If you need more than 4 levels of indentation, you're fucked anyway and should fix it
> 
>  	_Linus Torvalds_, sort of

### Why is Loupe only for the command line?

Loupe is written for the command line because I have failed to find a good, lightweight, stable, and cross-platform UI library for Go. I want Loupe to be easy to maintain by myself, so a dependency on a framework or a library with spotty maintenance doesn't appeal to me. I want Loupe to work when I'm 50 as well as it does today.

When planning Loupe, I also found that all the things I wanted it to do are very procedural. It lends its self to a very basic "ask the user questions one at a time and do things based on the answers", rather than a complex UI where it can be easy for a user to not check a box or be overwhelmed by options.



