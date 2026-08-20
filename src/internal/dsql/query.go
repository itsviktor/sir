package dsql

// Kind describes the return value of a query.
type Kind string

const (
	// KindOne queries return only one instance of the SELECTed fields.
	KindOne Kind = "one"

	// KindMany queries return a slice of the SELECTed fields.
	KindMany Kind = "many"

	// KindCount is a special query type.
	// It is applicable to queries that use SELECT COUNT().
	// The generated repository method returns int64.
	KindCount Kind = "count"

	// KindExec queries normally do not return anything.
	// However, the generated repository method may return a value when the query contains a RETURNING clause.
	KindExec Kind = "exec"
)

// Query is a DTO which contains meta information about user's query from a file.
type Query struct {
	SQL       string // SQL is the raw query text with the header line removed.
	Name      string // Name is the PascalCase query name parsed from the header comment.
	Kind      Kind   // Kind is the return kind parsed from the header comment (one, many, count, exec).
	File      string // File is the path to the query's file.
	StartLine int    // StartLine is the line number where the query starts.
}
