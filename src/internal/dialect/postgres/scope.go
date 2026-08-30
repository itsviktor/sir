package postgres

import (
	"github.com/antlr4-go/antlr/v4"
	"github.com/itsviktor/sir/src/internal/ir"
)

type scope struct {
	parent         *scope
	nodeToChildren map[antlr.Tree]*scope

	node antlr.Tree

	relations map[string]*ir.TableRelation
	aliases   map[string]*ir.TableRelation
}

func newScope(node antlr.Tree) *scope {
	return &scope{
		parent:         nil,
		nodeToChildren: make(map[antlr.Tree]*scope),
		node:           node,
		relations:      make(map[string]*ir.TableRelation),
		aliases:        make(map[string]*ir.TableRelation),
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

func (s *scope) addRelation(relation *ir.TableRelation) {
	s.relations[relation.Name] = relation
}

func (s *scope) hasAlias(alias string) bool {
	_, ok := s.aliases[alias]
	return ok
}

func (s *scope) addAlias(alias string, relation *ir.TableRelation) {
	s.aliases[alias] = relation
}
