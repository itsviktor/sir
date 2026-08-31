package postgres

import "github.com/itsviktor/sir/src/internal/schema"

type pgEnumType struct {
	Name   string
	Values []string
}

func (t pgEnumType) DbName() string {
	return t.Name
}

func (t pgEnumType) GoName() string {
	return t.Name
}

func (t pgEnumType) Imports() []string {
	return []string{}
}

func (t pgEnumType) Kind() schema.TypeKind {
	return schema.Enum
}
