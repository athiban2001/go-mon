package watch

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/athiban2001/go-mon/pkg/tree"
)

// Start : Initialize the watching event
func Start(ctx context.Context, foldername string, ignoreDotFiles bool, extensionsString string) (chan string, chan error, error) {
	stat, err := os.Stat(foldername)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil, err
	}
	if err != nil {
		return nil, nil, err
	}
	if !stat.IsDir() {
		return nil, nil, fmt.Errorf("not a directory")
	}

	events := make(chan string)
	errorsC := make(chan error)
	isValid := isValidDecorator(ignoreDotFiles, strings.Split(extensionsString, ","))
	wg := &sync.WaitGroup{}

	root, err := tree.Build(foldername, isValid)
	if err != nil {
		return nil, nil, err
	}

	watchTree(ctx, wg, root, isValid, events, errorsC)

	return events, errorsC, nil
}

func watchTree(ctx context.Context, wg *sync.WaitGroup, root *tree.Node, isValid func(string, bool) bool, events chan<- string, errorsC chan<- error) {
	wg.Add(1)
	defer wg.Done()

	for _, val := range root.Children {
		watchTree(ctx, wg, val, isValid, events, errorsC)
	}

	if root.IsDir {
		go watchDir(ctx, wg, root, isValid, events, errorsC)
		return
	}
	go watchFile(ctx, root, events, errorsC)
}

func watchDir(ctx context.Context, wg *sync.WaitGroup, root *tree.Node, isValid func(string, bool) bool, events chan<- string, errorsC chan<- error) {
	var (
		newlyAddedChildren      []*tree.Node
		stat                    fs.FileInfo
		entries                 []fs.DirEntry
		currentTime, cachedTime time.Time
		err                     error
	)

	wg.Wait()

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			stat, err = os.Stat(root.Name)
			if errors.Is(err, os.ErrNotExist) {
				events <- root.Name + " DELETED"
				return
			}
			if err != nil {
				errorsC <- err
				return
			}
			currentTime = stat.ModTime()
			if cachedTime != currentTime && !cachedTime.IsZero() {
				cachedTime = currentTime
				entries, err = os.ReadDir(root.Name)
				if err != nil {
					errorsC <- err
					return
				}

				for k, val := range entries {
					if !isValid(val.Name(), val.IsDir()) {
						entries = removeEntry(entries, k)
					}
				}

				if len(entries) > len(root.Children) {
					root.Children, newlyAddedChildren = AddChildren(root, entries)
					for _, val := range newlyAddedChildren {
						if val.IsDir {
							go watchDir(ctx, wg, val, isValid, events, errorsC)
						} else {
							go watchFile(ctx, val, events, errorsC)
						}
						events <- val.Name + " INSERTED"
					}
				} else if len(entries) < len(root.Children) {
					root.Children = RemoveChildren(root, entries)
				}
			} else if cachedTime != currentTime {
				cachedTime = currentTime
			}
		case <-ctx.Done():
			return
		}
	}
}

func watchFile(ctx context.Context, root *tree.Node, events chan<- string, errorsC chan<- error) {
	var (
		stat                    fs.FileInfo
		err                     error
		cachedTime, currentTime time.Time
	)

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			stat, err = os.Stat(root.Name)
			if errors.Is(err, os.ErrNotExist) {
				events <- root.Name + " DELETED"
				return
			}
			if err != nil {
				errorsC <- err
				return
			}

			currentTime = stat.ModTime()
			if cachedTime != currentTime && !cachedTime.IsZero() {
				cachedTime = currentTime
				events <- root.Name + " MODIFIED"
			} else if cachedTime != stat.ModTime() {
				cachedTime = currentTime
			}
		case <-ctx.Done():
			return
		}
	}
}
