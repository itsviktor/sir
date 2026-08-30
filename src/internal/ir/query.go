package ir

import (
	"github.com/itsviktor/sir/src/internal/utils"
)

type Query interface {
	query()
	Print(indent int)
	AddTarget(target Relation)
}

type SelectQuery struct {
	Targets []Relation
	Returns []ReturnExpr

	Where  Expr
	Offset Expr
	Limit  Expr
}

func NewSelectQuery() *SelectQuery {
	return &SelectQuery{
		Targets: make([]Relation, 0),
	}
}

func (SelectQuery) query() {}

func (s *SelectQuery) Print(indent int) {
	utils.IndentPrintf(indent, "SELECT QUERY\n")
	if len(s.Targets) > 0 {
		utils.IndentPrintf(indent, "FROM:\n")
		for _, target := range s.Targets {
			target.Print(indent + 2)
		}
	}

}
