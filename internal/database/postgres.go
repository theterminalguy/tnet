package database

import (
	"context"

	"entgo.io/ent/dialect"
	"github.com/10hourlabs/tentn/ent"
)

type DBPostgres struct {
	ConnectionString string
	client           *ent.Client
}

func (d *DBPostgres) Open() (*ent.Client, error) {
	return ent.Open(dialect.Postgres, d.GetConnectionString())
}

func (d *DBPostgres) GetConnectionString() string {
	return d.ConnectionString
}

func (d *DBPostgres) RunMigration() error {
	if err := d.client.Schema.Create(context.Background()); err != nil {
		return NewCreateSchemaError(err)
	}
	return nil
}

func NewPostgresClient() (*ent.Client, error) {
	db := &DBPostgres{
		ConnectionString: "",
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
