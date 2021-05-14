package watch

import (
	"io/fs"
	"path/filepath"
	"sort"
	"strings"
)

func ArrayDifference(root *TreeNode, X []fs.FileInfo, Y []*TreeNode, ignoreDotFiles bool) []*TreeNode {
	Xlen := len(X)
	Ylen := len(Y)
	difference := make([]*TreeNode, 0)

	i, j := 0, 0
	for i < Xlen && j < Ylen {
		absFileName := filepath.Join(root.Name, X[i].Name())
		if ignoreDotFiles && strings.Index(X[i].Name(), ".") == 0 {
			i++
		} else if absFileName < Y[j].Name {
			difference = append(difference, NewTreeNode(absFileName, X[i].ModTime(), X[i].IsDir()))
			i++
		} else if absFileName > Y[j].Name {
			j++
		} else {
			j++
			i++
		}
	}

	k := i
	appendArray := []*TreeNode{}
	for k < Xlen {
		appendArray = append(appendArray, NewTreeNode(filepath.Join(root.Name, X[k].Name()), X[k].ModTime(), X[k].IsDir()))
		k++
	}

	if i < Xlen {
		finalLength := Xlen - i + len(difference)
		if finalLength > cap(difference) {
			newDifference := make([]*TreeNode, finalLength)
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

func InsertChildren(oldChildren []*TreeNode, addition []*TreeNode) []*TreeNode {
	for k := range addition {
		i := sort.Search(len(oldChildren), func(i int) bool {
			return oldChildren[i].Name > addition[k].Name
		})
		oldChildren = append(oldChildren, &TreeNode{})
		copy(oldChildren[i+1:], oldChildren[i:])
		oldChildren[i] = addition[0]
	}

	return oldChildren
}
