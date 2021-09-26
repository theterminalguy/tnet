package database

import (
	"context"

	"entgo.io/ent/dialect"
	"github.com/10hourlabs/tentn/ent"
)

type DBPostgres struct {
	dsn    string
	client *ent.Client
}

var _ Databaser = (*DBPostgres)(nil)

func (d *DBPostgres) Open() (*ent.Client, error) {
	return ent.Open(dialect.Postgres, d.GetDSN())
}

func (d *DBPostgres) GetDSN() string {
	return d.dsn
}

func (d *DBPostgres) RunMigration() error {
	if err := d.client.Schema.Create(context.Background()); err != nil {
		return NewCreateSchemaError(err)
	}
	return nil
}

func NewPostgresClient() (*ent.Client, error) {
	db := &DBPostgres{
		// TODO: move dsn to environmental variable
		dsn: "host=localhost port=5341 user=theterminalguy dbname=tentn_dev",
	}
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
