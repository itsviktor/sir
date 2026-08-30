package postgres

import (
	"log"

	"github.com/antlr4-go/antlr/v4"
	"github.com/itsviktor/sir/src/internal/dsql"
	"github.com/itsviktor/sir/src/internal/ir"
	"github.com/itsviktor/sir/src/internal/parser"
	"github.com/itsviktor/sir/src/internal/schema"
	"github.com/itsviktor/sir/src/internal/transformer"
	"github.com/itsviktor/sir/src/internal/utils"
)

type PostgresTransformer struct{}

func (t PostgresTransformer) Transform(q dsql.Query, domainName string, tables map[string]*schema.Table) ir.Query {
	utils.Debugf("transforming %s query:\n%s\n\n", domainName, q.SQL)
	tctx := transformer.NewTransformContext(q.File, q.StartLine)

	input := antlr.NewInputStream(q.SQL)
	lexer := parser.NewPostgreSQLLexer(input)

	lexer.RemoveErrorListeners()
	lexer.AddErrorListener(&transformer.ErrorListener{})

	tokens := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

	p := parser.NewPostgreSQLParser(tokens)

	p.RemoveErrorListeners()
	p.AddErrorListener(&transformer.ErrorListener{})

	tree := p.Root()

	// First walk to build scopes tree.
	var rootScope *transformer.Scope
	transformer.WalkTree(tree, func(ctx antlr.Tree) {
		selectCtx, ok := ctx.(*parser.Select_no_parensContext)
		if ok {
			scope := transformer.NewScope(selectCtx)

			if rootScope != nil {
				rootScope.AddChild(scope)
			}
			rootScope = scope

			return
		}

		tableCtx, ok := ctx.(*parser.Table_refContext)
		if ok {
			if rootScope == nil {
				tctx.ErrOnToken(tableCtx.GetStart(), "table reference in nil scope")
			}

			// Parsing relation name for the table ref context.
			relation := getRelationName(tableCtx, tctx)

			// Finding table schema for that relation.
			table, ok := tables[relation.Name]
			if !ok {
				tctx.ErrOnToken(tableCtx.GetStart(), "cannot find table for the relation %q", relation.Name)
			}

			// Adding relation to the scope.
			err := rootScope.AddRelation(relation, table)
			if err != nil {
				tctx.ErrOnToken(tableCtx.GetStart(), "%v", err)
			}
		}
	}, func(ctx antlr.Tree) {
		_, ok := ctx.(*parser.Select_no_parensContext)
		if ok {
			if rootScope == nil {
				log.Fatalf("exit from nil scope")
			}

			parentScope := rootScope.Parent()
			if parentScope != nil {
				rootScope = parentScope
			}
		}
	})

	// Building query
	selectCtx, ok := transformer.FindFirstWide[*parser.Select_no_parensContext](tree, 100)
	if !ok {
		tctx.ErrOnToken(tree.GetStart(), "no select clause")
	}

	query := parseSelectQuery(selectCtx, rootScope, tctx)

	return query
}
