package inspector

import (
	"context"

	"github.com/itsviktor/sir/internal/schema"
)

type Inspector interface {
	Inspect(ctx context.Context, databaseName string) (map[string]*schema.Table, error)
}
