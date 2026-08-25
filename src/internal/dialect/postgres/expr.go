package postgres

import (
	"fmt"

	"github.com/itsviktor/sir/src/internal/ir"
	"github.com/itsviktor/sir/src/internal/parser"
	"github.com/itsviktor/sir/src/internal/transformer"
)

func parseExpr(aExpr *parser.A_exprContext, scope *scope, tctx *transformer.Context) ir.Expr {
	fmt.Printf("PARSE EXPRESSION: %s\n", aExpr.GetText())

	return parseQual(aExpr.A_expr_qual().(*parser.A_expr_qualContext), scope, tctx)
}

func parseQual(
	aExpr *parser.A_expr_qualContext,
	scope *scope,
	tctx *transformer.Context,
) ir.Expr {
	return parseLessLess(
		aExpr.A_expr_lessless().(*parser.A_expr_lesslessContext),
		scope,
		tctx,
	)
}

func parseLessLess(
	aExpr *parser.A_expr_lesslessContext,
	scope *scope,
	tctx *transformer.Context,
) ir.Expr {
	return parseOr(
		aExpr.A_expr_or(0).(*parser.A_expr_orContext),
		scope,
		tctx,
	)
}

func parseOr(
	aExpr *parser.A_expr_orContext,
	scope *scope,
	tctx *transformer.Context,
) ir.Expr {
	return parseAnd(
		aExpr.A_expr_and(0).(*parser.A_expr_andContext),
		scope,
		tctx,
	)
}

func parseAnd(
	aExpr *parser.A_expr_andContext,
	scope *scope,
	tctx *transformer.Context,
) ir.Expr {
	return parseBetween(
		aExpr.A_expr_between(0).(*parser.A_expr_betweenContext),
		scope,
		tctx,
	)
}

func parseBetween(
	aExpr *parser.A_expr_betweenContext,
	scope *scope,
	tctx *transformer.Context,
) ir.Expr {
	return parseIn(
		aExpr.A_expr_in(0).(*parser.A_expr_inContext),
		scope,
		tctx,
	)
}

func parseIn(
	aExpr *parser.A_expr_inContext,
	scope *scope,
	tctx *transformer.Context,
) ir.Expr {
	return parseUnaryNot(
		aExpr.A_expr_unary_not().(*parser.A_expr_unary_notContext),
		scope,
		tctx,
	)
}

func parseUnaryNot(
	aExpr *parser.A_expr_unary_notContext,
	scope *scope,
	tctx *transformer.Context,
) ir.Expr {
	return parseIsNull(
		aExpr.A_expr_isnull().(*parser.A_expr_isnullContext),
		scope,
		tctx,
	)
}

func parseIsNull(
	aExpr *parser.A_expr_isnullContext,
	scope *scope,
	tctx *transformer.Context,
) ir.Expr {
	return parseIsNot(
		aExpr.A_expr_is_not().(*parser.A_expr_is_notContext),
		scope,
		tctx,
	)
}

func parseIsNot(
	aExpr *parser.A_expr_is_notContext,
	scope *scope,
	tctx *transformer.Context,
) ir.Expr {
	return parseCompare(
		aExpr.A_expr_compare().(*parser.A_expr_compareContext),
		scope,
		tctx,
	)
}

func parseCompare(
	aExpr *parser.A_expr_compareContext,
	scope *scope,
	tctx *transformer.Context,
) ir.Expr {
	return &ir.LiteralExpr{
		Value: "a",
	}
}
