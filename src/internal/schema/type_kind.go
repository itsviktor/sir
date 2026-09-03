package schema

import "fmt"

type TypeKind int

const (
	Integer TypeKind = iota
	Float
	Numeric
	String
	Boolean
	Enum
	Unknown
	Variable
)

func TypeKindToString(tk TypeKind) string {
	switch tk {
	case Unknown:
		return "Unknown"
	case Variable:
		return "Variable"
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

func IsNumber(tk TypeKind) bool {
	return tk <= Numeric
}
