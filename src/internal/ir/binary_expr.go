package ir

type BinaryOp int

const (
	Equal BinaryOp = iota
	NotEqual
	Gt
	Gte
	Lt
	Lte
	Like
	ILike
	NotLike
	NotILike
)

type BinaryExpr struct {
	Right Expr
	Op    BinaryOp
	Left  Expr
}

func (BinaryExpr) expr() {}
