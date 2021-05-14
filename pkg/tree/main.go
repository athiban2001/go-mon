package tree

import (
	"fmt"
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
		Children: []*Node{},
	}
}

func Print(root *Node) {
	fmt.Println(root.Name)
	for _, v := range root.Children {
		Print(v)
	}
}
