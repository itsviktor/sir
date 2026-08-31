package transformer

import (
	"fmt"

	"github.com/antlr4-go/antlr/v4"
	"github.com/itsviktor/sir/src/internal/ir"
	"github.com/itsviktor/sir/src/internal/schema"
)

// Scope stores relations and aliases in the current query context.
// Every scope is tied to an ANTLR node, usually a SELECT clause context node.
type Scope struct {
	parent         *Scope                // parent stores a pointer to the parent scope.
	nodeToChildren map[antlr.Tree]*Scope // nodeToChildren maps SELECT clause context nodes to child scopes.

	node antlr.Tree // node is the SELECT clause context node that created this scope.

	relations map[string]*ir.TableRelation // relations maps relation names to their information.
	aliases   map[string]*ir.TableRelation // aliases maps relation aliases to their information.
}

// NewScope creates new query Scope.
func NewScope(node antlr.Tree) *Scope {
	return &Scope{
		parent:         nil,
		nodeToChildren: make(map[antlr.Tree]*Scope),
		node:           node,
		relations:      make(map[string]*ir.TableRelation),
		aliases:        make(map[string]*ir.TableRelation),
	}
}

// Parent returns pointer to the parent scope.
func (s *Scope) Parent() *Scope {
	return s.parent
}

// Node returns scope's ANTLR node.
func (s *Scope) Node() antlr.Tree {
	return s.node
}

// AddChild adds new child scope.
func (s *Scope) AddChild(scope *Scope) {
	s.nodeToChildren[scope.Node()] = scope
	scope.parent = s
}

// FindChildrenByNode returns child by its node and successful flag.
func (s *Scope) FindChildrenByNode(node antlr.Tree) (child *Scope, ok bool) {
	child, ok = s.nodeToChildren[node]
	return child, ok
}

// RelationNames returns every relation name registered in the scope.
func (s *Scope) RelationNames() []string {
	var names []string
	if s.parent != nil {
		names = s.parent.RelationNames()
	}

	for name := range s.relations {
		names = append(names, name)
	}

	return names
}

// AddRelation adds a new relation to the scope.
// It also modifies the relation by linking it to the given table.
func (s *Scope) AddRelation(relation *ir.TableRelation, table *schema.Table) error {
	if relation.Alias != "" {
		_, ok := s.aliases[relation.Alias]
		if ok {
			return fmt.Errorf("duplicate alias %q", relation.Alias)
		}

		s.aliases[relation.Alias] = relation
	}

	relation.LinkTable(table)

	s.relations[relation.Name] = relation

	return nil
}

// FindRelation tries to find relation by its name or alias.
// Propagates the search to the parent scope if it exists.
//
// Returns the relation and a success flag.
func (s *Scope) FindRelation(name string) (*ir.TableRelation, bool) {
	r, ok := s.relations[name]
	if ok {
		return r, true
	}

	ar, ok := s.aliases[name]
	if ok {
		return ar, true
	}

	if s.parent != nil {
		return s.parent.FindRelation(name)
	}

	return nil, false
}

// FindRelationCandidatesNames returns the names of relations that contain the
// specified column.
//
// If no matching relations are found in the current scope,
// the search is propagated to the parent scope.
func (s *Scope) FindRelationCandidatesNames(columnName string) []string {
	var candidates []string

	for name, relation := range s.relations {
		if relation.HasColumn(columnName) {
			candidates = append(candidates, name)
		}
	}

	if len(candidates) > 0 {
		return candidates
	}

	if s.parent != nil {
		return s.parent.FindRelationCandidatesNames(columnName)
	}

	return []string{}
}
