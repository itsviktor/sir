package ir

import "github.com/itsviktor/sir/src/internal/utils"

type LimitAllExpr struct {
}

func (LimitAllExpr) expr() {}

func (LimitAllExpr) Print(indent int) {
	utils.IndentPrintf(indent, "- ALL\n")
}
