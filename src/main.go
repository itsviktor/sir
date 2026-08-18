package main

import (
	"database/sql"
	"log"

	"github.com/itsviktor/sir/src/internal/connection"
	"github.com/itsviktor/sir/src/internal/dbtype"
	"github.com/itsviktor/sir/src/internal/dsql"
	"github.com/itsviktor/sir/src/internal/implementation/postgres"
)

func main() {
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
	dbType, err := dbtype.ParseAndValidate(driver)
	if err != nil {
		log.Fatalf("%v", err)
	}
	db, err := connection.Connect(driver, dsn)
	if err != nil {
		log.Fatalf("connection error: %v", err)
	}

	switch dbType {
	case dbtype.TypePostgres:
		doPostgres(db)
	default:
		log.Fatalf("unsupported database type: %s", dbType)
	}
}

func doPostgres(db *sql.DB) {
	// Getting database structure.
	tables, err := postgres.Inspect(db)
	if err != nil {
		log.Fatalf("inspecting postgres table structure: %v", err)
	}

	_ = tables

	// Parsing user queries.
	domains, err := dsql.LoadDir("queries")
	if err != nil {
		log.Fatalf("parsing query files: %v", err)
	}

	// Transforming queries into internal representation.
	for _, domain := range domains {
		for _, query := range domain.Queries {
			postgres.Transform(query, domain.Name)
		}
	}
}
