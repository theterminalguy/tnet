package database

import (
	"context"

	"entgo.io/ent/dialect"
	"github.com/10hourlabs/tentn/ent"
	_ "github.com/mattn/go-sqlite3"
)

type DBSQLite3 struct {
	ConnectionString string
	client           *ent.Client
}

func (d *DBSQLite3) Open() (*ent.Client, error) {
	return ent.Open(dialect.SQLite, d.GetConnectionString())
}

func (d *DBSQLite3) GetConnectionString() string {
	return d.ConnectionString
}

func (d *DBSQLite3) RunMigration() error {
	if err := d.client.Schema.Create(context.Background()); err != nil {
		return NewCreateSchemaError(err)
	}
	return nil
}

func NewSQLite3InMemoryClient() (*ent.Client, error) {
	db := &DBSQLite3{
		ConnectionString: "file:ent?mode=memory&cache=shared&_fk=1",
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
