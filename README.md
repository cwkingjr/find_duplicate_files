# Find Duplicate Files

Just a little CLI app to learn some golang.

It scans the files in the search directory paths, excluding any hidden unix directories and any directory names listed in your exclusion list, takes a sha256 hash of each file, and outputs any files where the sha's matched on one line, inluding the full path to the files.

## To Compile

With Go installed (brew install go):
```shell
go build find_duplicate_files.go
```

That will build an executable with the name of `find_duplicate_files`.

## To Run

### Help

```shell
./find_duplicate_files --help

Usage: find_duplicate_files -searchDirPaths <CSV of dir paths> [-excludedDirNames CSV of dir names]

Example: find_duplicate_files -searchDirPaths /Users/chuck/Documents,/Users/chuck/Desktop -excludedDirNames repos,node_modules

  -excludedDirNames value
    	A comma-seperated list of directory names that may be found in any path.
  -searchDirPaths value
    	A comma-seperated list of full paths to a directory. (Required)
```

### Example

```shell
./find_duplicate_files -searchDirPaths /Users/chuck/Documents,/Users/chuck/Desktop -excludedDirNames repos,node_modules
```

It prints lines of suspected duplicate files with the full path, pipe seperated, to stdout.

```shell
The files in each line below appear to be duplicates based upon a sha256 hash of their contents.
/Users/chuck/Documents/King-Charles-MITRE-20090712.odt|/Users/chuck/Documents/fromoldmacbook/King-Charles-MITRE-20090712.odt
```
