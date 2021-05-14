package tree

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

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
