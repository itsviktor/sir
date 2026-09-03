package ir

import "github.com/itsviktor/sir/src/internal/utils"

type ReturnExpr interface {
	Expr
	returnExpr()
	String() string
}

type AllReturnExpr struct {
	Pos utils.Position
}

func (AllReturnExpr) expr() {}

func (AllReturnExpr) returnExpr() {}

func (AllReturnExpr) Print(indent int) {
	utils.IndentPrintf(indent, "- all fields\n")
}

func (AllReturnExpr) String() string {
	return "*"
}

func (e *AllReturnExpr) GetPos() utils.Position {
	return e.Pos
}
