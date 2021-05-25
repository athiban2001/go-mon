package tree

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestNewNode(t *testing.T) {
	name := "test.go"
	isDir := true

	node1 := NewNode(name, isDir)
	node2 := NewNode(name, isDir)
	if node1 == nil || node2 == nil {
		t.Errorf("Expected not nil value")
	}
	if node1 == node2 {
		t.Errorf("Expected unequal values between subsequent calls")
	}
}

func TestBuild(t *testing.T) {
	testPaths := `/home/athiban/Programming/go-mon
	:/home/athiban/Programming/go-mon/Makefile
	:/home/athiban/Programming/go-mon/go.mod
	:/home/athiban/Programming/go-mon/go.sum
	:/home/athiban/Programming/go-mon/gomon
	:/home/athiban/Programming/go-mon/main.go
	:/home/athiban/Programming/go-mon/pkg
	:/home/athiban/Programming/go-mon/pkg/tree
	:/home/athiban/Programming/go-mon/pkg/tree/main.go
	:/home/athiban/Programming/go-mon/pkg/tree/main_test.go
	:/home/athiban/Programming/go-mon/pkg/tree/stack.go
	:/home/athiban/Programming/go-mon/pkg/tree/stack_test.go
	:/home/athiban/Programming/go-mon/pkg/watch
	:/home/athiban/Programming/go-mon/pkg/watch/main.go
	:/home/athiban/Programming/go-mon/pkg/watch/utils.go
	:/home/athiban/Programming/go-mon/pkg/watch/utils_test.go
	:`

	foldername, err := filepath.Abs("../../")
	if err != nil {
		t.Fatalf(err.Error())
	}

	root, err := Build(foldername, func(s string, b bool) bool {
		return s[0] != '.'
	})
	if err != nil {
		t.Fatalf(err.Error())
	}

	paths := new(string)
	allPaths(root, paths)

	if *paths != testPaths {
		t.Errorf("Expected \n%s, Received \n%s", testPaths, *paths)
	}
}

// 39818590463 ns/op
// 39.829 s
func BenchmarkBuild(b *testing.B) {
	foldername, err := filepath.Abs("../../../")
	if err != nil {
		b.Fatalf(err.Error())
	}

	for i := 0; i < b.N; i++ {
		_, err = Build(foldername, func(s string, b bool) bool {
			if strings.Contains(s, "postgres-data") {
				return false
			}
			return s[0] != '.'
		})

		if err != nil {
			b.Fatalf(err.Error())
		}
	}

}
