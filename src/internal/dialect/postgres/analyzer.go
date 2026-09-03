package postgres

import (
	"fmt"
	"strconv"
	"strings"

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

	if query.Offset != nil {
		tk, err := getExprTypeKind(query.Offset)

		if err != nil {
			return err
		}

		if tk != schema.Integer && tk != schema.Variable {
			return analyzer.NewErr(query.Offset.GetPos(), "invalid OFFSET clause return type: wait for Integer or Variable, got %s", schema.TypeKindToString(tk))
		}
	}

	if query.Limit != nil {
		tk, err := getExprTypeKind(query.Limit)

		if err != nil {
			return err
		}

		if tk != schema.Integer && tk != schema.Variable {
			return analyzer.NewErr(query.Offset.GetPos(), "invalid LIMIT clause return type: wait for Integer or Variable, got %s", schema.TypeKindToString(tk))
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
	case *ir.BinaryExpr:
		return getBinaryExprType(e)
	case *ir.DsqlExpr:
		return schema.Variable, nil
	case *ir.LimitAllExpr:
		return schema.Integer, nil
	case *ir.LogicalExpr:
		return getLogicalExprType(e)
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
	v := strings.ToLower(expr.Value)

	_, err := strconv.Atoi(v)
	if err == nil {
		return schema.Integer
	}

	_, err = strconv.ParseFloat(v, 64)
	if err == nil {
		return schema.Numeric
	}

	return schema.String
}

func getBinaryExprType(expr *ir.BinaryExpr) (schema.TypeKind, error) {
	left, err := getExprTypeKind(expr.Left)
	if err != nil {
		return 0, err
	}

	right, err := getExprTypeKind(expr.Right)
	if err != nil {
		return 0, err
	}

	if left == schema.Variable && right == schema.Variable {
		return 0, analyzer.NewErr(expr.GetPos(), "cannot perform binary operation between two variables")
	}

	if left == schema.Variable {
		return getVariableBinaryExprType(expr, right)
	}
	if right == schema.Variable {
		return getVariableBinaryExprType(expr, left)
	}

	op := expr.Op
	switch op.Type {
	case ir.Plus, ir.Minus, ir.Star, ir.Slash, ir.Caret, ir.Percent, ir.Gt, ir.Gte, ir.Lt, ir.Lte:
		if !schema.IsNumber(left) {
			return 0, analyzer.NewErr(expr.GetPos(), "left operand of the binary expression must be Number, got %s", schema.TypeKindToString(left))
		}
		if !schema.IsNumber(right) {
			return 0, analyzer.NewErr(expr.GetPos(), "right operand of the binary expression must be Number, got %s", schema.TypeKindToString(right))
		}

		return schema.Boolean, nil
	case ir.Like, ir.ILike, ir.NotLike, ir.NotILike:
		if left != schema.String {
			return 0, analyzer.NewErr(expr.GetPos(), "left operand of the LIKE expression must be String, got %s", schema.TypeKindToString(left))
		}
		if right != schema.String {
			return 0, analyzer.NewErr(expr.GetPos(), "right operand of the LIKE expression must be String, got %s", schema.TypeKindToString(right))
		}

		return schema.Boolean, nil
	case ir.Equal:
		if schema.IsNumber(left) && schema.IsNumber(right) {
			return schema.Boolean, nil
		}

		if left == right {
			return schema.Boolean, nil
		}

		return 0, analyzer.NewErr(expr.GetPos(), "cannot perform equal operation between the %s and %s", schema.TypeKindToString(left), schema.TypeKindToString(right))
	}

	return 0, analyzer.NewErr(expr.GetPos(), "unsupported binary expression operator: %s", expr.Op.String())
}

func getVariableBinaryExprType(expr *ir.BinaryExpr, second schema.TypeKind) (schema.TypeKind, error) {
	op := expr.Op.Type

	switch op {
	case ir.Plus, ir.Minus, ir.Star, ir.Slash, ir.Caret, ir.Percent, ir.Gt, ir.Gte, ir.Lt, ir.Lte:
		if schema.IsNumber(second) {
			return schema.Boolean, nil
		}

		return 0, analyzer.NewErr(expr.GetPos(), "unsupported operand for the binary expression: wait for Number, got %s", schema.TypeKindToString(second))
	case ir.Like, ir.ILike, ir.NotLike, ir.NotILike:
		if second == schema.String {
			return schema.Boolean, nil
		}

		return 0, analyzer.NewErr(expr.GetPos(), "unsupported operand for the like expression: wait for String, got %s", schema.TypeKindToString(second))
	case ir.Equal:
		return schema.Boolean, nil
	}

	return 0, analyzer.NewErr(expr.GetPos(), "unsupported binary expression operator: %s", expr.Op.String())
}

func getLogicalExprType(expr *ir.LogicalExpr) (schema.TypeKind, error) {
	left, err := getExprTypeKind(expr.Left)
	if err != nil {
		return 0, err
	}

	right, err := getExprTypeKind(expr.Right)
	if err != nil {
		return 0, err
	}

	if left == schema.Variable && right == schema.Variable {
		return 0, analyzer.NewErr(expr.GetPos(), "cannot perform logical operation between two variables")
	}

	switch right {
	case schema.Boolean, schema.Variable:
	default:
		return 0, analyzer.NewErr(expr.GetPos(), "right operand of the logical expression must be Boolean, got %s", schema.TypeKindToString(right))
	}

	switch left {
	case schema.Boolean, schema.Variable:
	default:
		return 0, analyzer.NewErr(expr.GetPos(), "left operand of the logical expression must be Boolean, got %s", schema.TypeKindToString(right))
	}

	return schema.Boolean, nil
}
