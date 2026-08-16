package schema

var registry = map[string]ColumnType{}

func RegisterType(name string, t ColumnType) {
	registry[name] = t
}

func GetRegisteredType(name string) (ColumnType, bool) {
	t, ok := registry[name]
	return t, ok
}
