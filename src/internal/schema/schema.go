package schema

import (
	"maps"
	"slices"
)

type Table struct {
	Name    string
	Columns map[string]Column
}

func (t Table) HasColumn(name string) bool {
	_, ok := t.Columns[name]
	return ok
}

func (t Table) ColumnNames() []string {
	return slices.Collect(maps.Keys(t.Columns))
}

type Column struct {
	Name         string
	DefaultValue *string
	Type         ColumnType
	IsNullable   bool
	IsPrimaryKey bool
}

type ColumnType interface {
	DbName() string
	GoName() string
	Kind() TypeKind
	Imports() []string
}

type DefaultType struct {
	dbName  string
	goName  string
	kind    TypeKind
	imports []string
}

func NewDefaultType(dbName, goName string, kind TypeKind, imports []string) ColumnType {
	return &DefaultType{dbName: dbName, goName: goName, kind: kind, imports: imports}
}

func (t DefaultType) DbName() string {
	return t.dbName
}

func (t DefaultType) GoName() string {
	return t.goName
}

func (t DefaultType) Imports() []string {
	return t.imports
}

func (t DefaultType) Kind() TypeKind {
	return t.kind
}
