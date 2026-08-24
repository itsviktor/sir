package ir

type IsNullExpr struct {
	Expr Expr
	Not  bool
}

func (IsNullExpr) expr() {}
