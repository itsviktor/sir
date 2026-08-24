package ir

type Expr interface {
	expr()
}

type LiteralExpr struct {
	Value string
}

func (LiteralExpr) expr() {}

type SqlExpr struct {
	Sql string
}

func (SqlExpr) expr() {}

type DsqlExpr struct {
	Path []string
}

func (DsqlExpr) expr() {}

type ColumnExpr struct {
	Name     string
	Relation Relation
}

func (ColumnExpr) expr() {}
