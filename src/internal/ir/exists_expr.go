package ir

type ExistsExpr struct {
	Query SelectQuery
	Not   bool
}

func (ExistsExpr) expr() {}
