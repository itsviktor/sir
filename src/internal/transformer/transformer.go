package transformer

import (
	"github.com/itsviktor/sir/src/internal/dsql"
	"github.com/itsviktor/sir/src/internal/ir"
	"github.com/itsviktor/sir/src/internal/schema"
)

type Transformer interface {
	// Transform transforms provided query to the Internal Representation.
	Transform(query dsql.Query, domainName string, tables map[string]*schema.Table) ir.Query
}
