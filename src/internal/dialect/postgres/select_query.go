package postgres

import (
	"fmt"

	"github.com/itsviktor/sir/src/internal/ir"
	"github.com/itsviktor/sir/src/internal/parser"
	"github.com/itsviktor/sir/src/internal/transformer"
)

// parseSelectQuery parses SELECT clause context and returns Internal Representation of the query.
func parseSelectQuery(selectCtx *parser.Select_clauseContext, scope *transformer.Scope, tctx *transformer.Context) *ir.SelectQuery {
	query := ir.NewSelectQuery()

	query.Returns = parseReturns(selectCtx, scope, tctx)
	query.Targets = parseTargets(selectCtx, scope, tctx)
	query.Where = parseWhere(selectCtx, scope, tctx)

	return query
}

// parseTargets parses FROM clause of the SELECT query.
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

// parseReturns parses return values of the SELECT query.
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
			continue
		}

		targetLabel, ok := child.(*parser.Target_labelContext)
		if ok {
			expr := parseExpr(targetLabel.A_expr().(*parser.A_exprContext), scope, tctx)
			targets = append(targets, expr.(ir.ReturnExpr))

			continue
		}
	}

	return targets
}

// parseWhere parses WHERE clause of the query.
func parseWhere(selectCtx *parser.Select_clauseContext, scope *transformer.Scope, tctx *transformer.Context) ir.Expr {
	whereClause, ok := transformer.FindFirstWide[*parser.Where_clauseContext](selectCtx, 1)
	if ok {
		fmt.Printf("FIND WHERE CLAUSE: %s\n", whereClause.GetText())
	}

	return nil
}
