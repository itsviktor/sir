package postgres

type pgEnumType struct {
	Name   string
	Values []string
}

func (t pgEnumType) DbName() string {
	return t.Name
}

func (t pgEnumType) GoName() string {
	return t.Name
}

func (t pgEnumType) Imports() []string {
	return []string{}
}
