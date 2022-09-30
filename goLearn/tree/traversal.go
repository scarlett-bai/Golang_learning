package tree

import "fmt"

func (node *Node) Traverse() {
	node.TraveseFunc(func(n *Node) {
		n.Print()
	})
	fmt.Println()
}

func (node *Node) TraveseFunc(f func(*Node)) {
	if node == nil {
		return
	}

	node.Left.TraveseFunc(f)
	f(node)
	node.Right.TraveseFunc(f)
}

func (node *Node) TraverseWithChannel() chan *Node {
	out := make(chan *Node)
	go func() {
		node.TraveseFunc(func(node *Node) {
			out <- node
		})
		close(out)
	}()
	return out
}
