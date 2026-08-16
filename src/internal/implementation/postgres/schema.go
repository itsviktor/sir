package postgres

import "github.com/itsviktor/sir/src/internal/schema"

type pgTable struct {
	NameC   string
	columns []pgColumn
}

func (t pgTable) Name() string {
	return t.NameC
}

func (t pgTable) Columns() []schema.Column {
	columns := make([]schema.Column, len(t.columns))

	for i, column := range t.columns {
		columns[i] = column
	}

	return columns
}

type pgColumn struct {
	NameC         string
	DefaultValueC *string
	DbTypeC       string
	IsNullableC   bool
	IsPrimaryKeyC bool
	t             schema.ColumnType
}

func (c pgColumn) Name() string {
	return c.NameC
}

func (c pgColumn) DefaultValue() *string {
	return c.DefaultValueC
}

func (c pgColumn) Type() schema.ColumnType {
	return c.t
}

func (c pgColumn) IsNullable() bool {
	return c.IsNullableC
}

func (c pgColumn) IsPrimaryKey() bool {
	return c.IsPrimaryKeyC
}

type pgEnumType struct {
	NameC  string
	Values []string
}

func (t pgEnumType) DbName() string {
	return t.NameC
}

func (t pgEnumType) GoName() string {
	return "string"
}

func (t pgEnumType) Imports() []string {
	return []string{}
}
