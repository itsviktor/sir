package postgres

import (
	"fmt"
	"strconv"

	"github.com/itsviktor/sir/src/internal/analyzer"
	"github.com/itsviktor/sir/src/internal/ir"
	"github.com/itsviktor/sir/src/internal/schema"
)

type PostgresAnalyzer struct {
}

func (PostgresAnalyzer) TypeCheck(query ir.Query) error {
	switch q := query.(type) {
	case *ir.SelectQuery:
		return validateSelectQuery(q)
	default:
		return fmt.Errorf("unknown query type: %T", query)
	}
}

func validateSelectQuery(query *ir.SelectQuery) error {
	if query.Where != nil {
		tk, err := getExprTypeKind(query.Where)

		if err != nil {
			return err
		}

		if tk != schema.Boolean {
			return analyzer.NewErr(query.Where.GetPos(), "invalid WHERE clause return type: wait for Boolean, got %s", schema.TypeKindToString(tk))
		}
	}

	return nil
}

func getExprTypeKind(expr ir.Expr) (schema.TypeKind, error) {
	switch e := expr.(type) {
	case *ir.LiteralExpr:
		return getLiteralType(e), nil
	case *ir.IsExpr:
		return getIsType(e)
	case *ir.ColumnExpr:
		return getColumnType(e)
	default:
		return 0, analyzer.NewErr(expr.GetPos(), "unexpected expression type: %T", expr)
	}
}

func getIsType(expr *ir.IsExpr) (schema.TypeKind, error) {
	switch expr.Predicate {
	case ir.IsTrue, ir.IsFalse, ir.IsUnknown:
		tk, err := getExprTypeKind(expr.Expr)

		if err != nil {
			return 0, err
		}

		if tk != schema.Boolean {
			return 0, analyzer.NewErr(expr.GetPos(), "invalid IS clause operand type: wait for Boolean, got %s", schema.TypeKindToString(tk))
		}

		return tk, nil
	}

	return schema.Boolean, nil
}

func getColumnType(expr *ir.ColumnExpr) (schema.TypeKind, error) {
	columnName := expr.Name
	column, ok := expr.Relation.GetColumn(columnName)
	if !ok {
		return 0, analyzer.NewErr(expr.GetPos(), "column %q does not found in the %q relation", columnName, expr.Relation.Name)
	}

	return column.Type.Kind(), nil
}

func getLiteralType(expr *ir.LiteralExpr) schema.TypeKind {
	_, err := strconv.ParseFloat(expr.Value, 64)
	if err == nil {
		return schema.Numeric
	}

	return schema.String
}
