package codegen

import (
	"github.com/itsviktor/sir/src/internal/ir"
	"github.com/itsviktor/sir/src/internal/transformer"
)

type GoDomain interface {
	AddQuery(query ir.Query, name string, scope *transformer.Scope) error
	Write() error
}

type GoCodegen interface {
	CreateDomain(outDir, name string) GoDomain
}
