package watch

import (
	"errors"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/athiban2001/go-mon/pkg/tree"
)

// Start : Initialize the watching event
func Start(foldername string, ignoreDotFiles bool) (chan string, chan error, error) {
	stat, err := os.Stat(foldername)
	if os.IsNotExist(err) {
		return nil, nil, err
	}
	if !stat.IsDir() {
		return nil, nil, errors.New("not a directory")
	}

	events := make(chan string)
	errors := make(chan error)

	root, err := tree.Build(foldername, ignoreDotFiles)
	if err != nil {
		return nil, nil, err
	}

	wg := &sync.WaitGroup{}
	watchTree(wg, root, ignoreDotFiles, events, errors)

	return events, errors, nil
}

func watchTree(wg *sync.WaitGroup, root *tree.Node, ignoreDotFiles bool, events chan<- string, errors chan<- error) {
	wg.Add(1)

	if !root.IsDir {
		go watchFile(root, events, errors)
		return
	}
	go watchDir(wg, root, ignoreDotFiles, events, errors)

	for _, val := range root.Children {
		watchTree(wg, val, ignoreDotFiles, events, errors)
	}

	wg.Done()
}

func watchDir(wg *sync.WaitGroup, root *tree.Node, ignoreDotFiles bool, events chan<- string, errors chan<- error) {
	wg.Wait()
	for {
		stat, err := os.Stat(root.Name)
		if os.IsNotExist(err) {
			events <- root.Name + " DELETED"
			return
		}
		if err != nil {
			errors <- err
			return
		}
		if root.ModTime != stat.ModTime() {
			infos, err := ioutil.ReadDir(root.Name)
			if err != nil {
				errors <- err
				return
			}

			for k, val := range infos {
				if strings.Index(val.Name(), ".") == 0 {
					infos = append(infos[:k], infos[k+1:]...)
				}
			}

			if len(infos) > len(root.Children) {
				newlyAddedFiles := ArrayDifference(root, infos)
				if len(newlyAddedFiles) != 0 {
					root.Children = InsertChildren(root.Children, newlyAddedFiles)

					for _, val := range newlyAddedFiles {
						if val.IsDir {
							go watchDir(wg, val, ignoreDotFiles, events, errors)
						} else {
							go watchFile(val, events, errors)
						}
						events <- val.Name + " INSERTED"
					}
				}
			} else if len(infos) < len(root.Children) {
				root.Children = GetRemainingChildren(root, infos)
			}
		}

		time.Sleep(500 * time.Millisecond)
	}
}

// watchFile : Watches for the file changes
func watchFile(root *tree.Node, events chan<- string, errors chan<- error) {
	for {
		stat, err := os.Stat(root.Name)
		if os.IsNotExist(err) {
			events <- root.Name + " DELETED"
			return
		}
		if err != nil {
			errors <- err
			return
		}
		if stat.ModTime() != root.ModTime {
			root.ModTime = stat.ModTime()
			events <- root.Name + " MODIFIED"
		}

		time.Sleep(500 * time.Millisecond)
	}
}
