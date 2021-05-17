package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/athiban2001/go-mon/pkg/watch"
	"github.com/fatih/color"
)

func main() {
	foldername := flag.String("foldername", ".", "Folder to watch changes for")
	ignoreDotFiles := flag.Bool("ignoredot", true, "Ignore files and folders that starts with .")
	flag.Parse()

	color.Yellow("[go-mon] version 1.0")
	color.Yellow("[go-mon] to restart at any time, press `rs`")
	color.Yellow("[go-mon] watching path: `%s`", *foldername)
	if *ignoreDotFiles {
		color.Yellow("[go-mon] ignoring .* files")
	}

	rootPath, err := filepath.Abs(*foldername)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error : %v\n", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	events, errors, err := watch.Start(ctx, rootPath, *ignoreDotFiles)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error : %v\n", err)
	}

	go func() {
		command := ""
		fmt.Scanln(&command)
		if command == "rs" {
			cancel()
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case data := <-events:
			fmt.Println(data)
		case err := <-errors:
			fmt.Fprintf(os.Stderr, "Error : %v\n", err)
		case <-ctx.Done():
			return
		case <-c:
			cancel()
			return
		}
	}
}
