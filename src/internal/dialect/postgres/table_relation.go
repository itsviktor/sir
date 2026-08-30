package postgres

import (
	"github.com/antlr4-go/antlr/v4"
	"github.com/itsviktor/sir/src/internal/ir"
	"github.com/itsviktor/sir/src/internal/parser"
	"github.com/itsviktor/sir/src/internal/transformer"
)

func isTableRelation(ctx *parser.Table_refContext) bool {
	_, ok := transformer.FindFirstWide[*parser.Join_qualContext](ctx)
	return !ok
}

func parseRelation(ctx *parser.Table_refContext, scope *scope, transformCtx *transformer.Context) ir.Relation {
	if isTableRelation(ctx) {
		return parseTableRelation(ctx, scope, transformCtx)
	} else {
		return parseJoinRelation(ctx, scope, transformCtx)
	}
}

func parseJoinRelation(ctx *parser.Table_refContext, scope *scope, transformCtx *transformer.Context) *ir.JoinRelation {
	rel := &ir.JoinRelation{}

	// Parsing left table relation.
	firstChild := ctx.GetChild(0)
	if firstChild == nil {
		transformCtx.ErrOnToken(ctx.GetStart(), "invalid join, first child is nil")
	}
	relExprCtx, ok := firstChild.(*parser.Relation_exprContext)
	if !ok {
		transformCtx.ErrOnToken(ctx.GetStart(), "invalid join, wait first child to be relation expression, got %T", firstChild)
	}
	name, _ := parseSchemaAndName(relExprCtx, transformCtx)

	left := &ir.TableRelation{
		Name: name,
	}
	rel.Left = left

	// parsing right table relation.
	var rightTableCtx *parser.Table_refContext
	for _, child := range ctx.GetChildren() {
		rightTableCtx, _ = child.(*parser.Table_refContext)
		if rightTableCtx != nil {
			break
		}
	}
	if rightTableCtx == nil {
		transformCtx.ErrOnToken(ctx.GetStart(), "invalid join, no table ref for the right table")
	}

	right := parseTableRelation(rightTableCtx, scope, transformCtx)
	rel.Right = right

	// Parsing join type.
	joinType := ir.InnerJoin
	joinTypeCtx, ok := transformer.FindFirstWide[*parser.Join_typeContext](ctx)
	if ok {
		joinType = parseJoinType(joinTypeCtx, scope, transformCtx)
	}
	rel.Type = joinType

	// Parsing ON expression.
	exprCtx, ok := transformer.FindFirstWide[*parser.Join_qualContext](ctx)
	if !ok {
		transformCtx.ErrOnToken(ctx.GetStart(), "invalid join, no join qual context")
	}
	aExpr, ok := transformer.FindFirst[*parser.A_exprContext](exprCtx)
	if !ok {
		transformCtx.ErrOnToken(ctx.GetStart(), "invalid join qual ctx, no A expression context")
	}
	rel.On = parseExpr(aExpr, scope, transformCtx)

	return rel
}

func parseJoinType(joinTypeCtx *parser.Join_typeContext, scope *scope, transformCtx *transformer.Context) ir.JoinType {
	firstChild := joinTypeCtx.GetChild(0)
	if firstChild == nil {
		transformCtx.ErrOnToken(joinTypeCtx.GetStart(), "invalid join type, first child is nil")
	}

	termNode, ok := firstChild.(*antlr.TerminalNodeImpl)
	if !ok {
		transformCtx.ErrOnToken(joinTypeCtx.GetStart(), "invalid join type, wait first child to be terminal node, got %T", firstChild)
	}

	switch termNode.GetText() {
	case "LEFT":
		return ir.LeftJoin
	case "RIGHT":
		return ir.RightJoin
	case "INNER":
		return ir.InnerJoin
	case "FULL":
		return ir.FullJoin
	}

	transformCtx.ErrOnToken(joinTypeCtx.GetStart(), "invalid join type: %s", joinTypeCtx.GetText())

	return ir.InnerJoin
}

func parseTableRelation(ctx *parser.Table_refContext, scope *scope, transformCtx *transformer.Context) *ir.TableRelation {
	rel := &ir.TableRelation{}

	firstChild := ctx.GetChild(0)
	if firstChild == nil {
		transformCtx.ErrOnToken(ctx.GetStart(), "invalid ref context, first child is nil")
	}

	relExprCtx, ok := firstChild.(*parser.Relation_exprContext)
	if !ok {
		transformCtx.ErrOnToken(ctx.GetStart(), "invalid ref context, wait first child to be relation expression, got %T", firstChild)
	}

	name, _ := parseSchemaAndName(relExprCtx, transformCtx)
	rel.Name = name

	aliasCtx, ok := transformer.FindFirstWide[*parser.Alias_clauseContext](ctx)
	if ok {
		colidCtx, ok := transformer.FindFirst[*parser.ColidContext](aliasCtx)
		if !ok {
			transformCtx.ErrOnToken(aliasCtx.GetStart(), "invalid alias clause, no colid context")
		}

		rel.Alias = new(unquoteIdentifier(colidCtx.GetText()))
	}

	return rel
}

func parseSchemaAndName(relExprCtx *parser.Relation_exprContext, transformCtx *transformer.Context) (string, *string) {
	qualNameCtx, ok := relExprCtx.GetChild(0).(*parser.Qualified_nameContext)
	if !ok {
		transformCtx.ErrOnToken(relExprCtx.GetStart(), "expected qualified_name, got %T", relExprCtx.GetChild(0))
	}

	colid, ok := qualNameCtx.GetChild(0).(*parser.ColidContext)
	if !ok {
		transformCtx.ErrOnToken(qualNameCtx.GetStart(), "expected colid as first child, got %T", qualNameCtx.GetChild(0))
	}

	parts := parseQualifiedParts(colid, qualNameCtx.Indirection())

	switch len(parts) {
	case 1:
		return parts[0], nil
	case 2:
		return parts[1], &parts[0]
	default:
		transformCtx.ErrOnToken(qualNameCtx.GetStart(), "unexpected amount of qualified name parts: %d", len(parts))
		return "", nil
	}
}
