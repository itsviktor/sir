package ir

import (
	"github.com/itsviktor/sir/src/internal/utils"
)

type Expr interface {
	expr()
	Print(indent int)
	GetPos() utils.Position
}

type LiteralExpr struct {
	Value string
	Pos   utils.Position
}

func (LiteralExpr) expr() {}

func (e *LiteralExpr) Print(indent int) {
	utils.IndentPrintf(indent, "- literal: %s\n", e.Value)
}

func (e *LiteralExpr) GetPos() utils.Position {
	return e.Pos
}

type DsqlExpr struct {
	Name  string
	Field *string
	Pos   utils.Position
}

func (DsqlExpr) expr() {}

func (e *DsqlExpr) Print(indent int) {
	utils.IndentPrintf(indent, "- DSQL param:\n")
	utils.IndentPrintf(indent+2, " name: %s\n", e.Name)
	if e.Field != nil {
		utils.IndentPrintf(indent+2, " field: %s\n", *e.Field)
	}
}

func (e *DsqlExpr) GetPos() utils.Position {
	return e.Pos
}

type ColumnExpr struct {
	Name     string
	Relation *TableRelation
	Pos      utils.Position
}

func (ColumnExpr) expr() {}

func (ColumnExpr) returnExpr() {}

func (e *ColumnExpr) Print(indent int) {
	utils.IndentPrintf(indent, "- columnref:\n")
	utils.IndentPrintf(indent+2, " relation:\n")
	e.Relation.Print(indent + 4)
	utils.IndentPrintf(indent+2, " name: %s\n", e.Name)
}

func (e *ColumnExpr) GetPos() utils.Position {
	return e.Pos
}

type WildcardColumnExpr struct {
	Relation Relation
	Pos      utils.Position
}

func (WildcardColumnExpr) expr() {}

func (WildcardColumnExpr) returnExpr() {}

func (e *WildcardColumnExpr) Print(indent int) {
	utils.IndentPrintf(indent, "- wildcard column ref:\n")
	utils.IndentPrintf(indent+2, " relation:\n")
	e.Relation.Print(indent + 4)
}

func (e *WildcardColumnExpr) GetPos() utils.Position {
	return e.Pos
}
