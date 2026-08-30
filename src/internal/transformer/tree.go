package transformer

import (
	"fmt"
	"strings"

	"github.com/antlr4-go/antlr/v4"
)

// WalkTree traverses provided antlr tree using enter and exit functions in this order:
// enter(node) - WalkTree(child) on each child - exit(node)
func WalkTree(node antlr.Tree, enter func(ctx antlr.Tree), exit func(ctx antlr.Tree)) {
	enter(node)

	for i := 0; i < node.GetChildCount(); i++ {
		WalkTree(node.GetChild(i), enter, exit)
	}

	exit(node)
}

func PrintTree(node antlr.Tree, limit int) {
	if node == nil {
		return
	}

	printLeaf(node, 0, -1, limit)
}

func printLeaf(node antlr.Tree, offset int, depth int, limit int) {
	if node == nil {
		return
	}

	if depth == limit {
		return
	}

	indent := strings.Repeat(" ", offset)

	fmt.Printf("%senter node: %T\n", indent, node)

	for i := 0; i < node.GetChildCount(); i++ {
		printLeaf(node.GetChild(i), offset+2, depth+1, limit)
	}

	fmt.Printf("%sexit node: %T\n", indent, node)
}

func FindFirst[T antlr.Tree](node antlr.Tree) (T, bool) {
	if target, ok := node.(T); ok {
		return target, true
	}

	for i := 0; i < node.GetChildCount(); i++ {
		target, ok := FindFirst[T](node.GetChild(i))
		if ok {
			return target, true
		}
	}

	var zero T
	return zero, false
}

// FindFirstWide searches the ANTLR tree in breadth-first order and returns the first node
// of type T found at a depth not exceeding maxDeep. The root node has depth 0.
func FindFirstWide[T antlr.Tree](node antlr.Tree, maxDeep int) (T, bool) {
	type queueItem struct {
		node  antlr.Tree
		depth int
	}

	queue := []queueItem{{node: node, depth: 0}}

	for i := 0; i < len(queue); i++ {
		current := queue[i]

		if target, ok := current.node.(T); ok {
			return target, true
		}

		if current.depth >= maxDeep {
			continue
		}

		for j := 0; j < current.node.GetChildCount(); j++ {
			queue = append(queue, queueItem{
				node:  current.node.GetChild(j),
				depth: current.depth + 1,
			})
		}
	}

	var zero T
	return zero, false
}
