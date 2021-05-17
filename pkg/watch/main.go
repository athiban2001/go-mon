package watch

import (
	"context"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"sync"
	"time"

	"github.com/athiban2001/go-mon/pkg/tree"
)

// Start : Initialize the watching event
func Start(ctx context.Context, foldername string, ignoreDotFiles bool) (chan string, chan error, error) {
	stat, err := os.Stat(foldername)
	if os.IsNotExist(err) {
		return nil, nil, err
	}
	if err != nil {
		return nil, nil, err
	}
	if !stat.IsDir() {
		return nil, nil, fmt.Errorf("not a directory")
	}

	events := make(chan string)
	errors := make(chan error)
	wg := &sync.WaitGroup{}

	root, err := tree.Build(foldername, ignoreDotFiles)
	if err != nil {
		return nil, nil, err
	}

	watchTree(ctx, wg, root, ignoreDotFiles, events, errors)

	return events, errors, nil
}

func watchTree(ctx context.Context, wg *sync.WaitGroup, root *tree.Node, ignoreDotFiles bool, events chan<- string, errors chan<- error) {
	wg.Add(1)
	defer wg.Done()

	for _, val := range root.Children {
		watchTree(ctx, wg, val, ignoreDotFiles, events, errors)
	}

	if root.IsDir {
		go watchDir(ctx, wg, root, ignoreDotFiles, events, errors)
		return
	}
	go watchFile(ctx, root, events, errors)
}

func watchDir(ctx context.Context, wg *sync.WaitGroup, root *tree.Node, ignoreDotFiles bool, events chan<- string, errors chan<- error) {
	wg.Wait()
	fmt.Println("HERE")

	var newlyAddedChildren []*tree.Node
	var stat fs.FileInfo
	var infos []fs.FileInfo
	var err error
	var modTime time.Time

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			stat, err = os.Stat(root.Name)
			if os.IsNotExist(err) {
				events <- root.Name + " DELETED"
				return
			}
			if err != nil {
				errors <- err
				return
			}
			modTime = stat.ModTime()
			if root.ModTime != modTime {
				infos, err = ioutil.ReadDir(root.Name)
				if err != nil {
					errors <- err
					return
				}

				for k, val := range infos {
					if ignoreDotFiles && val.Name()[0] == '.' {
						infos = append(infos[:k], infos[k+1:]...)
					}
				}

				if len(infos) > len(root.Children) {
					root.Children, newlyAddedChildren = AddChildren(root, infos)
					for _, val := range newlyAddedChildren {
						if val.IsDir {
							go watchDir(ctx, wg, val, ignoreDotFiles, events, errors)
						} else {
							go watchFile(ctx, val, events, errors)
						}
						events <- val.Name + " INSERTED"
					}
				} else if len(infos) < len(root.Children) {
					root.Children = RemoveChildren(root, infos)
				}
			}
		case <-ctx.Done():
			return
		}
	}
}

func watchFile(ctx context.Context, root *tree.Node, events chan<- string, errors chan<- error) {
	var stat fs.FileInfo
	var err error
	var modTime time.Time

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			stat, err = os.Stat(root.Name)
			if os.IsNotExist(err) {
				events <- root.Name + " DELETED"
				return
			}
			if err != nil {
				errors <- err
				return
			}
			modTime = stat.ModTime()
			if modTime != root.ModTime {
				root.ModTime = modTime
				events <- root.Name + " MODIFIED"
			}
		case <-ctx.Done():
			return
		}
	}
}
