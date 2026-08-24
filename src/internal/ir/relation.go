package ir

type Relation interface {
	relation()
}

type TableRelation struct {
	Name  string
	Alias *string
}

func (TableRelation) relation() {}

type SubqueryRelation struct {
	Query Query
	Alias string
}

func (SubqueryRelation) relation() {}
