package ir

type UnaryOp int

const (
	Not UnaryOp = iota
)

type UnaryExpr struct {
	Op    UnaryOp
	Right Expr
}

func (UnaryExpr) expr() {}
