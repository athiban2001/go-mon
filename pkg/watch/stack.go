package watch

type Stack []*TreeNode

func (s Stack) push(node *TreeNode) Stack {
	return append(s, node)
}

func (s Stack) pop() (Stack, *TreeNode) {
	stackLength := len(s)
	if stackLength == 0 {
		return s, nil
	}

	return s[1:], s[0]
}
