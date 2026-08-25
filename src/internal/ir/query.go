package ir

import "fmt"

type Query interface {
	query()
	Print()
	AddTarget(target Relation)
}

type SelectQuery struct {
	Targets []Relation

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

func (s *SelectQuery) Print() {
	fmt.Printf("SELECT QUERY\n")
	if len(s.Targets) > 0 {
		fmt.Printf("FROM:\n")
		for _, target := range s.Targets {
			target.Print(4)
		}
	}

}

func (s *SelectQuery) AddTarget(target Relation) {
	s.Targets = append(s.Targets, target)
}
