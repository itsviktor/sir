package ir

import (
	"fmt"

	"github.com/itsviktor/sir/src/internal/utils"
)

type IsPredicate int

const (
	IsNull IsPredicate = iota
	IsTrue
	IsFalse
	IsUnknown
	IsDocument
	IsNormalized
)

func isPredicateToString(predicate IsPredicate) string {
	switch predicate {
	case IsNull:
		return "null"
	case IsTrue:
		return "true"
	case IsFalse:
		return "false"
	case IsUnknown:
		return "unknown"
	case IsDocument:
		return "document"
	case IsNormalized:
		return "normalized"
	}
	return fmt.Sprintf("unknown predicate: %d", predicate)
}

type IsExpr struct {
	Expr      Expr
	Predicate IsPredicate
	Not       bool
	Pos       utils.Position
}

func (IsExpr) expr() {}

func (e *IsExpr) Print(indent int) {
	if e.Not {
		utils.IndentPrintf(indent, "- is not %s expr\n", isPredicateToString(e.Predicate))
	} else {
		utils.IndentPrintf(indent, "- is %s expr\n", isPredicateToString(e.Predicate))
	}

	utils.IndentPrintf(indent+2, " expr:\n")

	e.Expr.Print(indent + 4)
}

func (e *IsExpr) GetPos() utils.Position {
	return e.Pos
}
