package postgres

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/itsviktor/sir/src/internal/schema"
)

// Inspect returns tables and columns using the provided connection.
func Inspect(db *sql.DB) ([]schema.Table, error) {
	var tables []schema.Table

	err := parseEnumTypes(db)
	if err != nil {
		return tables, fmt.Errorf("getting types: %w", err)
	}

	rows, err := db.Query(`
		SELECT c.relname AS table_name
		FROM pg_catalog.pg_class c
		JOIN pg_catalog.pg_namespace n
			ON n.oid = c.relnamespace
		WHERE n.nspname = 'public'
		AND c.relkind = 'r';`,
	)
	if err != nil {
		return tables, fmt.Errorf("getting tables: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var table pgTable
		if err := rows.Scan(&table.NameC); err != nil {
			return tables, fmt.Errorf("scanning tables: %w", err)
		}

		columns, err := parseColumns(db, table.NameC)
		if err != nil {
			return tables, fmt.Errorf("getting columns for table %s: %w", table.NameC, err)
		}

		table.columns = columns

		tables = append(tables, table)
	}
	if err = rows.Err(); err != nil {
		return tables, err
	}

	return tables, nil
}

// parseEnumTypes parses and registers all enums types in the provided connection.
func parseEnumTypes(db *sql.DB) error {
	const query = `
		SELECT
			t.typname AS name,
			json_agg(e.enumlabel ORDER BY e.enumsortorder) AS values
		FROM pg_type t
		JOIN pg_enum e
			ON e.enumtypid = t.oid
		JOIN pg_namespace n
			ON n.oid = t.typnamespace
		WHERE n.nspname = 'public'
		GROUP BY t.typname
		ORDER BY t.typname;
	`

	rows, err := db.Query(query)
	if err != nil {
		return err
	}

	for rows.Next() {
		var t pgEnumType
		var values []byte

		if err := rows.Scan(&t.NameC, &values); err != nil {
			return fmt.Errorf("scanning type: %w", err)
		}

		if err := json.Unmarshal(values, &t.Values); err != nil {
			return fmt.Errorf("parsing enum values: %w", err)
		}

		schema.RegisterType(t.NameC, t)
	}
	if err := rows.Err(); err != nil {
		return err
	}

	return nil
}

// parseColumns parses all columns of the provided table.
func parseColumns(db *sql.DB, tableName string) ([]pgColumn, error) {
	var columns []pgColumn

	const query = `
		SELECT
			a.attname AS name,
			pg_get_expr(d.adbin, d.adrelid) AS default_value,
			format_type(a.atttypid, a.atttypmod) AS type,
			EXISTS (
				SELECT 1
				FROM pg_index i
				WHERE i.indrelid = a.attrelid
				AND i.indisprimary
				AND a.attnum = ANY(i.indkey)
			) AS is_primary_key,
			NOT a.attnotnull AS is_nullable
		FROM pg_attribute a
		LEFT JOIN pg_attrdef d
			ON d.adrelid = a.attrelid
			AND d.adnum = a.attnum
		WHERE a.attrelid = $1::regclass
		AND a.attnum > 0
		AND NOT a.attisdropped
		ORDER BY a.attnum;
	`

	rows, err := db.Query(query, fmt.Sprintf("public.%s", tableName))
	if err != nil {
		return columns, err
	}

	for rows.Next() {
		var column pgColumn
		if err := rows.Scan(&column.NameC, &column.DefaultValueC, &column.DbTypeC, &column.IsPrimaryKeyC, &column.IsNullableC); err != nil {
			return columns, fmt.Errorf("scanning column: %w", err)
		}

		t, err := parseType(column.DbTypeC, column.IsNullableC)
		if err != nil {
			return columns, fmt.Errorf("parsing type for column %s: %w", column.NameC, err)
		}
		column.t = t

		columns = append(columns, column)
	}
	if err := rows.Err(); err != nil {
		return columns, err
	}

	return columns, nil
}
