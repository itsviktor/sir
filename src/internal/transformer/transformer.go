package transformer

import (
	"github.com/itsviktor/sir/src/internal/dsql"
)

type Transformer interface {
	// Transform transforms provided query to the Internal Representation.
	Transform(query dsql.Query, domainName string)
}
