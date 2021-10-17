package database

import (
	"context"
	"fmt"
	"os"

	"entgo.io/ent/dialect"
	"github.com/10hourlabs/tentn/ent"
	_ "github.com/lib/pq"
)

type DBPostgres struct {
	dsn    string
	client *ent.Client
}

var _ Databaser = (*DBPostgres)(nil)

func (d *DBPostgres) Open() (*ent.Client, error) {
	return ent.Open(dialect.Postgres, d.GetDSN())
}

func (*DBPostgres) GetDSN() string {
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

func (d *DBPostgres) RunMigration() error {
	if err := d.client.Schema.Create(context.Background()); err != nil {
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
