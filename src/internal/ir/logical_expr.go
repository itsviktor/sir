package ir

type LogicalOp int

const (
	And LogicalOp = iota
	Or
)

type LogicalExpr struct {
	Left  Expr
	Op    LogicalOp
	Right Expr
}

func (LogicalExpr) expr() {}
