package postgres

import (
	"github.com/antlr4-go/antlr/v4"
	"github.com/itsviktor/sir/src/internal/ir"
	"github.com/itsviktor/sir/src/internal/parser"
	"github.com/itsviktor/sir/src/internal/transformer"
)

// getRelationName parses and returns information about table relation - name and alias.
// It doesn't parse ON clause or anything else. Returns info about a left relation in a JOIN clause.
func getRelationName(ctx *parser.Table_refContext, tctx *transformer.Context) *ir.TableRelation {
	rel := &ir.TableRelation{}

	// Relation name parsing.
	first := ctx.GetChild(0)
	if first == nil {
		tctx.ErrOnToken(ctx.GetStart(), "invalid relation: first child is nil")
	}

	qualifiedName, ok := transformer.FindFirst[*parser.Qualified_nameContext](first)
	if !ok {
		tctx.ErrOnToken(ctx.GetStart(), "invalid relation: first child doesn't contain qualifiedName node")
	}

	parts := parseQualifiedParts(qualifiedName.Colid().(*parser.ColidContext), qualifiedName.Indirection())
	if len(parts) == 0 {
		tctx.ErrOnToken(ctx.GetStart(), "invalid relation: qualified parts size is nil")
	}

	rel.Name = sanitizeIdentifier(parts[0])

	// Alias parsing
	second := ctx.GetChild(1)
	aliasClause, ok := second.(*parser.Alias_clauseContext)
	if ok {
		colid := aliasClause.Colid()
		if colid != nil {
			rel.Alias = sanitizeIdentifier(colid.GetText())
		}
	}

	return rel
}

// parseRelation parses a table relation, including a chain of JOIN clauses.
func parseRelation(
	ctx *parser.Table_refContext,
	scope *transformer.Scope,
	tctx *transformer.Context,
) ir.Relation {
	children := ctx.GetChildren()

	var current ir.Relation
	current = getRelationName(ctx, tctx)

	for i := 0; i < len(children); i++ {
		if _, ok := children[i].(*parser.Join_typeContext); !ok {
			continue
		}

		var next int

		current, next = parseJoin(
			children,
			i,
			current,
			scope,
			tctx,
		)

		i = next - 1
	}

	return current
}

// parseJoin parses a single JOIN clause and returns the resulting relation
// and the index of the next JOIN clause.
func parseJoin(
	children []antlr.Tree,
	joinIndex int,
	left ir.Relation,
	scope *transformer.Scope,
	tctx *transformer.Context,
) (ir.Relation, int) {
	joinTypeCtx := children[joinIndex].(*parser.Join_typeContext)
	joinType := parseJoinType(joinTypeCtx, tctx)

	var (
		right ir.Relation
		on    ir.Expr
	)

	for i := joinIndex + 1; i < len(children); i++ {
		switch child := children[i].(type) {
		case *parser.Join_typeContext:
			return buildJoin(
				left,
				joinType,
				right,
				on,
				joinTypeCtx,
				tctx,
			), i

		case *parser.Table_refContext:
			if right != nil {
				tctx.ErrOnToken(
					child.GetStart(),
					"invalid JOIN: multiple right relations",
				)
				continue
			}

			right = getRelationName(child, tctx)

		case *parser.Join_qualContext:
			if on != nil {
				tctx.ErrOnToken(
					child.GetStart(),
					"invalid JOIN: multiple join qualifications",
				)
				continue
			}

			exprCtx, ok := transformer.FindFirst[*parser.A_exprContext](child)
			if !ok {
				tctx.ErrOnToken(
					child.GetStart(),
					"invalid JOIN qualification: A expression not found",
				)
				continue
			}

			on = parseExpr(exprCtx, scope, tctx)
		}
	}

	return buildJoin(
		left,
		joinType,
		right,
		on,
		joinTypeCtx,
		tctx,
	), len(children)
}

func buildJoin(
	left ir.Relation,
	joinType ir.JoinType,
	right ir.Relation,
	on ir.Expr,
	joinTypeCtx *parser.Join_typeContext,
	tctx *transformer.Context,
) ir.Relation {
	if right == nil {
		tctx.ErrOnToken(
			joinTypeCtx.GetStart(),
			"invalid JOIN: right relation not found",
		)
	}

	return &ir.JoinRelation{
		Left:  left,
		Type:  joinType,
		Right: right,
		On:    on,
	}
}

func parseJoinType(joinTypeCtx *parser.Join_typeContext, transformCtx *transformer.Context) ir.JoinType {
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

	return 0
}
