package ir

import (
	"github.com/itsviktor/sir/src/internal/utils"
)

type Query interface {
	query()
	Print(indent int)
	AddTarget(target Relation)
}

// Internal Representation for the SELECT query.
// Contains every part of the query, such as target, returns, where clause, offset clause, limit clause.
type SelectQuery struct {
	Targets []Relation
	Returns []ReturnExpr

	Where  Expr
	Offset Expr
	Limit  Expr
}

// Creates new Internal Representation for the SELECT query.
func NewSelectQuery() *SelectQuery {
	return &SelectQuery{
		Targets: make([]Relation, 0),
		Returns: make([]ReturnExpr, 0),
	}
}

func (SelectQuery) query() {}

func (s *SelectQuery) Print(indent int) {
	utils.IndentPrintf(indent, "SELECT QUERY\n")

	if len(s.Returns) > 0 {
		utils.IndentPrintf(indent, "RETURNS:\n")
		for _, ret := range s.Returns {
			ret.Print(indent + 2)
		}
	}

	if len(s.Targets) > 0 {
		utils.IndentPrintf(indent, "FROM:\n")
		for _, target := range s.Targets {
			target.Print(indent + 2)
		}
	}

	if s.Where != nil {
		utils.IndentPrintf(indent, "WHERE:\n")
		s.Where.Print(indent + 2)
	}
}
