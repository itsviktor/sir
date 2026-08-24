package ir

type SubqueryExpr struct {
	Query Query
}

func (SubqueryExpr) expr() {}
