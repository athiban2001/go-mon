package tree

import (
	"testing"
	"time"
)

func TestNewNode(t *testing.T) {
	name := "test"
	modTime := time.Now()
	isDir := true

	node1 := NewNode(name, modTime, isDir)
	node2 := NewNode(name, modTime, isDir)
	if node1 == nil || node2 == nil {
		t.Errorf("Expected not nil value")
	}
	if node1 == node2 {
		t.Errorf("Expected unequal values between subsequent calls")
	}
}
