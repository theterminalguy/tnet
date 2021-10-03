package repository

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/10hourlabs/tentn/ent"
	"github.com/10hourlabs/tentn/internal/database"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gosimple/slug"
	_ "github.com/joho/godotenv/autoload"
)

var (
	// database connection
	dBConn *ent.Client

	// database query context
	dBContext context.Context

	// errors
	RecordNotFoundError error = errors.New("entity not found")
)

func init() {
	dBContext = context.Background()
	client, err := database.NewPostgresClient(os.Getenv("TENTN_POSTGRES_DSN"))
	if err != nil {
		panic(fmt.Sprintf("Database Error %v", err))
	}
	dBConn = client
}

func slugify(title string, id uuid.UUID) string {
	return slug.Make(fmt.Sprintf("%v %v", title, id))
}

func validateParams(s interface{}, fields ...string) error {
	validate := validator.New()
	if len(fields) > 0 {
		if err := validate.StructPartial(s, fields...); err != nil {
			return err
		}
	} else {
		if err := validate.Struct(s); err != nil {
			return err
		}
	}
	return nil
}
