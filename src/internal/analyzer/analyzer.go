package analyzer

import "github.com/itsviktor/sir/src/internal/ir"

type Analyzer interface {
	TypeCheck(query ir.Query) error
}
