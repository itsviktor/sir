package postgres

import (
	"fmt"
	"strings"

	"github.com/antlr4-go/antlr/v4"
	"github.com/itsviktor/sir/src/internal/dsql"
	"github.com/itsviktor/sir/src/internal/parser"
	"github.com/itsviktor/sir/src/internal/schema"
	"github.com/itsviktor/sir/src/internal/transformer"
	"github.com/itsviktor/sir/src/internal/utils"
)

type tableRelation struct {
	schema string
	name   string
	alias  string
}

type scope struct {
	parent    *scope
	children  []*scope
	relations []string
	aliases   map[string]string
}

func newScope() *scope {
	return &scope{
		parent:    nil,
		children:  make([]*scope, 0),
		relations: make([]string, 0),
		aliases:   make(map[string]string),
	}
}

func (s *scope) addChild(children *scope) {
	s.children = append(s.children, children)
	children.parent = s
}

func (s *scope) printRoot(offset int) {
	utils.Debugf("%sscope relations: \n", strings.Repeat(" ", offset))
	for _, rel := range s.relations {
		utils.Debugf("%s- %s\n", strings.Repeat(" ", offset), rel)
	}
	utils.Debugf("%sscope aliases: \n", strings.Repeat(" ", offset))
	for alias, rel := range s.aliases {
		utils.Debugf("%s- %s -> %s\n", strings.Repeat(" ", offset), alias, rel)
	}

	for _, child := range s.children {
		child.printRoot(offset + 4)
	}
}

type QueryPart struct {
	SQL string
}

type QueryWhereClause struct {
}

func Transform(query dsql.Query, tables map[string]schema.Table, domainName string) {
	utils.Debugf("transforming %s query:\n%s\n\n", domainName, query.SQL)
	transformCtx := transformer.NewTransformContext(query.File, query.StartLine)

	input := antlr.NewInputStream(query.SQL)
	lexer := parser.NewPostgreSQLLexer(input)

	lexer.RemoveErrorListeners()
	lexer.AddErrorListener(&transformer.ErrorListener{})

	tokens := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

	p := parser.NewPostgreSQLParser(tokens)

	p.RemoveErrorListeners()
	p.AddErrorListener(&transformer.ErrorListener{})

	tree := p.Root()

	// First traversal to build scopes.
	utils.Debug("DEBUG create scopes\n")
	var queryScope *scope
	scopeByContext := map[*parser.Select_clauseContext]*scope{}
	transformer.WalkAntlrTree(tree, func(ctx antlr.Tree) {
		selectContext, ok := ctx.(*parser.Select_clauseContext)
		if ok {
			utils.Debugf("create new scope\n")

			nScope := newScope()
			if queryScope != nil {
				queryScope.addChild(nScope)
			}
			queryScope = nScope

			scopeByContext[selectContext] = nScope

			return
		}

		tableRefCtx, ok := ctx.(*parser.Table_refContext)
		if ok {
			if queryScope == nil {
				transformCtx.PositionToToken(tableRefCtx.GetStart())
				utils.TraceErr(transformCtx.Pos, "table ref inside nil scope")
			}

			rel, err := parseTableRelation(tableRefCtx)
			if err != nil {
				transformCtx.PositionToToken(tableRefCtx.GetStart())
				utils.TraceErr(transformCtx.Pos, "transforming query: %v", err)
			}

			utils.Debugf("relation: %+v\n", rel)

			if rel.alias == "" {
				queryScope.relations = append(queryScope.relations, rel.name)
			} else {
				queryScope.aliases[rel.alias] = rel.name
			}

			return
		}
	}, func(ctx antlr.Tree) {
		_, ok := ctx.(*parser.Select_clauseContext)
		if ok {
			parentScope := queryScope.parent
			if parentScope != nil {
				utils.Debugf("exit from scope\n")
				queryScope = parentScope
			}
		}
	})

	utils.Debugf("\nDEBUG scopes\n")
	queryScope.printRoot(0)

	// Second traversal to analyze dsql tokens.
	utils.Debugf("\nDEBUG where analyze\n")
	var currentScope *scope
	transformer.WalkAntlrTree(tree, func(ctx antlr.Tree) {
		selectContext, ok := ctx.(*parser.Select_clauseContext)
		if ok {
			utils.Debugf("enter scope\n")

			scope, ok := scopeByContext[selectContext]
			if !ok {
				transformCtx.PositionToToken(selectContext.GetStart())
				utils.TraceErr(transformCtx.Pos, "scope not found for the select context")
			}
			currentScope = scope

			return
		}

		whereCtx, ok := ctx.(*parser.Where_clauseContext)
		if ok {
			analyzeWhere(whereCtx, currentScope, transformCtx)
		}
	}, func(ctx antlr.Tree) {
		_, ok := ctx.(*parser.Select_clauseContext)
		if ok {
			utils.Debugf("exit scope\n")
			currentScope = currentScope.parent
		}
	})
}

func parseTableRelation(ctx *parser.Table_refContext) (tableRelation, error) {
	var rel tableRelation

	aliasCtx := ctx.Alias_clause()
	if aliasCtx != nil {
		colidCtx := aliasCtx.Colid()
		if colidCtx == nil {
			return rel, fmt.Errorf("failed to find colid context in the alias context")
		}

		rel.alias = colidCtx.GetText()
	}

	relCtx := ctx.Relation_expr()
	if relCtx == nil {
		return rel, fmt.Errorf("failed to find relation expression in the query")
	}

	qnameCtx := relCtx.Qualified_name()
	if qnameCtx == nil {
		return rel, fmt.Errorf("failed to find qualified name context in the query")
	}

	qnameChildren := qnameCtx.GetChildren()

	switch len(qnameChildren) {
	case 1:
		child := qnameChildren[0]

		colidCtx, ok := child.(*parser.ColidContext)
		if !ok {
			return rel, fmt.Errorf(
				"expected ColidContext, got %T",
				child,
			)
		}

		rel.name = colidCtx.GetText()
	case 2:
		child := qnameChildren[0]

		colidCtx, ok := child.(*parser.ColidContext)
		if !ok {
			return rel, fmt.Errorf(
				"expected ColidContext, got %T",
				child,
			)
		}

		rel.schema = colidCtx.GetText()

		child = qnameChildren[1]

		indirCtx, ok := child.(*parser.IndirectionContext)
		if !ok {
			return rel, fmt.Errorf("mismatched qname children amount and child type; wait for IndirectionContext, got %T", child)
		}

		indirEl := indirCtx.Indirection_el(0)
		if indirEl == nil {
			return rel, fmt.Errorf("failed to parse indirection element child")
		}

		attrName := indirEl.Attr_name()
		if attrName == nil {
			return rel, fmt.Errorf("failed to parse attribute name from the indirection element")
		}

		rel.name = attrName.GetText()
	default:
		return rel, fmt.Errorf("invalid qname children amount %d", len(qnameChildren))
	}

	return rel, nil
}

func analyzeWhere(whereCtx *parser.Where_clauseContext, scope *scope, transformCtx *transformer.TransformContext) {
	utils.Debugf("where ctx: %s\n", whereCtx.GetText())

	expr := whereCtx.A_expr()
	if expr == nil {
		transformCtx.PositionToToken(whereCtx.GetStart())
		utils.TraceErr(transformCtx.Pos, "where condition without expression (probably an ANTLR issue)")
	}

}
