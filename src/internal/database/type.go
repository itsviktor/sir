package database

import "fmt"

// Dialect is a type of the database.
type Dialect string

const (
	Postgres Dialect = "postgres"
	Mysql    Dialect = "mysql"
	Sqlite   Dialect = "sqlite"
)

var supportedDrivers = map[string]Dialect{
	"pgx": Postgres,
}

// DialectFromDriver gets the driver name, checks whether it is in the list of supported drivers, and returns the appropriate database dialect.
func DialectFromDriver(driver string) (Dialect, error) {
	dbType, ok := supportedDrivers[driver]
	if !ok {
		return Postgres, fmt.Errorf(
			"unsupported driver: %s. Currently only the following drivers are supported: %v",
			driver,
			supportedDrivers,
		)
	}

	return dbType, nil
}
