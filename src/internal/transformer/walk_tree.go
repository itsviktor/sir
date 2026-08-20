package transformer

import "github.com/antlr4-go/antlr/v4"

// WalkTree traverses provided antlr tree using enter and exit functions in this order:
// enter(node) - WalkTree(child) on each child - exit(node)
func WalkTree(node antlr.Tree, enter func(ctx antlr.Tree), exit func(ctx antlr.Tree)) {
	enter(node)

	for i := 0; i < node.GetChildCount(); i++ {
		WalkTree(node.GetChild(i), enter, exit)
	}

	exit(node)
}
