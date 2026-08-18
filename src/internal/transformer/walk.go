package transformer

import "github.com/antlr4-go/antlr/v4"

func WalkAntlrTree(node antlr.Tree, visit func(antlr.ParserRuleContext)) {
	if ctx, ok := node.(antlr.ParserRuleContext); ok {
		visit(ctx)
	}
	for i := 0; i < node.GetChildCount(); i++ {
		WalkAntlrTree(node.GetChild(i), visit)
	}
}
