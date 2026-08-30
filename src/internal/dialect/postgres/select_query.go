package postgres

import (
	"github.com/itsviktor/sir/src/internal/ir"
	"github.com/itsviktor/sir/src/internal/parser"
	"github.com/itsviktor/sir/src/internal/transformer"
)

func parseSelectQuery(selectCtx *parser.Select_clauseContext, scope *transformer.Scope, tctx *transformer.Context) *ir.SelectQuery {
	query := ir.NewSelectQuery()

	query.Returns = parseReturns(selectCtx, scope, tctx)
	query.Targets = parseTargets(selectCtx, scope, tctx)

	return query
}

func parseTargets(selectCtx *parser.Select_clauseContext, scope *transformer.Scope, tctx *transformer.Context) []ir.Relation {
	relations := make([]ir.Relation, 0)

	fromClause, ok := transformer.FindFirst[*parser.From_listContext](selectCtx)
	if !ok {
		tctx.ErrOnToken(selectCtx.GetStart(), "no from clause in the select query")
	}

	for _, child := range fromClause.GetChildren() {
		tableRefCtx, ok := child.(*parser.Table_refContext)
		if !ok {
			continue
		}

		rel := parseRelation(tableRefCtx, scope, tctx)
		relations = append(relations, rel)
	}

	return relations
}

func parseReturns(selectCtx *parser.Select_clauseContext, scope *transformer.Scope, tctx *transformer.Context) []ir.ReturnExpr {
	targetListCtx, ok := transformer.FindFirstWide[*parser.Target_listContext](selectCtx)
	if !ok {
		tctx.ErrOnToken(selectCtx.GetStart(), "no target list in the select query")
	}

	var targets []ir.ReturnExpr
	for _, child := range targetListCtx.GetChildren() {
		_, ok := child.(*parser.Target_starContext)
		if ok {
			targets = append(targets, &ir.AllReturnExpr{})
		}

		targetLabel, ok := child.(*parser.Target_labelContext)
		if ok {
			aExpr := targetLabel.A_expr()
			expr := parseExpr(aExpr.(*parser.A_exprContext), scope, tctx)

			columnRefExpr, ok := expr.(*ir.ColumnExpr)
			if !ok {
				tctx.ErrOnToken(aExpr.GetStart(), "invalid SELECT clause: wait for ColumnRef, got %T", expr)
			}
			targets = append(targets, columnRefExpr)
		}
	}

	return targets
}
