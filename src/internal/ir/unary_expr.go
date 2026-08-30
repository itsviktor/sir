package ir

import (
	"github.com/itsviktor/sir/src/internal/utils"
)

type UnaryExpr struct {
	Expr Expr
	Op   Op
}

func (UnaryExpr) expr() {}

func (e *UnaryExpr) Print(indent int) {
	utils.IndentPrintf(indent, "- unary expr:\n")

	utils.IndentPrintf(indent+2, " expr:\n")
	e.Expr.Print(indent + 4)

	utils.IndentPrintf(indent+2, " op: %s\n", e.Op.String())
}
