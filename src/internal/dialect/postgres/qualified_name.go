package postgres

import (
	"slices"
	"strings"

	"github.com/antlr4-go/antlr/v4"
	"github.com/itsviktor/sir/src/internal/parser"
)

// parseQualifiedParts parses relations and columns names from the colid and indirection parts.
// First part before dot always belongs to the colid context, while everything else falls
// into the indirection context.
//
// Supports star identifiers.
//
// Returns an array of sanitized identifiers in reversed order.
//
// Examples:
//
//	"public"."Table" -> [public, Table].
//	"table".Column -> [table, column]
//	table.* -> [table, *]
func parseQualifiedParts(colid *parser.ColidContext, indirection parser.IIndirectionContext) []string {
	parts := []string{sanitizeIdentifier(colid.GetText())}

	if indirection == nil {
		return parts
	}

	indCtx := indirection.(*parser.IndirectionContext)
	for _, el := range indCtx.AllIndirection_el() {
		elCtx := el.(*parser.Indirection_elContext)
		secondChild := elCtx.GetChild(1).(antlr.ParseTree)

		parts = append(parts, sanitizeIdentifier(secondChild.GetText()))
	}

	slices.Reverse(parts)

	return parts
}

// sanitizeIdentifier removes quotes from a string identifier, such as a database
// or alias name. PostgreSQL folds unquoted identifiers to lowercase, so
// TableName and tablename refer to the same identifier unless quoted.
//
// Examples:
//
//	"Name"         -> Name
//	"name""valid"  -> name"valid
//	TableName      -> tablename
func sanitizeIdentifier(s string) string {
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		return strings.ReplaceAll(s[1:len(s)-1], `""`, `"`)
	}

	return strings.ToLower(s)
}
