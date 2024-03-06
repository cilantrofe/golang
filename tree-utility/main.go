package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
)

func dirTree(out io.Writer, path string, printFiles bool) error {
	return printDir(out, path, printFiles, "")
}

func printDir(out io.Writer, path string, printFiles bool, prefix string) error {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	var printableFiles []os.FileInfo
	if !printFiles {
		for _, file := range files {
			if file.IsDir() {
				printableFiles = append(printableFiles, file)
			}
		}
	} else {
		printableFiles = files
	}

	for i, file := range printableFiles {
		if !file.IsDir() {
			if i == len(printableFiles)-1 {
				if file.Size() == 0 {
					fmt.Fprintf(out, "%s└───%s (%s)\n", prefix, file.Name(), "empty")
				} else {
					fmt.Fprintf(out, "%s└───%s (%db)\n", prefix, file.Name(), file.Size())
				}
			} else {
				if file.Size() == 0 {
					fmt.Fprintf(out, "%s├───%s (%s)\n", prefix, file.Name(), "empty")
				} else {
					fmt.Fprintf(out, "%s├───%s (%db)\n", prefix, file.Name(), file.Size())
				}
			}
		} else {
			if i == len(printableFiles)-1 {
				fmt.Fprintf(out, "%s└───%s\n", prefix, file.Name())
				printDir(out, filepath.Join(path, file.Name()), printFiles, prefix+"\t")
			} else {
				fmt.Fprintf(out, "%s├───%s\n", prefix, file.Name())
				printDir(out, filepath.Join(path, file.Name()), printFiles, prefix+"│\t")
			}
		}
	}

	return nil
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
