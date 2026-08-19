package schema

var registry = map[string]ColumnType{}

// RegisterType registers new global type definition for a database type.
// It is used to register PostgreSQL enum types and custom overrides for default type resolution.
func RegisterType(databaseTypeName string, t ColumnType) {
	registry[databaseTypeName] = t
}

// GetRegisteredType returns type definition for the provided database type name and a success flag.
func GetRegisteredType(databaseTypeName string) (ColumnType, bool) {
	t, ok := registry[databaseTypeName]
	return t, ok
}
