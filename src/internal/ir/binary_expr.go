package ir

import (
	"github.com/itsviktor/sir/src/internal/utils"
)

type BinaryExpr struct {
	Right Expr
	Op    Op
	Left  Expr
	Pos   utils.Position
}

func (BinaryExpr) expr() {}

func (e *BinaryExpr) Print(indent int) {
	utils.IndentPrintf(indent, "- binary expr:\n")

	utils.IndentPrintf(indent+2, " left: \n")
	e.Left.Print(indent + 4)

	utils.IndentPrintf(indent+2, " right: \n")
	e.Right.Print(indent + 4)

	utils.IndentPrintf(indent+2, " op: %s\n", e.Op.String())
}

func (e *BinaryExpr) GetPos() utils.Position {
	return e.Pos
}
