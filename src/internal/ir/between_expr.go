package ir

type BetweenExpr struct {
	Left Expr
	From Expr
	To   Expr
	Not  bool
}

func (BetweenExpr) expr() {}
