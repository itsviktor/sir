package ir

import "github.com/itsviktor/sir/src/internal/utils"

type OrderByNulls int

const (
	NullsFirst OrderByNulls = iota
	NullsLast
)

func orderByNullsToString(nulls OrderByNulls) string {
	if nulls == NullsFirst {
		return "NULLS FIRST"
	} else {
		return "NULLS LAST"
	}
}

type OrderByOrder int

const (
	Asc OrderByOrder = iota
	Desc
)

func orderByOrderToString(order OrderByOrder) string {
	if order == Asc {
		return "ASC"
	} else {
		return "DESC"
	}
}

type OrderByExpr struct {
	Expr    Expr
	Order   OrderByOrder
	UsingOp *Op
	Nulls   *OrderByNulls

	Pos utils.Position
}

func (e OrderByExpr) Print(indent int) {
	utils.IndentPrintf(indent, "- ORDER BY:\n")

	utils.IndentPrintf(indent+2, " expr:\n")
	e.Expr.Print(indent + 4)

	if e.UsingOp == nil {
		utils.IndentPrintf(indent+2, " order: %s\n", orderByOrderToString(e.Order))
	} else {
		utils.IndentPrintf(indent+2, " using op: %s\n", e.UsingOp.String())
	}

	if e.Nulls != nil {
		utils.IndentPrintf(indent+2, " nulls: %s\n", orderByNullsToString(*e.Nulls))
	}
}
