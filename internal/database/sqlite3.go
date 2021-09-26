package database

import (
	"context"

	"entgo.io/ent/dialect"
	"github.com/10hourlabs/tentn/ent"
	_ "github.com/mattn/go-sqlite3"
)

type DBSQLite3 struct {
	dsn    string
	client *ent.Client
}

var _ Databaser = (*DBSQLite3)(nil)

func (d *DBSQLite3) Open() (*ent.Client, error) {
	return ent.Open(dialect.SQLite, d.GetDSN())
}

func (d *DBSQLite3) GetDSN() string {
	return d.dsn
}

func (d *DBSQLite3) RunMigration() error {
	if err := d.client.Schema.Create(context.Background()); err != nil {
		return NewCreateSchemaError(err)
	}
	return nil
}

func NewSQLite3InMemoryClient() (*ent.Client, error) {
	db := &DBSQLite3{
		// TODO: move dsn to environmental variable
		dsn: "file:ent?mode=memory&cache=shared&_fk=1",
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
