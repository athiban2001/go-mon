package tree

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/fatih/color"
)

type Node struct {
	Name     string
	IsDir    bool
	Children []*Node
}

// NewNode : Creates a node ready to be inserted into the tree
func NewNode(name string, isDir bool) *Node {
	return &Node{
		Name:     name,
		IsDir:    isDir,
		Children: make([]*Node, 0),
	}
}

func Build(foldername string, isValid func(string, bool) bool) (*Node, error) {
	root := NewNode(foldername, true)
	wg := &sync.WaitGroup{}

	buildRecurse(wg, root, isValid)
	wg.Wait()

	return root, nil
}

func buildRecurse(wg *sync.WaitGroup, root *Node, isValid func(string, bool) bool) {
	var (
		filename string
		isDir    bool
	)
	entries, err := os.ReadDir(root.Name)
	if err != nil {
		color.Red("[go-mon] %v", err)
		color.Red("[go-mon] Ignoring Errored Files")
	}

	for _, entry := range entries {
		filename = entry.Name()
		isDir = entry.IsDir()
		if !isValid(filename, isDir) {
			continue
		}

		childNode := NewNode(filepath.Join(root.Name, filename), isDir)
		root.Children = append(root.Children, childNode)
		if childNode.IsDir {
			wg.Add(1)
			go func() {
				buildRecurse(wg, childNode, isValid)
				wg.Done()
			}()
		}
	}
}

func allPaths(root *Node, paths *string) string {
	*paths += root.Name + "\n\t:"
	for _, v := range root.Children {
		allPaths(v, paths)
	}

	return *paths
}
