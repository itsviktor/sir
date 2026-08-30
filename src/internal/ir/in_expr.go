package ir

import "github.com/itsviktor/sir/src/internal/utils"

type InExpr struct {
	Left  Expr
	Not   bool
	Right InRight
}

func (InExpr) expr() {}

func (e *InExpr) Print(indent int) {
	utils.IndentPrintf(indent, "- in expr:\n")

	switch r := e.Right.(type) {
	case InList:
		utils.IndentPrintf(indent+2, " items:\n")
		for _, item := range r.Items {
			item.Print(indent + 4)
		}
	case InSubquery:
		utils.IndentPrintf(indent+2, " subquery:\n")
		r.Query.Print(indent + 4)
	}
}

type InRight interface {
	inRight()
}

type InList struct{ Items []Expr }

func (InList) inRight() {}

type InSubquery struct{ Query *SelectQuery }

func (InSubquery) inRight() {}
