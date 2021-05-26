package watch

import (
	"io/fs"
	"path/filepath"
	"sort"

	"github.com/athiban2001/go-mon/pkg/tree"
)

func isValidDecorator(ignoreDotFiles bool, extensions []string) func(string, bool) bool {
	return func(filename string, isDir bool) bool {
		if ignoreDotFiles && filename[0] == '.' {
			return false
		}
		if !isDir {
			extension := filepath.Ext(filename)
			for _, val := range extensions {
				if val == extension {
					return true
				}
			}
			return false
		}
		return true
	}
}

func removeEntry(entries []fs.DirEntry, index int) []fs.DirEntry {
	newEntries := make([]fs.DirEntry, len(entries)-1)
	i, k := 0, 0

	for i = 0; i < len(entries); i++ {
		if i != index {
			newEntries[k] = entries[i]
			k++
		}
	}

	return newEntries
}

func insertChildrenInOrder(oldChildren []*tree.Node, addition []*tree.Node) []*tree.Node {
	for k := range addition {
		i := sort.Search(len(oldChildren), func(i int) bool {
			return oldChildren[i].Name > addition[k].Name
		})
		oldChildren = append(oldChildren, &tree.Node{})
		copy(oldChildren[i+1:], oldChildren[i:])
		oldChildren[i] = addition[k]
	}

	return oldChildren
}

func AddChildren(root *tree.Node, infos []fs.DirEntry) ([]*tree.Node, []*tree.Node) {
	infosLen := len(infos)
	childrenLen := len(root.Children)
	newChildren := make([]*tree.Node, 0)
	i, j := 0, 0

	for i < infosLen && j < childrenLen {
		absFileName := filepath.Join(root.Name, infos[i].Name())
		if absFileName < root.Children[j].Name {
			newChildren = append(newChildren, tree.NewNode(absFileName, infos[i].IsDir()))
			i++
		} else if absFileName > root.Children[j].Name {
			j++
		} else {
			j++
			i++
		}
	}

	k := i
	nodesFromInfo := make([]*tree.Node, 0)
	for k < infosLen {
		nodesFromInfo = append(nodesFromInfo, tree.NewNode(filepath.Join(root.Name, infos[k].Name()), infos[k].IsDir()))
		k++
	}

	if i < infosLen {
		finalLength := infosLen - i + len(newChildren)
		if finalLength > cap(newChildren) {
			newDifference := make([]*tree.Node, finalLength)
			copy(newDifference, newChildren)
			copy(newDifference[len(newChildren):], nodesFromInfo)
			newChildren = newDifference
		} else {
			differenceLen := len(newChildren)
			newChildren = newChildren[:finalLength]
			copy(newChildren[differenceLen:], nodesFromInfo)
		}
	}

	if len(newChildren) != 0 {
		return insertChildrenInOrder(root.Children, newChildren), newChildren
	}

	return root.Children, newChildren
}

func RemoveChildren(root *tree.Node, infos []fs.DirEntry) []*tree.Node {
	children := root.Children
	remainingChildren := make([]*tree.Node, 0)
	i, j := 0, 0
	infosLen, childrenLen := len(infos), len(children)

	for i < infosLen && j < childrenLen {
		absFileName := filepath.Join(root.Name, infos[i].Name())
		if absFileName > children[j].Name {
			j++
		} else if absFileName < children[j].Name {
			i++
		} else {
			remainingChildren = append(remainingChildren, children[j])
			i++
			j++
		}
	}

	return remainingChildren
}
