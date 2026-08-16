package postgres

import (
	"fmt"
	"strings"

	"github.com/itsviktor/sir/src/internal/schema"
)

var prefixes = map[string]string{
	// Integer
	"smallint": "int16",
	"integer":  "int32",
	"bigint":   "int64",

	// Floating point
	"real":             "float32",
	"double precision": "float64",

	// Exact numeric
	"numeric": "decimal.Decimal",
	"decimal": "decimal.Decimal",

	// String
	"character varying": "string",
	"varchar":           "string",
	"character":         "string",
	"char":              "string",
	"text":              "string",

	// Boolean
	"boolean": "bool",

	// UUID
	"uuid": "uuid.UUID",

	// Network
	"inet":     "netip.Prefix",
	"cidr":     "netip.Prefix",
	"macaddr":  "net.HardwareAddr",
	"macaddr8": "net.HardwareAddr",

	// Date / time
	"date":                        "time.Time",
	"timestamp without time zone": "time.Time",
	"timestamp with time zone":    "time.Time",
	"time without time zone":      "time.Time",
	"time with time zone":         "time.Time",
	"interval":                    "pgtype.Interval",

	// Binary
	"bytea": "[]byte",

	// JSON
	"json":  "[]byte",
	"jsonb": "[]byte",
}

var imports = map[string]string{
	"numeric": "github.com/shopspring/decimal",
	"decimal": "github.com/shopspring/decimal",

	"uuid": "github.com/google/uuid",

	"inet":     "net/netip",
	"cidr":     "net/netip",
	"macaddr":  "net",
	"macaddr8": "net",

	"date":                        "time",
	"timestamp without time zone": "time",
	"timestamp with time zone":    "time",
	"time without time zone":      "time",
	"time with time zone":         "time",
	"interval":                    "github.com/jackc/pgx/v5/pgtype",
}

// parseType parses type definition from postgres type name.
func parseType(dbType string, isNullable bool) (schema.ColumnType, error) {
	for prefix, goType := range prefixes {
		if !strings.HasPrefix(dbType, prefix) {
			continue
		}

		isArrayByItself := strings.HasPrefix(goType, "[]")
		isArrayType := strings.HasSuffix(dbType, "[]")

		goTypeBuilder := strings.Builder{}
		if isArrayType {
			if isNullable {
				goTypeBuilder.WriteString("*")
			}
			goTypeBuilder.WriteString("[]")

			if !isArrayByItself {
				goTypeBuilder.WriteString("*")
			}
			goTypeBuilder.WriteString(goType)
		} else {
			if isNullable && !isArrayByItself {
				goTypeBuilder.WriteString("*")
			}
			goTypeBuilder.WriteString(goType)
		}

		var requiredImports []string
		if i, ok := imports[dbType]; ok {
			requiredImports = append(requiredImports, i)
		}

		return schema.NewDefaultType(dbType, goTypeBuilder.String(), requiredImports), nil
	}

	t, ok := schema.GetRegisteredType(dbType)
	if ok {
		return t, nil
	}

	return nil, fmt.Errorf("unsupported pg type: %s", dbType)
}
