package transformer

import (
	"github.com/antlr4-go/antlr/v4"
	"github.com/itsviktor/sir/src/internal/utils"
)

type Context struct {
	Filepath  string
	StartLine int
}

func NewTransformContext(filepath string, startLine int) *Context {
	return &Context{
		StartLine: startLine,
		Filepath:  filepath,
	}
}

func (ctx *Context) TokenToPosition(token antlr.Token) utils.FilePosition {
	return utils.FilePosition{
		Filepath: ctx.Filepath,
		Line:     token.GetLine() + ctx.StartLine - 1,
		Column:   token.GetColumn(),
	}

}

func (ctx *Context) ErrOnToken(token antlr.Token, format string, v ...any) {
	utils.TraceErr(ctx.TokenToPosition(token), format, v...)
}
