package ir

import "github.com/itsviktor/sir/src/internal/utils"

type BetweenExpr struct {
	Left Expr
	From Expr
	To   Expr
	Not  bool
	Pos  utils.Position
}

func (BetweenExpr) expr() {}

func (e *BetweenExpr) Print(indent int) {
	if e.Not {
		utils.IndentPrintf(indent, "- not between expr:\n")
	} else {
		utils.IndentPrintf(indent, "- between expr:\n")
	}

	utils.IndentPrintf(indent+2, " left: \n")
	e.Left.Print(indent + 4)

	utils.IndentPrintf(indent+2, " from: \n")
	e.From.Print(indent + 4)

	utils.IndentPrintf(indent+2, " to: \n")
	e.To.Print(indent + 4)
}

func (e *BetweenExpr) GetPos() utils.Position {
	return e.Pos
}
