package analyzer

import (
	"fmt"

	"github.com/itsviktor/sir/src/internal/utils"
)

type AnalyzerError struct {
	Pos utils.Position
	Msg string
}

func (e *AnalyzerError) Error() string {
	return e.Msg
}

func NewErr(pos utils.Position, format string, v ...any) *AnalyzerError {
	return &AnalyzerError{
		Pos: pos,
		Msg: fmt.Sprintf(format, v...),
	}
}
