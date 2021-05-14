package watch

import (
	"io/fs"
	"path/filepath"
	"sort"

	"github.com/athiban2001/go-mon/pkg/tree"
)

func ArrayDifference(root *tree.Node, infos []fs.FileInfo) []*tree.Node {
	infosLen := len(infos)
	childrenLen := len(root.Children)
	difference := make([]*tree.Node, 0)

	i, j := 0, 0
	for i < infosLen && j < childrenLen {
		absFileName := filepath.Join(root.Name, infos[i].Name())
		if absFileName < root.Children[j].Name {
			difference = append(difference, tree.NewNode(absFileName, infos[i].ModTime(), infos[i].IsDir()))
			i++
		} else if absFileName > root.Children[j].Name {
			j++
		} else {
			j++
			i++
		}
	}

	k := i
	appendArray := make([]*tree.Node, 0)
	for k < infosLen {
		appendArray = append(appendArray, tree.NewNode(filepath.Join(root.Name, infos[k].Name()), infos[k].ModTime(), infos[k].IsDir()))
		k++
	}

	if i < infosLen {
		finalLength := infosLen - i + len(difference)
		if finalLength > cap(difference) {
			newDifference := make([]*tree.Node, finalLength)
			copy(newDifference, difference)
			copy(newDifference[len(difference):], appendArray)
			difference = newDifference
		} else {
			differenceLen := len(difference)
			difference = difference[:finalLength]
			copy(difference[differenceLen:], appendArray)
		}
	}

	return difference
}

func GetRemainingChildren(root *tree.Node, infos []fs.FileInfo) []*tree.Node {
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

func InsertChildren(oldChildren []*tree.Node, addition []*tree.Node) []*tree.Node {
	for k := range addition {
		i := sort.Search(len(oldChildren), func(i int) bool {
			return oldChildren[i].Name > addition[k].Name
		})
		oldChildren = append(oldChildren, &tree.Node{})
		copy(oldChildren[i+1:], oldChildren[i:])
		oldChildren[i] = addition[0]
	}

	return oldChildren
}
