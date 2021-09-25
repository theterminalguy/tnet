package database

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect"
	"github.com/10hourlabs/tentn/ent"
)

type DBPostgres struct {
	ConnectionString string
	client           *ent.Client
}

func (d *DBPostgres) Open() (*ent.Client, error) {
	return ent.Open(dialect.SQLite, d.GetConnectionString())
}

func (d *DBPostgres) GetConnectionString() string {
	return d.ConnectionString
}

func (d *DBPostgres) RunMigration() error {
	if err := d.client.Schema.Create(context.Background()); err != nil {
		// TODO refactor error string
		// see https://travix.io/errors-derived-from-constants-in-go-fda6748b4072
		return errors.New(fmt.Sprintf("failed creating schema resources: %v", err))
	}
	return nil
}

func NewPostgresClient() (*ent.Client, error) {
	db := &DBPostgres{
		ConnectionString: "",
	}
	client, err := db.Open()
	if err != nil {
		// TODO refactor error string
		// see https://travix.io/errors-derived-from-constants-in-go-fda6748b4072
		return nil, errors.New(fmt.Sprintf("failed opening connection to sqlite: %v", err))
	}
	db.client = client
	// Run the automatic migration tool to create all schema resources
	if err = db.RunMigration(); err != nil {
		return nil, err
	}
	return db.client, nil
}
