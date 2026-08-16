package dbtype

import (
	"fmt"
)

type DbType string

const (
	TypePostgres DbType = "postgres"
	TypeSql      DbType = "sql"
	TypeSqlite   DbType = "sqlite"
)

var supportedDrivers = map[string]DbType{
	"pgx": TypePostgres,
}

// ParseAndValidate gets the driver name, checks whether it is in the list of supported drivers, and returns the appropriate database type.
func ParseAndValidate(driver string) (DbType, error) {
	dbType, ok := supportedDrivers[driver]
	if !ok {
		return TypePostgres, fmt.Errorf(
			"unsupported driver: %s. Currently only the following drivers are supported: %v",
			driver,
			supportedDrivers,
		)
	}

	return dbType, nil
}
