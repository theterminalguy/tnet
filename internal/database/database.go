package database

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect"
	"github.com/10hourlabs/tentn/ent"
	_ "github.com/mattn/go-sqlite3"
)

type DBSQLite3 struct {
	ConnectionString string
	client           *ent.Client
}

func (d *DBSQLite3) GetConnectionString() string {
	return d.ConnectionString
}

func (d *DBSQLite3) RunMigration() error {
	if err := d.client.Schema.Create(context.Background()); err != nil {
		// TODO refactor error string
		// see https://travix.io/errors-derived-from-constants-in-go-fda6748b4072
		return errors.New(fmt.Sprintf("failed creating schema resources: %v", err))
	}
	return nil
}

func NewSQLite3InMemoryClient() (*ent.Client, error) {
	adapter := &DBSQLite3{
		ConnectionString: "file:ent?mode=memory&cache=shared&_fk=1",
	}
	client, err := ent.Open(dialect.SQLite, adapter.GetConnectionString())
	if err != nil {
		// TODO refactor error string
		// see https://travix.io/errors-derived-from-constants-in-go-fda6748b4072
		return nil, errors.New(fmt.Sprintf("failed opening connection to sqlite: %v", err))
	}
	adapter.client = client
	// Run the automatic migration tool to create all schema resources
	if err = adapter.RunMigration(); err != nil {
		return nil, err
	}
	return adapter.client, nil
}
