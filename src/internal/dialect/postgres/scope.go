package postgres

import (
	"github.com/antlr4-go/antlr/v4"
	"github.com/itsviktor/sir/src/internal/ir"
	"github.com/itsviktor/sir/src/internal/schema"
)

type tableRelation struct {
	relation *ir.TableRelation
	table    *schema.Table
}

type scope struct {
	parent         *scope
	nodeToChildren map[antlr.Tree]*scope

	node antlr.Tree

	relations map[string]*tableRelation
	aliases   map[string]*tableRelation
}

func newScope(node antlr.Tree) *scope {
	return &scope{
		parent:         nil,
		nodeToChildren: make(map[antlr.Tree]*scope),
		node:           node,
		relations:      make(map[string]*tableRelation),
		aliases:        make(map[string]*tableRelation),
	}
}

func (s *scope) addChildren(node antlr.Tree, other *scope) {
	s.nodeToChildren[node] = other
	other.parent = s
}

func (s *scope) childrenByNode(node antlr.Tree) (*scope, bool) {
	child, ok := s.nodeToChildren[node]
	return child, ok
}

func (s *scope) hasRelation(name string) bool {
	_, ok := s.relations[name]
	return ok
}

func (s *scope) relationNames() []string {
	var names []string
	if s.parent != nil {
		names = s.parent.relationNames()
	}

	for name, _ := range s.relations {
		names = append(names, name)
	}

	return names
}

func (s *scope) findRelation(name string) (*ir.TableRelation, *schema.Table, bool) {
	r, ok := s.relations[name]
	if ok {
		return r.relation, r.table, true
	}

	ar, ok := s.aliases[name]
	if ok {
		return ar.relation, ar.table, true
	}

	if s.parent != nil {
		return s.parent.findRelation(name)
	}

	return nil, nil, false
}

func (s *scope) findRelationCandidates(columnName string) []string {
	var candidates []string

	for name, data := range s.relations {
		if data.table.HasColumn(columnName) {
			candidates = append(candidates, name)
		}
	}

	if len(candidates) > 0 {
		return candidates
	}

	if s.parent != nil {
		return s.parent.findRelationCandidates(columnName)
	}

	return []string{}
}

func (s *scope) addRelation(relation *ir.TableRelation, table *schema.Table) {
	s.relations[relation.Name] = &tableRelation{
		relation,
		table,
	}
}

func (s *scope) hasAlias(alias string) bool {
	_, ok := s.aliases[alias]
	return ok
}

func (s *scope) addAlias(alias string, relation *ir.TableRelation, table *schema.Table) {
	s.aliases[alias] = &tableRelation{
		relation,
		table,
	}
}
