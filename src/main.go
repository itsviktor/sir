package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/antlr4-go/antlr/v4"
	"github.com/itsviktor/sir/src/internal/dsqlloader"
)

type errorListener struct {
	*antlr.DefaultErrorListener
}

func (l *errorListener) SyntaxError(
	recognizer antlr.Recognizer,
	offendingSymbol any,
	line int,
	column int,
	msg string,
	e antlr.RecognitionException,
) {
	fmt.Printf("error at %d:%d: %s\n", line, column, msg)
}

func drillDown(ctx antlr.ParserRuleContext) antlr.ParserRuleContext {
	for ctx.GetChildCount() == 1 {
		if child, ok := ctx.GetChild(0).(antlr.ParserRuleContext); ok {
			ctx = child
			continue
		}
		break
	}
	return ctx
}

func walk(node antlr.Tree, visit func(antlr.ParserRuleContext)) {
	if ctx, ok := node.(antlr.ParserRuleContext); ok {
		visit(ctx)
	}
	for i := 0; i < node.GetChildCount(); i++ {
		walk(node.GetChild(i), visit)
	}
}

func main() {
	// dsn := os.Getenv("SIR_DSN")
	// dsn := "postgres://postgres:postgres@localhost:5432/sir"

	// if dsn == "" {
	// 	log.Fatalf("empty connection string. Please, provide the DSN string using SIR_DSN environment variable")
	// }

	// // driver := os.Getenv("SIR_DRIVER")
	// driver := "pgx"

	// if driver == "" {
	// 	log.Fatalf("empty driver string. Please, provide the driver name using SIR_DRIVER environment variable")
	// }

	// dbType, err := dbtype.ParseAndValidate(driver)
	// if err != nil {
	// 	log.Fatalf("%v", err)
	// }
	// db, err := connection.Connect(driver, dsn)
	// if err != nil {
	// 	log.Fatalf("connection error: %v", err)
	// }

	// var tables []schema.Table
	// switch dbType {
	// case dbtype.TypePostgres:
	// 	tables, err = postgres.Inspect(db)
	// 	if err != nil {
	// 		log.Fatalf("inspecting postgres table structure: %v", err)
	// 	}
	// default:
	// 	log.Fatalf("unsupported database type: %s", dbType)
	// }

	// for _, table := range tables {
	// 	fmt.Printf("table: %s\n", table.Name())

	// 	for _, column := range table.Columns() {
	// 		fmt.Printf("column %s: db type: %s, go type: %s, default: %s, primaryKey: %t, nullable: %t\n", column.Name(), column.Type().DbName(), column.Type().GoName(), derefStr(column.DefaultValue(), "<empty>"), column.IsPrimaryKey(), column.IsNullable())
	// 	}

	// 	fmt.Printf("\n")
	// }

	// input := antlr.NewInputStream("select * from users where $data.name = users.name")
	// lexer := parser.NewPostgreSQLLexer(input)

	// lexer.RemoveErrorListeners()
	// lexer.AddErrorListener(&errorListener{})

	// tokens := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

	// p := parser.NewPostgreSQLParser(tokens)

	// p.RemoveErrorListeners()
	// p.AddErrorListener(&errorListener{})

	// tree := p.Root()

	// walk(tree, func(ctx antlr.ParserRuleContext) {
	// 	cmp, ok := ctx.(*parser.A_expr_compareContext)
	// 	if !ok {
	// 		return
	// 	}
	// 	if cmp.GetChildCount() != 3 {
	// 		return
	// 	}

	// 	leftChild := cmp.GetChild(0)
	// 	rightChild := cmp.GetChild(2)

	// 	left := drillDown(leftChild.(antlr.ParserRuleContext))
	// 	right := drillDown(rightChild.(antlr.ParserRuleContext))

	// 	fmt.Printf("compare: [%s] [%s] [%s]\n",
	// 		leftChild.(antlr.ParseTree).GetText(),
	// 		cmp.GetChild(1).(antlr.ParseTree).GetText(),
	// 		rightChild.(antlr.ParseTree).GetText(),
	// 	)
	// 	fmt.Printf("left type: %T, right type: %T\n", left, right)

	// 	right = right.(*parser.ColumnrefContext)
	// })

	domains, err := dsqlloader.Load("queries")
	if err != nil {
		log.Fatalf("parsing query files: %v", err)
	}

	for _, domain := range domains {
		fmt.Printf("domain: %+v\n", domain.Name)
		fmt.Printf("\n")

		for _, query := range domain.Queries {
			fmt.Printf("query: %s, kind: %s\n%s\n", query.Name, query.Kind, query.SQL)
			fmt.Printf("%s\n", strings.Repeat("-", 30))
		}

		fmt.Printf("\n\n")
	}
}

func derefStr(ptr *string, fallback string) string {
	if ptr != nil {
		return *ptr
	}
	return fallback
}
