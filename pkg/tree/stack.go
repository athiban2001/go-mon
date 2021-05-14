package tree

type Stack []*Node

func (s Stack) push(node *Node) Stack {
	return append(s, node)
}

func (s Stack) pop() (Stack, *Node) {
	stackLength := len(s)
	if stackLength == 0 {
		return s, nil
	}

	return s[1:], s[0]
}
