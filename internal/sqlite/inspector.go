package sqlite

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"time"

	"github.com/itsviktor/sir/internal/schema"
	"github.com/itsviktor/sir/internal/utils"
	_ "modernc.org/sqlite"
)

//go:embed sqlite_setup.sql
var setupSQL string

type Inspector struct{}

func NewInspector() *Inspector {
	return &Inspector{}
}

func connect(ctx context.Context, filename string) (*sql.DB, error) {
	dsn := "file:" + filename +
		"?_pragma=foreign_keys(1)" +
		"&_pragma=busy_timeout(5000)"

	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := db.PingContext(pingCtx); err != nil {
		db.Close()
		return nil, fmt.Errorf("ping sqlite: %w", err)
	}

	return db, nil
}

func (i *Inspector) Inspect(ctx context.Context, databaseName string) (map[string]*schema.Table, error) {
	db, err := connect(ctx, databaseName)
	if err != nil {
		utils.Fatal(err.Error())
	}

	err = utils.ExecSQLScript(ctx, db, setupSQL)
	if err != nil {
		utils.Fatal(err.Error())
	}

	return nil, nil
}
