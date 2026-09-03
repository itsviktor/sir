package main

import (
	"errors"
	"log"

	"github.com/itsviktor/sir/src/internal/analyzer"
	"github.com/itsviktor/sir/src/internal/codegen"
	"github.com/itsviktor/sir/src/internal/database"
	"github.com/itsviktor/sir/src/internal/dialect/postgres"
	"github.com/itsviktor/sir/src/internal/loader"
	"github.com/itsviktor/sir/src/internal/schema"
	"github.com/itsviktor/sir/src/internal/transformer"
	"github.com/itsviktor/sir/src/internal/utils"
)

func main() {
	utils.SetDebug(true)

	// dsn := os.Getenv("SIR_DSN")
	dsn := "postgres://postgres:postgres@localhost:5432/sir"

	if dsn == "" {
		log.Fatalf("empty connection string. Please, provide the DSN string using SIR_DSN environment variable")
	}

	// driver := os.Getenv("SIR_DRIVER")
	driver := "pgx"

	if driver == "" {
		log.Fatalf("empty driver string. Please, provide the driver name using SIR_DRIVER environment variable")
	}

	// Validating driver name.
	dialect, err := database.DialectFromDriver(driver)
	if err != nil {
		log.Fatalf("%v", err)
	}

	// Connecting to the database.
	db, err := database.Connect(driver, dsn)
	if err != nil {
		log.Fatalf("connection error: %v", err)
	}

	// Creating schema inspector.
	var inspector schema.Inspector
	switch dialect {
	case database.Postgres:
		inspector = &postgres.PostgresInspector{}
	default:
		log.Fatalf("unsupported database dialect: %s", dialect)
	}

	// Getting the database schema.
	tables, err := inspector.Inspect(db)
	if err != nil {
		log.Fatalf("inspecting database schema: %v", err)
	}

	// Loading user queries.
	domains, err := loader.LoadDir("queries")
	if err != nil {
		log.Fatalf("loading query files: %v", err)
	}

	// Creating IR transformer, analyzer and code generator.
	var t transformer.Transformer
	var a analyzer.Analyzer
	var cg codegen.GoCodegen
	switch dialect {
	case database.Postgres:
		t = &postgres.PostgresTransformer{}
		a = &postgres.PostgresAnalyzer{}
		cg = postgres.NewPostgresGoCodegen("out")
	default:
		log.Fatalf("unsupported database dialect: %s", dialect)
	}

	// Transforming queries to internal representation, then typechecking them and generate a code.
	for _, domain := range domains {
		d := cg.CreateDomain("out", domain.Name)

		for _, query := range domain.Queries {
			queryScope, queryInternalRepresentation := t.Transform(query, domain.Name, tables)

			if err := a.TypeCheck(queryInternalRepresentation); err != nil {
				if typeErr, ok := errors.AsType[*analyzer.AnalyzerError](err); ok {
					utils.TraceErr(query.Filepath, typeErr.Pos, "%s", typeErr.Error())
				} else {
					log.Fatalf("typecheck error: %s", err.Error())
				}
			}

			err := d.AddQuery(queryInternalRepresentation, query.Name, queryScope)
			if err != nil {
				log.Fatalf("failed to generate query code: %v", err)
			}
		}

		if err := d.Write(); err != nil {
			log.Fatalf("failed to write domain code into the file: %v", err)
		}
	}
}
