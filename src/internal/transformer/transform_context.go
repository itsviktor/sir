package transformer

import (
	"github.com/antlr4-go/antlr/v4"
	"github.com/itsviktor/sir/src/internal/utils"
)

type TransformContext struct {
	Filepath  string
	StartLine int
}

func NewTransformContext(filepath string, startLine int) *TransformContext {
	return &TransformContext{
		StartLine: startLine,
		Filepath:  filepath,
	}
}

func (ctx *TransformContext) TokenToPosition(token antlr.Token) utils.FilePosition {
	return utils.FilePosition{
		Filepath: ctx.Filepath,
		Line:     token.GetLine() + ctx.StartLine - 1,
		Column:   token.GetColumn(),
	}

}

func (ctx *TransformContext) ErrOnToken(token antlr.Token, format string, v ...any) {
	utils.TraceErr(ctx.TokenToPosition(token), format, v...)
}
