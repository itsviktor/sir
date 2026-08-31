package ir

import "github.com/itsviktor/sir/src/internal/utils"

type LogicalExpr struct {
	Left  Expr
	Op    Op
	Right Expr
	Pos   utils.Position
}

func (LogicalExpr) expr() {}

func (e *LogicalExpr) Print(indent int) {
	if e.Op.Type == And {
		utils.IndentPrintf(indent, "- AND expr:\n")
	} else {
		utils.IndentPrintf(indent, "- OR op:\n")
	}

	utils.IndentPrintf(indent+2, " left:\n")
	e.Left.Print(indent + 4)

	utils.IndentPrintf(indent+2, " right:\n")
	e.Right.Print(indent + 4)
}

func (e *LogicalExpr) GetPos() utils.Position {
	return e.Pos
}
