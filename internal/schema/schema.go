package schema

type Table struct {
	Name       string
	Columns    map[string]*Column
	PrimaryKey *Column
}

type Column struct {
	Name string
	Type *Type
}

type Type struct {
	Name       string
	Nullable   bool
	PrimaryKey bool
}
