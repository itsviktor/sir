package ir

type InExpr struct {
	Left  Expr
	Items []Expr
	Not   bool
}

func (InExpr) expr() {}
