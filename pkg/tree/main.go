package tree

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Node struct {
	Name     string
	ModTime  time.Time
	IsDir    bool
	Children []*Node
}

// NewNode : Creates a node ready to be inserted into the tree
func NewNode(name string, modTime time.Time, isDir bool) *Node {
	return &Node{
		Name:     name,
		ModTime:  modTime,
		IsDir:    isDir,
		Children: make([]*Node, 0),
	}
}

func Build(folderName string, ignoreDotFiles bool) (*Node, error) {
	stat, err := os.Stat(folderName)
	if err != nil {
		return nil, err
	}
	root := NewNode(folderName, stat.ModTime(), stat.IsDir())
	stack := (Stack)([]*Node{root})
	var currentNode *Node

	for len(stack) != 0 {
		stack, currentNode = stack.pop()

		if currentNode != nil {
			infos, err := ioutil.ReadDir(currentNode.Name)
			if err != nil {
				return nil, err
			}

			for _, val := range infos {
				if ignoreDotFiles && strings.Index(val.Name(), ".") == 0 {
					continue
				}

				childNode := NewNode(filepath.Join(currentNode.Name, val.Name()), val.ModTime(), val.IsDir())
				currentNode.Children = append(currentNode.Children, childNode)
				if val.IsDir() {
					stack = stack.push(childNode)
				}
			}

		}
	}

	return root, nil
}

func Print(root *Node) {
	fmt.Println(root.Name)
	for _, v := range root.Children {
		Print(v)
	}
}
