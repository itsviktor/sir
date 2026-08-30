package ir

import "github.com/itsviktor/sir/src/internal/utils"

type ReturnExpr interface {
	Expr
	returnExpr()
}

type AllReturnExpr struct {
}

func (AllReturnExpr) expr() {}

func (AllReturnExpr) returnExpr() {}

func (AllReturnExpr) Print(indent int) {
	utils.IndentPrintf(indent, "- all fields")
}
