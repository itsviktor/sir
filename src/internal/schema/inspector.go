package schema

import "database/sql"

// Inspector inspects a database and returns its table definitions.
type Inspector interface {
	// Inspect retrieves the database schema using the provided database connection.
	Inspect(db *sql.DB) (map[string]Table, error)
}
