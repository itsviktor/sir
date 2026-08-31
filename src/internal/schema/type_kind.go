package schema

import "fmt"

type TypeKind int

const (
	Unknown TypeKind = iota
	Integer
	Float
	Numeric
	String
	Boolean
	Enum
)

func TypeKindToString(tk TypeKind) string {
	switch tk {
	case Unknown:
		return "Unknown"
	case Integer:
		return "Integer"
	case Float:
		return "Float"
	case Numeric:
		return "Numeric"
	case String:
		return "String"
	case Boolean:
		return "Boolean"
	case Enum:
		return "Enum"
	default:
		return fmt.Sprintf("unknown type kind: %d", tk)
	}
}
