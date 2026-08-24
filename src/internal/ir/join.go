package ir

type JoinType int

const (
	LeftJoin JoinType = iota
	RightJoin
	InnerJoin
)

type Join struct {
	Type  JoinType
	Table Relation
	On    Expr
}
