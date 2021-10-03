package repo

import (
	"context"
	"fmt"
	"os"

	"github.com/10hourlabs/tentn/ent"
	"github.com/10hourlabs/tentn/internal/database"
	"github.com/google/uuid"
	"github.com/gosimple/slug"
)

var (
	// database connection
	dBConn *ent.Client

	// database query context
	dBContext context.Context
)

func init() {
	dBContext = context.Background()
	client, err := database.NewPostgresClient(os.Getenv("TENTN_POSTGRES_DSN"))
	if err != nil {
		// TODO should we panic or switch to in memory sql database?
		panic("Failed to initialize database")
	}
	dBConn = client
}

func slugify(title string, id uuid.UUID) string {
	return slug.Make(fmt.Sprintf("%v %v", title, id))
}
