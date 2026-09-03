package ir

import (
	"fmt"
	"strings"

	"github.com/itsviktor/sir/src/internal/schema"
)

type Relation interface {
	relation()
	Print(indent int)
	GetName() string
}

type TableRelation struct {
	Name  string
	Alias string
	table *schema.Table
}

func (TableRelation) relation() {}

func (r *TableRelation) LinkTable(table *schema.Table) {
	r.table = table
}

func (r *TableRelation) HasColumn(name string) bool {
	if r.table == nil {
		panic("cannot check column name: relation has no linked table")
	}

	return r.table.HasColumn(name)
}

func (r *TableRelation) GetColumn(name string) (schema.Column, bool) {
	column, ok := r.table.Columns[name]
	return column, ok
}

func (r *TableRelation) ColumnNames() []string {
	if r.table == nil {
		panic("cannot get column names: relation has no linked table")
	}

	return r.table.ColumnNames()
}

func (r *TableRelation) Print(indent int) {
	i := strings.Repeat(" ", indent)

	fmt.Printf("%s- table relation: name=%s ", i, r.Name)
	if r.Alias != "" {
		fmt.Printf("alias=%s", r.Alias)
	}
	fmt.Printf("\n")
}

func (r *TableRelation) GetName() string {
	return r.Name
}

type SubqueryRelation struct {
	Query Query
	Alias string
}

func (SubqueryRelation) relation() {}

func (SubqueryRelation) Print(indent int) {}

type JoinType int

const (
	LeftJoin JoinType = iota
	RightJoin
	InnerJoin
	FullJoin
)

func joinTypeToString(t JoinType) string {
	switch t {
	case LeftJoin:
		return "LEFT JOIN"
	case RightJoin:
		return "RIGHT JOIN"
	case InnerJoin:
		return "INNER JOIN"
	case FullJoin:
		return "FULL JOIN"
	}

	return "UNKNOWN JOIN"
}

type JoinRelation struct {
	Left  Relation
	Type  JoinType
	Right Relation
	On    Expr
}

func (JoinRelation) relation() {}

func (r *JoinRelation) Print(indent int) {
	i := strings.Repeat(" ", indent)

	fmt.Printf("%s- join relation, type: %s\n%s   left:\n", i, joinTypeToString(r.Type), i)
	r.Left.Print(indent + 4)
	fmt.Printf("%s   right:\n", i)
	r.Right.Print(indent + 4)

	if r.On != nil {
		fmt.Printf("%s   on:\n", i)
		r.On.Print(indent + 4)
	}
}

func (r *JoinRelation) GetName() string {
	return r.Left.GetName()
}
