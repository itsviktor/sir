package postgres

import (
	"fmt"

	"github.com/itsviktor/sir/src/internal/ir"
	"github.com/itsviktor/sir/src/internal/parser"
	"github.com/itsviktor/sir/src/internal/transformer"
)

func parseExpr(aExpr *parser.A_exprContext, tctx *transformer.Context) ir.Expr {
	fmt.Printf("PARSE EXPRESSION: %s\n", aExpr.GetText())

	transformer.PrintTree(aExpr, 100)

	fmt.Printf("\n\n")

	return &ir.LiteralExpr{
		Value: "a",
	}
}
