package main

import (
    "crypto/sha256"
    "encoding/hex"
    "flag"
    "fmt"
    "io"
    "log"
    "os"
    "path/filepath"
    "strings"
)

// Custom flag type to hold strings

type stringList []string

func (s *stringList) String() string {
    return fmt.Sprintf("%v", *s)
}

func (s *stringList) Set(value string) error {
    *s = strings.Split(value, ",")
    return nil
}

// Takes a path and generates a sha256 hash of the file at that path
func getHash(path string) string {
    myfile, err := os.Open(path)
    if err != nil {
        log.Fatal("ERROR: problem opening path:", err)
    }
    defer myfile.Close()
  
    hasher := sha256.New()

    if _, err := io.Copy(hasher, myfile); err != nil {
        log.Fatal("ERROR: problem copying file into hasher:", err)
    }
  
    return hex.EncodeToString(hasher.Sum(nil))
}

// Creates a closure over the params and returns a WalkFunc function
// that gets called for each file/directory.
func walker(files *[]string, excludedDirNames stringList) filepath.WalkFunc {
    return func(path string, info os.FileInfo, err error) error {
        if err != nil {
            log.Fatal("ERROR: rx error on the walk function err param:", err)
        }

        if info.IsDir() {

            // *nix hidden dir should be skipped completely
            // Not accounting for any Windows stuff in this app.
            if info.Name()[0:1] == "." {
                // fmt.Println("## --skipping hidden dir path: ", path, ", filename: ", info.Name())

                // skip the directory and do not descend into it
                return filepath.SkipDir
            }

            // If this is an excluded directory name, skip it and don't
            // decend any further.
            for _, dirName := range excludedDirNames {
                if dirName == info.Name() {
                    return filepath.SkipDir
                }
            }

            // fmt.Println("## --skipping dir path: ", path, ", filename: ", info.Name())

            // skip the directory file but go ahead and descend into the directory
            return nil
        }

        // skip hidden files
        if info.Name()[0:1] == "." {
            return nil
        }

        // Add this file to the found file list
        *files = append(*files, path)

        return nil
    }
}

func main() {

    // Example invocation
    // 

    var files []string

    // Needs full paths, like "/Users/chuck/Desktop"
    var searchDirPaths stringList

    // Just needs a directory name, which can show up in any path
    var excludedDirNames stringList

    flag.Usage = func() {
        fmt.Printf("Usage: find_duplicate_files -searchDirPaths <CSV of dir paths> [-excludedDirNames CSV of dir names]\n\n")
        fmt.Printf("Example: find_duplicate_files -searchDirPaths /Users/chuck/Documents,/Users/chuck/Desktop -excludedDirNames repos,node_modules\n\n")
        flag.PrintDefaults()
    }

    flag.Var(&searchDirPaths, "searchDirPaths", "A comma-seperated list of full paths to a directory. (Required)")
    flag.Var(&excludedDirNames, "excludedDirNames", "A comma-seperated list of directory names that may be found in any path.")
    flag.Parse()

    if (&searchDirPaths).String() == "[]" {
        fmt.Println("searchDirPaths is required and not provided.")
        flag.PrintDefaults()
        os.Exit(1)
    }

    // keys are the hashes, values are slices of paths
    mymap := make(map[string][]string)

    for _, searchDirPath := range searchDirPaths {
        err := filepath.Walk(searchDirPath, walker(&files, excludedDirNames))
        if err != nil {
            panic(err)
        }
    }

    for _, file := range files {
        myhash := getHash(file)
        mymap[myhash] = append(mymap[myhash], file)
        //fmt.Println(myhash,"|",file)
    }

    fmt.Println("The files in each line below appear to be duplicates based upon a sha256 hash of their contents.")
    for _, val := range mymap {
        if len(val) > 1 {
            fmt.Println(strings.Join(val[:], "|"))
        }
    }
}
