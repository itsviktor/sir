package postgres

import (
	"strings"

	"github.com/antlr4-go/antlr/v4"
	"github.com/itsviktor/sir/src/internal/ir"
	"github.com/itsviktor/sir/src/internal/parser"
	"github.com/itsviktor/sir/src/internal/transformer"
)

func parseRelation(ctx *parser.Table_refContext, scope *scope, transformCtx *transformer.Context) ir.Relation {
	_, ok := transformer.FindFirstWide[*parser.Join_qualContext](ctx)
	if ok {
		return parseJoinRelation(ctx, scope, transformCtx)
	} else {
		return parseTableRelation(ctx, scope, transformCtx)
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
	rightTableCtx, ok := transformer.FindFirstWide[*parser.Table_refContext](ctx)
	if !ok {
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
	firstChild := relExprCtx.GetChild(0)
	if firstChild == nil {
		transformCtx.ErrOnToken(relExprCtx.GetStart(), "invalid relation expression context, first child is nil")
	}

	qualNameCtx, ok := firstChild.(*parser.Qualified_nameContext)
	if !ok {
		transformCtx.ErrOnToken(relExprCtx.GetStart(), "invalid relation expression context, wait for the first child to be qualified name context, got %T", firstChild)
	}

	switch qualNameCtx.GetChildCount() {
	case 1:
		nameCtx := qualNameCtx.GetChild(0)
		if nameCtx == nil {
			transformCtx.ErrOnToken(qualNameCtx.GetStart(), "invalid qualified name context, no first child")
		}

		colidCtx, ok := nameCtx.(*parser.ColidContext)
		if !ok {
			transformCtx.ErrOnToken(qualNameCtx.GetStart(), "invalid qualified name context, wait first children to be colid context, got %T", nameCtx)
		}

		return unquoteIdentifier(colidCtx.GetText()), nil
	case 2:
		nameCtx := qualNameCtx.GetChild(0)
		if nameCtx == nil {
			transformCtx.ErrOnToken(qualNameCtx.GetStart(), "invalid qualified name context, no first child")
		}

		colidCtx, ok := nameCtx.(*parser.ColidContext)
		if !ok {
			transformCtx.ErrOnToken(qualNameCtx.GetStart(), "invalid qualified name context, wait first children to be colid context, got %T", nameCtx)
		}

		schemaCtx := qualNameCtx.GetChild(1)
		if schemaCtx == nil {
			transformCtx.ErrOnToken(qualNameCtx.GetStart(), "invalid qualified name context, no second child")
		}

		colLabelCtx, ok := transformer.FindFirst[*parser.ColLabelContext](schemaCtx)
		if !ok {
			transformCtx.ErrOnToken(qualNameCtx.GetStart(), "invalid schema name, no col label context")
		}

		return unquoteIdentifier(colLabelCtx.GetText()), new(unquoteIdentifier(colidCtx.GetText()))
	}

	transformCtx.ErrOnToken(qualNameCtx.GetStart(), "invalid amount of qualified name context children")

	return "", nil
}

func unquoteIdentifier(s string) string {
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		return strings.ReplaceAll(s[1:len(s)-1], `""`, `"`)
	}

	return strings.ToLower(s)
}
