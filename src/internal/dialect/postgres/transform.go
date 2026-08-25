package postgres

import (
	"log"

	"github.com/antlr4-go/antlr/v4"
	"github.com/itsviktor/sir/src/internal/dsql"
	"github.com/itsviktor/sir/src/internal/ir"
	"github.com/itsviktor/sir/src/internal/parser"
	"github.com/itsviktor/sir/src/internal/transformer"
	"github.com/itsviktor/sir/src/internal/utils"
)

type PostgresTransformer struct{}

func (t PostgresTransformer) Transform(q dsql.Query, domainName string) {
	utils.Debugf("transforming %s query:\n%s\n\n", domainName, q.SQL)
	transformCtx := transformer.NewTransformContext(q.File, q.StartLine)

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
	var rootScope *scope
	nodeToScope := map[antlr.Tree]*scope{}
	transformer.WalkTree(tree, func(ctx antlr.Tree) {
		selectCtx, ok := ctx.(*parser.Select_clauseContext)
		if ok {
			nscope := newScope(selectCtx)

			if rootScope == nil {
				rootScope = nscope
			} else {
				rootScope.add(nscope)
				rootScope = nscope
			}

			nodeToScope[ctx] = rootScope

			return
		}

		tableCtx, ok := ctx.(*parser.Table_refContext)
		if ok {
			if rootScope == nil {
				transformCtx.ErrOnToken(tableCtx.GetStart(), "table reference in empty scope")
			}

			rel := parseRelation(tableCtx, rootScope, transformCtx)

			tableRel, ok := rel.(*ir.TableRelation)
			if ok {
				rootScope.relations[tableRel.Name] = tableRel
				if tableRel.Alias != nil {
					alias := *tableRel.Alias
					if rootScope.hasAlias(alias) {
						transformCtx.ErrOnToken(tableCtx.GetStart(), "duplicate alias %s", alias)
					}

					rootScope.aliases[alias] = tableRel
				}
			}
		}
	}, func(ctx antlr.Tree) {
		_, ok := ctx.(*parser.Select_clauseContext)
		if ok {
			if rootScope == nil {
				log.Fatalf("exit from nil scope")
			}

			if rootScope.parent != nil {
				rootScope = rootScope.parent
			}
		}
	})

	// Second walk to build query.
	var query ir.Query
	transformer.WalkTree(tree, func(rawCtx antlr.Tree) {
		switch ctx := rawCtx.(type) {
		case *parser.Select_clauseContext:
			scope, ok := nodeToScope[ctx]
			if !ok {
				transformCtx.ErrOnToken(ctx.GetStart(), "select clause inside empty scope")
			}

			query = parseSelectQuery(ctx, scope, transformCtx)
		}
	}, func(ctx antlr.Tree) {})

	query.Print()
}
