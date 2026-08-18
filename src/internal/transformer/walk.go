package transformer

import "github.com/antlr4-go/antlr/v4"

func WalkAntlrTree(node antlr.Tree, enter func(ctx antlr.Tree), exit func(ctx antlr.Tree)) {
	enter(node)

	for i := 0; i < node.GetChildCount(); i++ {
		WalkAntlrTree(node.GetChild(i), enter, exit)
	}

	exit(node)
}
