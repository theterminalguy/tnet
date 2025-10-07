package database

import (
	"context"
	"fmt"
	"os"

	"entgo.io/ent/dialect"
	"github.com/theterminalguy/tentn/ent"
	"github.com/theterminalguy/tentn/ent/migrate"
	_ "github.com/lib/pq"
)

type DBPostgres struct {
	client *ent.Client
}

var _ Databaser = (*DBPostgres)(nil)

func (*DBPostgres) GetDSN() string {
	if os.Getenv("ENV") == "production" {
		return fmt.Sprintf(
			"user=%s password=%s database=%s host=%s/%s",
			os.Getenv("CLOUDSQL_PG_USER"),
			os.Getenv("CLOUDSQL_PG_PASSWORD"),
			os.Getenv("CLOUDSQL_PG_DBNAME"),
			os.Getenv("CLOUDSQL_PG_SOCKET_DIR"),
			os.Getenv("CLOUDSQL_PG_INSTANCE"),
		)
	}
	return fmt.Sprintf(
		"postgres://%v:%v@%v:%v/%v?sslmode=%v",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_SSL_MODE"),
	)
}

func (d *DBPostgres) Open() (*ent.Client, error) {
	return ent.Open(dialect.Postgres, d.GetDSN())
}

func (d *DBPostgres) RunMigration() error {
	if err := d.client.Schema.Create(
		context.Background(), // TODO: avoid Background context
		migrate.WithDropIndex(true),
		migrate.WithDropColumn(true),
	); err != nil {
		return NewCreateSchemaError(err)
	}
	return nil
}

func NewPostgresClient() (*ent.Client, error) {
	db := &DBPostgres{}
	client, err := db.Open()
	if err != nil {
		return nil, NewConnectionError(err)
	}
	db.client = client
	// Run the automatic migration tool to create all schema resources
	if err = db.RunMigration(); err != nil {
		return nil, err
	}
	return db.client, nil
}
