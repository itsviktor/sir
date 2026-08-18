package postgres

import "github.com/itsviktor/sir/src/internal/schema"

type pgTable struct {
	NameC   string
	columns map[string]*pgColumn
}

func (t pgTable) Name() string {
	return t.NameC
}

func (t pgTable) HasColumn(name string) bool {
	_, ok := t.columns[name]
	return ok
}

func (t pgTable) GetColumn(name string) schema.Column {
	return t.columns[name]
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
