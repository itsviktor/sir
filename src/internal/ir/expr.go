package ir

import "github.com/itsviktor/sir/src/internal/utils"

type Expr interface {
	expr()
	Print(indent int)
}

type LiteralExpr struct {
	Value string
}

func (LiteralExpr) expr() {}

func (e *LiteralExpr) Print(indent int) {
	utils.IndentPrintf(indent, "- literal: %s\n", e.Value)
}

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

func (e *ColumnExpr) Print(indent int) {
	utils.IndentPrintf(indent, "- columnref:\n")
	utils.IndentPrintf(indent+2, " relation:\n")
	e.Relation.Print(indent + 4)
	utils.IndentPrintf(indent+2, " name: %s\n", e.Name)
}
