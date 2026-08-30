package ir

import "github.com/itsviktor/sir/src/internal/utils"

type IsDistinctExpr struct {
	Left  Expr
	Right Expr
	Not   bool
}

func (IsDistinctExpr) expr() {}

func (e *IsDistinctExpr) Print(indent int) {
	if e.Not {
		utils.IndentPrintf(indent, "- is not distinct expr:\n")
	} else {
		utils.IndentPrintf(indent, "- is distinct expr:\n")
	}

	utils.IndentPrintf(indent+2, " left:\n")
	e.Left.Print(indent + 4)

	utils.IndentPrintf(indent+2, " right:\n")
	e.Right.Print(indent + 4)
}
