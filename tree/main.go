//Package main implements a simply tree unix util.
package main

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

const (
	LASTELEM = "└───"
	VERTICAL = "│"
	ELEM     = "├───"
	TAB      = "\t"
	NEWLINE  = "\n"
)

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

// dirTree handle params.
func dirTree(out io.Writer, path string, printFiles bool) error {
	// TODO add more params and handle them
	return showTree(out, path, printFiles, "")
}

// showTree will do recursive traversal tree.
func showTree(out io.Writer, path string, printFiles bool, prefix string) error {
	files, err := ioutil.ReadDir(path)
	var lastElem int
	// find last element
	if printFiles {
		lastElem = len(files) - 1
	} else {
		for i, f := range files {
			if f.IsDir() {
				lastElem = i
			}
		}
	}
	if err == nil {
		// we go through all the files in the directory, if the file is a directory, then we call diTree again
		for i, file := range files {
			// skip files if [-f] not exist
			if !printFiles && !file.IsDir() {
				continue
			}
			printFile(out, file, prefix, lastElem == i)
			if file.IsDir() {
				if i == lastElem {
					err = showTree(out, filepath.Join(path, file.Name()), printFiles, prefix+TAB)
				} else {
					err = showTree(out, filepath.Join(path, file.Name()), printFiles, prefix+VERTICAL+TAB)
				}
			}
		}
	}
	return err
}

// printFile prints file or directory.
func printFile(out io.Writer, f os.FileInfo, prefix string, IsLastElem bool) {
	if IsLastElem {
		prefix += LASTELEM
	} else {
		prefix += ELEM
	}
	name := f.Name()
	size := f.Size()
	if f.IsDir() {
		out.Write([]byte(prefix + name + "\n"))
	} else {
		if size != 0 {
			out.Write([]byte(prefix + name + " (" + strconv.FormatInt(size, 10) + "b)" + "\n"))
		} else {
			out.Write([]byte(prefix + name + " (empty)" + "\n"))
		}
	}
}
