package postgres

import (
	"github.com/itsviktor/sir/src/internal/ir"
	"github.com/itsviktor/sir/src/internal/parser"
	"github.com/itsviktor/sir/src/internal/transformer"
)

func parseSelectQuery(selectCtx *parser.Select_clauseContext, scope *scope, tctx *transformer.Context) *ir.SelectQuery {
	query := ir.NewSelectQuery()

	query.Targets = parseTargets(selectCtx, scope, tctx)

	return query
}

func parseTargets(selectCtx *parser.Select_clauseContext, scope *scope, tctx *transformer.Context) []ir.Relation {
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
