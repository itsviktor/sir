package dsql

// Domain is a group of queries that will be transformed into a repository.
// Each query becomes a repository method.
type Domain struct {
	Name    string  // Name is the PascalCase name derived from the .dsql filename.
	Queries []Query // Queries holds all queries found in the file.
}
