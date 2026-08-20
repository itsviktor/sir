package ir

import "github.com/itsviktor/sir/src/internal/utils"

type SQLPart interface {
	isSqlPart()
}

type SQLStaticPart struct {
	Start int
	End   int
}

func (SQLStaticPart) isSqlPart() {}

type Expr interface {
	isExpr()
}

type BinaryExpr struct {
	Left  Expr
	Op    string
	Right Expr
}

func (BinaryExpr) isExpr() {}

type VariableExpr struct {
	Name  string
	Field *string
	Pos   utils.FilePosition
}

func (VariableExpr) isExpr() {}

type ColumnExpr struct {
	Name     string
	Relation *string
}

func (ColumnExpr) isExpr() {}

type SQLWherePart struct {
	Expression Expr
}

func (SQLWherePart) isSqlPart() {}

type SQLLimitPart struct {
	Variable VariableExpr
}

func (SQLLimitPart) isSqlPart() {}

type SQLOffsetPart struct {
	Variable VariableExpr
}

func (SQLOffsetPart) isSqlPart() {}

type QueryIR struct {
	Name      string
	SQL       string
	Parts     []SQLPart
	Relations map[string]any
	Aliases   map[string]string
}
