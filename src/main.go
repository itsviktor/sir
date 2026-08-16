package main

import (
	"fmt"

	"github.com/antlr4-go/antlr/v4"
	"github.com/itsviktor/sir/src/internal/parser"
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

	input := antlr.NewInputStream("select * from users where name = $data.name or surname = $data.surname or id = $id")
	lexer := parser.NewPostgreSQLLexer(input)

	lexer.RemoveErrorListeners()
	lexer.AddErrorListener(&errorListener{})

	tokens := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

	parser := parser.NewPostgreSQLParser(tokens)

	parser.RemoveErrorListeners()
	parser.AddErrorListener(&errorListener{})

	// Запускаем entry rule PostgreSQL grammar.
	parser.Root()
}

func derefStr(ptr *string, fallback string) string {
	if ptr != nil {
		return *ptr
	}
	return fallback
}
