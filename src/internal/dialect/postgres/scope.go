package postgres

import (
	"github.com/antlr4-go/antlr/v4"
	"github.com/itsviktor/sir/src/internal/ir"
)

type scope struct {
	parent   *scope
	children []*scope

	node antlr.Tree

	relations map[string]*ir.TableRelation
	aliases   map[string]*ir.TableRelation
}

func newScope(node antlr.Tree) *scope {
	return &scope{
		parent:    nil,
		children:  make([]*scope, 0),
		node:      node,
		relations: make(map[string]*ir.TableRelation),
		aliases:   make(map[string]*ir.TableRelation),
	}
}

func (s *scope) add(other *scope) {
	s.children = append(s.children, other)
	other.parent = s
}

func (s *scope) hasAlias(alias string) bool {
	_, ok := s.aliases[alias]
	return ok
}
