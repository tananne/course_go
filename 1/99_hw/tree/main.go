package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func printDir(out io.Writer, path string, printFiles bool, prefix string) error {
	docs, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	baseSep := "├───"
	dirSep := "│	"
	lastSep := "└───"
	fillSep := "	"

	lastDir := -1
	lastFile := -1

	for i := 0; i < len(docs); i++ {
		if !docs[i].IsDir() {
			lastFile = i
		} else {
			lastDir = i
		}
	}

	for i := 0; i < len(docs); i++ {
		if !docs[i].IsDir() {
			if printFiles {
				fileInfo, err := docs[i].Info()
				if err != nil {
					return err
				}

				var fileSize string
				if fileInfo.Size() == 0 {
					fileSize = "empty"
				} else {
					fileSize = fmt.Sprintf("%d", fileInfo.Size()) + "b"
				}
				if lastFile == i && lastDir < lastFile {
					fmt.Fprintln(out, prefix+lastSep+docs[i].Name()+" ("+fileSize+")")
				} else {
					fmt.Fprintln(out, prefix+baseSep+docs[i].Name()+" ("+fileSize+")")
				}
			}
		} else {
			if lastDir == i && (!printFiles || printFiles && lastFile < lastDir) {
				fmt.Fprintln(out, prefix+lastSep+docs[i].Name())
				printDir(out, path+"/"+docs[i].Name(), printFiles, prefix+fillSep)
			} else {
				fmt.Fprintln(out, prefix+baseSep+docs[i].Name())
				printDir(out, path+"/"+docs[i].Name(), printFiles, prefix+dirSep)
			}
		}
	}

	return nil
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	absDirPath, err := filepath.Abs(path)

	if err != nil {
		return err
	}

	return printDir(out, absDirPath, printFiles, "")
}

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
