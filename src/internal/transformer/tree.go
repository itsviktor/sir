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

func FindFirstWide[T antlr.Tree](node antlr.Tree) (T, bool) {
	queue := []antlr.Tree{node}

	for i := 0; i < len(queue); i++ {
		current := queue[i]

		if target, ok := current.(T); ok {
			return target, true
		}

		for j := 0; j < current.GetChildCount(); j++ {
			queue = append(queue, current.GetChild(j))
		}
	}

	var zero T
	return zero, false
}
