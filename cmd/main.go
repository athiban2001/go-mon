package main

import (
	"flag"
	"fmt"
	"path/filepath"

	"github.com/athiban2001/go-mon/pkg/watch"
)

func main() {
	foldername := flag.String("foldername", ".", "Folder to watch changes for")
	ignoreDotFiles := flag.Bool("ignoredot", true, "Ignore files and folders that starts with .")
	flag.Parse()

	rootPath, err := filepath.Abs(*foldername)
	if err != nil {
		fmt.Printf("Error : %v\n", err)
	}

	fileChanges, err := watch.StartWatch(rootPath, *ignoreDotFiles)
	if err != nil {
		fmt.Printf("Error : %v\n", err)
	}

	for {
		data := <-fileChanges
		fmt.Println(data)
	}
}
