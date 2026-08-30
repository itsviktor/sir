package ir

import "fmt"

type Op struct {
	Type OpType
	Text string
}

func NewOp(t OpType, text string) Op {
	return Op{
		Type: t,
		Text: text,
	}
}

type OpType int

const (
	// comparison
	Equal OpType = iota
	NotEqual
	Gt
	Gte
	Lt
	Lte

	// pattern matching
	Like
	ILike
	NotLike
	NotILike

	// arithmetic (binary)
	Plus
	Minus
	Percent
	Slash
	Star
	Caret

	// bitwise
	BitAnd
	BitOr
	BitXor
	BitShiftLeft
	BitShiftRight
	BitNot

	// logical
	And
	Or
	Not

	// unary arithmetic
	UnaryMinus
	UnaryPlus

	// fallback
	CustomUnary
	CustomBinary
)

func (op Op) String() string {
	switch op.Type {
	case Equal:
		return "="
	case NotEqual:
		return "<>"
	case Gt:
		return ">"
	case Gte:
		return ">="
	case Lt:
		return "<"
	case Lte:
		return "<="
	case Like:
		return "LIKE"
	case ILike:
		return "ILIKE"
	case NotLike:
		return "NOT LIKE"
	case NotILike:
		return "NOT ILIKE"
	case Plus, UnaryPlus:
		return "+"
	case Minus, UnaryMinus:
		return "-"
	case Percent:
		return "%"
	case Slash:
		return "/"
	case Star:
		return "*"
	case Caret:
		return "^"
	case BitAnd:
		return "&"
	case BitOr:
		return "|"
	case BitXor:
		return "#"
	case BitShiftLeft:
		return "<<"
	case BitShiftRight:
		return ">>"
	case BitNot:
		return "~"
	case And:
		return "AND"
	case Or:
		return "OR"
	case Not:
		return "NOT"
	case CustomUnary:
		return "CUSTOM UNARY: " + op.Text
	case CustomBinary:
		return "CUSTOM BINARY: " + op.Text
	}
	return fmt.Sprintf("unknown op type: %d", op.Type)
}

func StringToOpType(s string) (OpType, bool) {
	switch s {
	case "=":
		return Equal, true
	case "<>", "!=":
		return NotEqual, true
	case ">":
		return Gt, true
	case ">=":
		return Gte, true
	case "<":
		return Lt, true
	case "<=":
		return Lte, true
	case "+":
		return Plus, true
	case "-":
		return Minus, true
	case "%":
		return Percent, true
	case "/":
		return Slash, true
	case "*":
		return Star, true
	case "^":
		return Caret, true
	case "&":
		return BitAnd, true
	case "|":
		return BitOr, true
	case "#":
		return BitXor, true
	case "<<":
		return BitShiftLeft, true
	case ">>":
		return BitShiftRight, true
	case "~":
		return BitNot, true
	default:
		return 0, false
	}
}
