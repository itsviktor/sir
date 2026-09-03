package postgres

import (
	"bytes"
	"fmt"
	"go/format"
	"log"
	"os"
	"path"
	"strings"

	"github.com/danielgtaylor/casing"
	"github.com/itsviktor/sir/src/internal/codegen"
	"github.com/itsviktor/sir/src/internal/ir"
	"github.com/itsviktor/sir/src/internal/transformer"
)

type PostgresGoDomain struct {
	outDir               string
	name                 string
	repositoryStructName string
	b                    *bytes.Buffer
}

type PostgresGoCodegen struct {
	outDir string
}

func NewPostgresGoCodegen(outDir string) PostgresGoCodegen {
	c := PostgresGoCodegen{
		outDir,
	}

	c.createOutDir(outDir)

	return c
}

func (PostgresGoCodegen) createOutDir(outDir string) {
	err := os.MkdirAll(outDir, 0755)
	if err != nil {
		log.Fatalf("failed to create out dir: %v", err)
	}
}

func (PostgresGoCodegen) CreateDomain(outDir, name string) codegen.GoDomain {
	repositoryStructName := casing.Camel(name + "Repository")

	initialCode := `
		package db

		import "database/sql"

		type %rsn struct {
			db *sql.DB
		}

		func New%rsn(db *sql.DB) *%rsn {
			return &%rsn{db}
		}
	`
	initialCode = strings.ReplaceAll(initialCode, "%rsn", repositoryStructName)

	return &PostgresGoDomain{
		outDir:               outDir,
		name:                 name,
		repositoryStructName: repositoryStructName,
		b:                    bytes.NewBuffer([]byte(initialCode)),
	}
}

func (d *PostgresGoDomain) AddQuery(query ir.Query, name string, scope *transformer.Scope) error {
	switch q := query.(type) {
	case *ir.SelectQuery:
		return d.addSelectQuery(q, name, scope)
	}

	return fmt.Errorf("unsupported query type: %T", query)
}

func (d *PostgresGoDomain) addSelectQuery(q *ir.SelectQuery, name string, scope *transformer.Scope) error {

	return nil
}

func (d *PostgresGoDomain) Write() error {
	outFilename := path.Join(d.outDir, casing.Snake(d.name)+".go")

	f, err := os.Create(outFilename)
	if err != nil {
		return err
	}
	defer f.Close()

	sourceCode := d.b.Bytes()
	formatted, err := format.Source(sourceCode)
	if err != nil {
		return fmt.Errorf("failed to format: %v", err)
	}

	_, err = f.Write(formatted)
	return err
}
