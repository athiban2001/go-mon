package watch

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type TreeNode struct {
	Name     string
	ModTime  time.Time
	IsDir    bool
	Children []*TreeNode
	*sync.Mutex
}

func NewTreeNode(name string, modTime time.Time, isDir bool) *TreeNode {
	return &TreeNode{
		Name:     name,
		ModTime:  modTime,
		IsDir:    isDir,
		Children: []*TreeNode{},
		Mutex:    &sync.Mutex{},
	}
}

// StartWatch : Initialize the watching event
func StartWatch(foldername string, ignoreDotFiles bool) (chan string, error) {
	stat, err := os.Stat(foldername)
	if os.IsNotExist(err) {
		return nil, err
	}
	if !stat.IsDir() {
		return nil, errors.New("not a directory")
	}

	fileChanges := make(chan string)
	eventChanges := make(chan string)

	root, err := BuildTree(foldername, ignoreDotFiles)
	if err != nil {
		return nil, err
	}

	wg := &sync.WaitGroup{}
	WatchTreeFiles(wg, root, ignoreDotFiles, fileChanges)

	go func() {
		for {
			data := <-fileChanges
			eventChanges <- data
		}
	}()

	return eventChanges, nil
}

func BuildTree(folderName string, ignoreDotFiles bool) (*TreeNode, error) {
	var currentNode *TreeNode

	stat, err := os.Stat(folderName)
	if err != nil {
		return nil, err
	}
	root := NewTreeNode(folderName, stat.ModTime(), stat.IsDir())
	stack := (Stack)([]*TreeNode{root})

	for len(stack) != 0 {
		stack, currentNode = stack.pop()

		if currentNode != nil {
			infos, err := ioutil.ReadDir(currentNode.Name)
			if err != nil {
				return nil, err
			}

			for _, val := range infos {
				if strings.Index(val.Name(), ".") == 0 {
					continue
				}

				childNode := NewTreeNode(filepath.Join(currentNode.Name, val.Name()), val.ModTime(), val.IsDir())
				currentNode.Children = append(currentNode.Children, childNode)
				if val.IsDir() {
					stack = stack.push(childNode)
				}
			}

		}
	}

	return root, nil
}

func PrintTree(root *TreeNode) {
	fmt.Println(root.Name)
	for _, v := range root.Children {
		PrintTree(v)
	}
}

func WatchTreeFiles(wg *sync.WaitGroup, root *TreeNode, ignoreDotFiles bool, fileChanges chan<- string) {
	wg.Add(1)
	if !root.IsDir {
		go WatchFile(root.Name, root.ModTime, fileChanges)
	} else {
		go WatchDir(wg, root, ignoreDotFiles, fileChanges)
	}
	for _, val := range root.Children {
		WatchTreeFiles(wg, val, ignoreDotFiles, fileChanges)
	}
	wg.Done()
}

func WatchDir(wg *sync.WaitGroup, root *TreeNode, ignoreDotFiles bool, fileChanges chan<- string) {
	wg.Wait()
	for {
		_, err := os.Stat(root.Name)
		if os.IsNotExist(err) {
			fileChanges <- root.Name + " DELETED"
			return
		}
		if err != nil {
			return
		}

		infos, err := ioutil.ReadDir(root.Name)
		if err != nil {
			return
		}

		for k, val := range infos {
			if strings.Index(val.Name(), ".") == 0 {
				infos = append(infos[:k], infos[k+1:]...)
			}
		}

		if len(infos) > len(root.Children) {
			newlyAddedFiles := ArrayDifference(root, infos, root.Children, true)
			fmt.Println(root.Name, " ", len(newlyAddedFiles), len(infos), " ", len(root.Children))
			if len(newlyAddedFiles) != 0 {
				root.Children = InsertChildren(root.Children, newlyAddedFiles)

				for _, val := range newlyAddedFiles {
					fmt.Println(val)
					if val.IsDir {
						go WatchDir(wg, val, ignoreDotFiles, fileChanges)
					} else {
						go WatchFile(val.Name, val.ModTime, fileChanges)
					}
					fileChanges <- val.Name + " INSERTED"
				}
			}

		}
		time.Sleep(500 * time.Millisecond)
	}
}

// WatchFile : Watches for the file changes
func WatchFile(filename string, modTime time.Time, fileChanges chan<- string) {
	for {
		stat, err := os.Stat(filename)
		if os.IsNotExist(err) {
			fileChanges <- filename + " DELETED"
			return
		}
		if err != nil {
			return
		}

		if stat.ModTime() != modTime {
			modTime = stat.ModTime()
			fileChanges <- filename + " MODIFIED"
		}
		time.Sleep(500 * time.Millisecond)
	}
}
