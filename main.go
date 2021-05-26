package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/athiban2001/go-mon/pkg/watch"
	"github.com/fatih/color"
)

func main() {
	foldername := flag.String("f", ".", "Folder to watch changes for")
	ignoreDotFiles := flag.Bool("i", true, "Ignore files and folders that starts with .")
	extensions := flag.String("e", ".go", "Extensions to watch out for comma-separated eg: .go,.html")
	command := flag.String("c", "make run", "Command to re-run the system")
	flag.Parse()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	rootPath, err := filepath.Abs(*foldername)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error : %v\n", err)
	}

	color.Yellow("[go-mon] version 1.0")
	color.Yellow("[go-mon] to restart at any time, press `rs`")
	color.Yellow("[go-mon] watching path: `%s`", rootPath)
	if *ignoreDotFiles {
		color.Yellow("[go-mon] ignoring .* files")
	}
	color.Yellow("[go-mon] watching extensions : `%s`", *extensions)

	watchCtx, watchCancel := context.WithCancel(context.Background())
	events, errors, err := watch.Start(watchCtx, rootPath, *ignoreDotFiles, *extensions)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error : %v\n", err)
	}

	execCtx, execCancel := context.WithCancel(context.Background())
	triggerC := make(chan string)

	go execute(execCtx, rootPath, *command, triggerC)
	triggerC <- "INIT"

	go func() {
		command := ""
		for {
			fmt.Scanln(&command)
			if command == "rs" {
				triggerC <- "RESTART"
			}
		}
	}()

	for {
		select {
		case data := <-events:
			triggerC <- data
		case err := <-errors:
			fmt.Fprintf(os.Stderr, "Error : %v\n", err)
		case <-watchCtx.Done():
			return
		case <-c:
			watchCancel()
			execCancel()
			return
		}
	}
}

func execute(ctx context.Context, folderpath string, command string, triggerC <-chan string) {
	var cmd *exec.Cmd
	commandList := strings.Split(command, " ")
	executable := commandList[0]
	args := commandList[1:]
	errOut := &bytes.Buffer{}

	for val := range triggerC {
		if cmd != nil {
			cmd.Process.Kill()
		}
		errOut.Reset()
		if val == "INIT" {
			color.Green("[go-mon] starting `%s`", command)
		} else {
			color.Green("[go-mon] restarting `%s`", command)
		}
		cmd = exec.CommandContext(ctx, executable, args...)
		cmd.Dir = folderpath
		cmd.Env = os.Environ()
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = errOut
		if err := cmd.Run(); err != nil {
			if strings.Contains(errOut.String(), "no such file or directory") {
				return
			}
			fmt.Fprintf(os.Stderr, "%s", errOut)
		}
	}
}
