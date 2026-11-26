package main

import (
	"context"

	"github.com/itsviktor/sir/internal/config"
	"github.com/itsviktor/sir/internal/inspector"
	"github.com/itsviktor/sir/internal/sqlite"
	"github.com/itsviktor/sir/internal/utils"
)

func main() {
	cfg := config.Load()

	var i inspector.Inspector
	switch cfg.Driver {
	case config.DriverSqlite:
		i = sqlite.NewInspector()
	default:
		utils.Fatal("unsupported driver", "name", cfg.Driver)
	}

	i.Inspect(context.Background(), "test.db")
}
