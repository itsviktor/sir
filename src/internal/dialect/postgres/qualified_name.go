package postgres

import (
	"strings"

	"github.com/itsviktor/sir/src/internal/parser"
)

func parseQualifiedParts(colid *parser.ColidContext, indirection parser.IIndirectionContext) []string {
	parts := []string{unquoteIdentifier(colid.GetText())}

	if indirection == nil {
		return parts
	}

	indCtx := indirection.(*parser.IndirectionContext)
	for _, el := range indCtx.AllIndirection_el() {
		elCtx := el.(*parser.Indirection_elContext)
		if attrName := elCtx.Attr_name(); attrName != nil {
			parts = append(parts, unquoteIdentifier(attrName.GetText()))
		}
	}

	return parts
}

func unquoteIdentifier(s string) string {
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		return strings.ReplaceAll(s[1:len(s)-1], `""`, `"`)
	}

	return strings.ToLower(s)
}
