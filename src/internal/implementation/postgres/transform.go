package postgres

import (
	"fmt"
	"log"

	"github.com/antlr4-go/antlr/v4"
	"github.com/itsviktor/sir/src/internal/dsql"
	"github.com/itsviktor/sir/src/internal/parser"
	"github.com/itsviktor/sir/src/internal/transformer"
)

type tableRelation struct {
	schema string
	name   string
	alias  string
}

func Transform(query dsql.Query, domainName string) {
	fmt.Printf("transforming %s query:\n%s\n\n", domainName, query.SQL)

	input := antlr.NewInputStream(query.SQL)
	lexer := parser.NewPostgreSQLLexer(input)

	lexer.RemoveErrorListeners()
	lexer.AddErrorListener(&transformer.ErrorListener{})

	tokens := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

	p := parser.NewPostgreSQLParser(tokens)

	p.RemoveErrorListeners()
	p.AddErrorListener(&transformer.ErrorListener{})

	tree := p.Root()

	relations := make([]tableRelation, 0)

	transformer.WalkAntlrTree(tree, func(ctx antlr.ParserRuleContext) {
		// Getting relation info.
		tableRefCtx, ok := ctx.(*parser.Table_refContext)
		if ok {
			rel, err := parseTableRelation(tableRefCtx)
			if err != nil {
				log.Fatalf("transforming query: %v\n%s", err, query.SQL)
			}

			fmt.Printf("relation: %+v\n", rel)
			relations = append(relations, rel)

			return
		}

		cmp, ok := ctx.(*parser.Where_clauseContext)
		if ok {
			fmt.Printf("Where clause: %s\n", cmp.GetText())
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
