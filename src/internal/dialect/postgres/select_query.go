package postgres

import (
	"github.com/itsviktor/sir/src/internal/ir"
	"github.com/itsviktor/sir/src/internal/parser"
	"github.com/itsviktor/sir/src/internal/transformer"
)

// parseSelectQuery parses SELECT clause context and returns Internal Representation of the query.
func parseSelectQuery(root *parser.Select_no_parensContext, scope *transformer.Scope, tctx *transformer.Context) *ir.SelectQuery {
	query := ir.NewSelectQuery()

	selectCtx := root.Select_clause()
	if selectCtx == nil {
		tctx.ErrOnToken(root.GetStart(), "invalid select query: no select clause context")
	}

	query.Returns = parseReturns(selectCtx.(*parser.Select_clauseContext), scope, tctx)
	query.Targets = parseTargets(selectCtx.(*parser.Select_clauseContext), scope, tctx)
	query.Where = parseWhere(selectCtx.(*parser.Select_clauseContext), scope, tctx)

	limitCtx := root.Select_limit()
	if limitCtx != nil {
		query.Offset = parseOffset(limitCtx.(*parser.Select_limitContext), scope, tctx)
		query.Limit = parseLimit(limitCtx.(*parser.Select_limitContext), scope, tctx)
	}

	sortCtx := root.Sort_clause_()
	if sortCtx != nil {
		query.OrderBy = parseOrderBy(sortCtx.(*parser.Sort_clause_Context), scope, tctx)
	}

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
	targetListCtx, ok := transformer.FindFirstWide[*parser.Target_listContext](selectCtx, 100)
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
	whereClause, ok := transformer.FindFirstWide[*parser.Where_clauseContext](selectCtx, 10)
	if !ok {
		return nil
	}

	return parseExpr(whereClause.A_expr().(*parser.A_exprContext), scope, tctx)
}

// parseOffset parses OFFSET clause of the query.
func parseOffset(limitCtx *parser.Select_limitContext, scope *transformer.Scope, tctx *transformer.Context) ir.Expr {
	offsetCtx := limitCtx.Offset_clause()
	if offsetCtx == nil {
		return nil
	}

	offsetValue := offsetCtx.Select_offset_value()
	if offsetValue != nil {
		aExpr := offsetValue.A_expr()
		return parseExpr(aExpr.(*parser.A_exprContext), scope, tctx)
	}

	fetchFirst := offsetCtx.Select_fetch_first_value()
	if fetchFirst != nil {
		cExpr := fetchFirst.C_expr()
		return parseCExpr(cExpr, scope, tctx)
	}

	return nil
}

// parseLimit parses LIMIT clause of the query.
func parseLimit(limitCtx *parser.Select_limitContext, scope *transformer.Scope, tctx *transformer.Context) ir.Expr {
	limitClause := limitCtx.Limit_clause()
	if limitClause == nil {
		return nil
	}

	limitValue := limitClause.Select_limit_value()
	if limitValue != nil {
		aExpr := limitValue.A_expr()
		if aExpr != nil {
			return parseExpr(aExpr.(*parser.A_exprContext), scope, tctx)
		}

		if limitClause.Select_limit_value().ALL() != nil {
			return &ir.LimitAllExpr{}
		}

	}

	return nil
}

// parseOrderBy parses ORDER BY clause of the query.
func parseOrderBy(sortCtx *parser.Sort_clause_Context, scope *transformer.Scope, tctx *transformer.Context) []*ir.OrderByExpr {
	sortClause := sortCtx.Sort_clause()
	if sortClause == nil {
		return nil
	}

	sortList := sortClause.Sortby_list()
	if sortList == nil {
		tctx.ErrOnToken(sortCtx.GetStart(), "invalid sort context: sort clause without sortby list")
	}

	var items []*ir.OrderByExpr
	for _, child := range sortList.AllSortby() {
		orderBy := &ir.OrderByExpr{}

		aExpr := child.A_expr()
		if aExpr == nil {
			tctx.ErrOnToken(child.GetStart(), "invalid order by item: no A expr")
		}
		orderBy.Expr = parseExpr(aExpr.(*parser.A_exprContext), scope, tctx)

		if child.USING() == nil {
			ascDescCtx := child.Asc_desc_()
			if ascDescCtx == nil {
				tctx.ErrOnToken(child.GetStart(), "invalid order by item: asc desc ctx is nil")
			}

			if ascDescCtx.ASC() != nil {
				orderBy.Order = ir.Asc
			} else {
				orderBy.Order = ir.Desc
			}
		} else {
			qualCtx := child.Qual_all_op()
			opText := qualCtx.GetText()
			opType, ok := ir.StringToOpType(opText)
			if !ok {
				tctx.ErrOnToken(qualCtx.GetStart(), "invalid order by USING: unknown operator %s", opText)
			}

			orderBy.UsingOp = new(ir.NewOp(opType, opText))
		}

		nullsCtx := child.Nulls_order_()
		if nullsCtx != nil {
			if nullsCtx.FIRST_P() != nil {
				orderBy.Nulls = new(ir.NullsFirst)
			} else {
				orderBy.Nulls = new(ir.NullsLast)
			}
		}

		items = append(items, orderBy)
	}

	return items
}
