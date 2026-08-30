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

	utils.IndentPrintf(indent+2, " left:\n")
	e.Left.Print(indent + 4)

	switch r := e.Right.(type) {
	case InList:
		utils.IndentPrintf(indent+2, " in items:\n")
		for _, item := range r.Items {
			item.Print(indent + 4)
		}
	case InSubquery:
		utils.IndentPrintf(indent+2, " in subquery:\n")
		r.Query.Print(indent + 4)
	case InDsql:
		utils.IndentPrintf(indent+2, " in dsql:\n")
		r.Expr.Print(indent + 4)
	}
}

type InRight interface {
	inRight()
}

type InList struct{ Items []Expr }

func (InList) inRight() {}

type InSubquery struct{ Query *SelectQuery }

func (InSubquery) inRight() {}

type InDsql struct{ Expr Expr }

func (InDsql) inRight() {}
