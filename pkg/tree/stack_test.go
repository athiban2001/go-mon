package tree

import "testing"

func TestPush(t *testing.T) {
	stack := make(Stack, 0)
	tests := []struct {
		input    *Node
		expected int
	}{
		{input: &Node{}, expected: 1},
		{input: &Node{}, expected: 2},
		{input: &Node{}, expected: 3},
	}

	for _, test := range tests {
		stack = stack.push(test.input)
		if len(stack) != test.expected {
			t.Errorf("Expected Length : %v, Received Length : %v", test.expected, len(stack))
		}
	}
}

func TestPop(t *testing.T) {
	var node *Node
	stack := make(Stack, 0)
	stack = stack.push(&Node{})
	stack = stack.push(&Node{})
	stack = stack.push(&Node{})

	tests := []struct {
		expected *Node
	}{
		{expected: stack[0]},
		{expected: stack[1]},
		{expected: stack[2]},
		{expected: nil},
	}

	for _, test := range tests {
		stack, node = stack.pop()
		if node != test.expected {
			t.Errorf("Expected : %v, Received : %v", test.expected, node)
		}
	}
}
