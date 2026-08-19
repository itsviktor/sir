package transformer

import (
	"github.com/antlr4-go/antlr/v4"
	"github.com/itsviktor/sir/src/internal/utils"
)

type TransformContext struct {
	Pos       utils.Position
	StartLine int
}

func NewTransformContext(filepath string, startLine int) *TransformContext {
	return &TransformContext{
		StartLine: startLine,
		Pos: utils.Position{
			Filepath: filepath,
			Line:     0,
			Column:   0,
		},
	}
}

func (ctx *TransformContext) PositionToToken(token antlr.Token) {
	ctx.Pos.Line = token.GetLine() + ctx.StartLine - 1
	ctx.Pos.Column = token.GetColumn()
}
