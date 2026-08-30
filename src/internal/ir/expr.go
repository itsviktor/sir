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

type DsqlExpr struct {
	Name  string
	Field *string
}

func (DsqlExpr) expr() {}

func (e *DsqlExpr) Print(indent int) {
	utils.IndentPrintf(indent, "- DSQL param:\n")
	utils.IndentPrintf(indent+2, " name: %s\n", e.Name)
	if e.Field != nil {
		utils.IndentPrintf(indent+2, " field: %s\n", *e.Field)
	}
}

type ColumnExpr struct {
	Name     string
	Relation Relation
}

func (ColumnExpr) expr() {}

func (ColumnExpr) returnExpr() {}

func (e *ColumnExpr) Print(indent int) {
	utils.IndentPrintf(indent, "- columnref:\n")
	utils.IndentPrintf(indent+2, " relation:\n")
	e.Relation.Print(indent + 4)
	utils.IndentPrintf(indent+2, " name: %s\n", e.Name)
}

type WildcardColumnExpr struct {
	Relation Relation
}

func (WildcardColumnExpr) expr() {}

func (WildcardColumnExpr) returnExpr() {}

func (e *WildcardColumnExpr) Print(indent int) {
	utils.IndentPrintf(indent, "- wildcard column ref:\n")
	utils.IndentPrintf(indent+2, " relation:\n")
	e.Relation.Print(indent + 4)
}
