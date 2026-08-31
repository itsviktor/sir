package postgres

import (
	"fmt"
	"strings"

	"github.com/itsviktor/sir/src/internal/schema"
)

type typeInfo struct {
	goType          string
	requiredImports []string
	kind            schema.TypeKind
}

var typesInfo = map[string]typeInfo{
	// Integer
	"smallint": {
		goType: "int16",
		kind:   schema.Integer,
	},
	"integer": {
		goType: "int32",
		kind:   schema.Integer,
	},
	"bigint": {
		goType: "int64",
		kind:   schema.Integer,
	},

	// Floating point
	"real": {
		goType: "float32",
		kind:   schema.Float,
	},
	"double precision": {
		goType: "float64",
		kind:   schema.Float,
	},

	// Exact numeric
	"numeric": {
		goType:          "decimal.Decimal",
		requiredImports: []string{"github.com/shopspring/decimal"},
		kind:            schema.Numeric,
	},
	"decimal": {
		goType:          "decimal.Decimal",
		requiredImports: []string{"github.com/shopspring/decimal"},
		kind:            schema.Numeric,
	},

	// String
	"character varying": {
		goType: "string",
		kind:   schema.String,
	},
	"varchar": {
		goType: "string",
		kind:   schema.String,
	},
	"character": {
		goType: "string",
		kind:   schema.String,
	},
	"char": {
		goType: "string",
		kind:   schema.String,
	},
	"text": {
		goType: "string",
		kind:   schema.String,
	},

	// Boolean
	"boolean": {
		goType: "bool",
		kind:   schema.Boolean,
	},

	// UUID
	"uuid": {
		goType:          "uuid.UUID",
		requiredImports: []string{"github.com/google/uuid"},
		kind:            schema.String,
	},

	// Network
	"inet": {
		goType:          "netip.Prefix",
		requiredImports: []string{"net/netip"},
		kind:            schema.String,
	},
	"cidr": {
		goType:          "netip.Prefix",
		requiredImports: []string{"net/netip"},
		kind:            schema.String,
	},
	"macaddr": {
		goType:          "net.HardwareAddr",
		requiredImports: []string{"net"},
		kind:            schema.String,
	},
	"macaddr8": {
		goType:          "net.HardwareAddr",
		requiredImports: []string{"net"},
		kind:            schema.String,
	},

	// Date / time
	"date": {
		goType:          "time.Time",
		requiredImports: []string{"time"},
		kind:            schema.String,
	},
	"timestamp without time zone": {
		goType:          "time.Time",
		requiredImports: []string{"time"},
		kind:            schema.String,
	},
	"timestamp with time zone": {
		goType:          "time.Time",
		requiredImports: []string{"time"},
		kind:            schema.String,
	},
	"time without time zone": {
		goType:          "time.Time",
		requiredImports: []string{"time"},
		kind:            schema.String,
	},
	"time with time zone": {
		goType:          "time.Time",
		requiredImports: []string{"time"},
		kind:            schema.String,
	},
	"interval": {
		goType:          "pgtype.Interval",
		requiredImports: []string{"github.com/jackc/pgx/v5/pgtype"},
		kind:            schema.String,
	},

	// Binary
	"bytea": {
		goType: "[]byte",
		kind:   schema.String,
	},

	// JSON
	"json": {
		goType: "[]byte",
		kind:   schema.String,
	},
	"jsonb": {
		goType: "[]byte",
		kind:   schema.String,
	},
}

// parseType parses type definition from postgres type name.
func parseType(dbType string, isNullable bool) (schema.ColumnType, error) {
	for prefix, info := range typesInfo {
		if !strings.HasPrefix(dbType, prefix) {
			continue
		}

		isArrayByItself := strings.HasPrefix(info.goType, "[]")
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
			goTypeBuilder.WriteString(info.goType)
		} else {
			if isNullable && !isArrayByItself {
				goTypeBuilder.WriteString("*")
			}
			goTypeBuilder.WriteString(info.goType)
		}

		return schema.NewDefaultType(dbType, goTypeBuilder.String(), info.kind, info.requiredImports), nil
	}

	t, ok := schema.GetRegisteredType(dbType)
	if ok {
		return t, nil
	}

	return nil, fmt.Errorf("unsupported pg type: %s", dbType)
}
