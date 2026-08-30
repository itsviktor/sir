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

func (t PostgresTransformer) Transform(q dsql.Query, domainName string, tables map[string]*schema.Table) {
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
	transformer.WalkTree(tree, func(ctx antlr.Tree) {
		selectCtx, ok := ctx.(*parser.Select_clauseContext)
		if ok {
			nscope := newScope(selectCtx)

			if rootScope != nil {
				rootScope.addChildren(selectCtx, nscope)
			}
			rootScope = nscope

			return
		}

		tableCtx, ok := ctx.(*parser.Table_refContext)
		if ok {
			if rootScope == nil {
				transformCtx.ErrOnToken(tableCtx.GetStart(), "table reference in empty scope")
			}

			relation := parseRelation(tableCtx, rootScope, transformCtx)

			switch rel := relation.(type) {
			case *ir.TableRelation:
				addTableRelationToScope(tableCtx, rel, rootScope, tables, transformCtx)
			case *ir.JoinRelation:
				addTableRelationToScope(tableCtx, rel.Left.(*ir.TableRelation), rootScope, tables, transformCtx)
				addTableRelationToScope(tableCtx, rel.Right.(*ir.TableRelation), rootScope, tables, transformCtx)
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
	selectCtx, ok := transformer.FindFirstWide[*parser.Select_clauseContext](tree)
	if !ok {
		transformCtx.ErrOnToken(tree.GetStart(), "no select clause")
	}

	query := parseSelectQuery(selectCtx, rootScope, transformCtx)
	query.Print(0)
}
