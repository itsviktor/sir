package ir

type Query interface {
	query()
}

type SelectQuery struct {
	From Relation

	Joins []Join

	Where  Expr
	Offset Expr
	Limit  Expr
}

func (SelectQuery) query() {}
