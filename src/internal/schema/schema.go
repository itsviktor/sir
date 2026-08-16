package schema

type Table interface {
	Name() string
	Columns() []Column
}

type Column interface {
	Name() string
	DefaultValue() *string
	Type() ColumnType
	IsNullable() bool
	IsPrimaryKey() bool
}

type ColumnType interface {
	DbName() string
	GoName() string
	Imports() []string
}

type DefaultType struct {
	dbName  string
	goName  string
	imports []string
}

func NewDefaultType(dbName, goName string, imports []string) ColumnType {
	return &DefaultType{dbName: dbName, goName: goName, imports: imports}
}

func (t DefaultType) DbName() string {
	return t.dbName
}

func (t DefaultType) GoName() string {
	return t.goName
}

func (t DefaultType) Imports() []string {
	return t.imports
}
