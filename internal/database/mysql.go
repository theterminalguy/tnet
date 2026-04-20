package database

import (
	"context"
	"fmt"
	"os"

	"entgo.io/ent/dialect"
	_ "github.com/go-sql-driver/mysql"
	"github.com/theterminalguy/tnet/ent"
	"github.com/theterminalguy/tnet/ent/migrate"
)

type DBMySQL struct {
	client *ent.Client
}

var _ Databaser = (*DBMySQL)(nil)

func (*DBMySQL) GetDSN() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=True",
		os.Getenv("MYSQL_USER"),
		os.Getenv("MYSQL_PASSWORD"),
		os.Getenv("MYSQL_HOST"),
		os.Getenv("MYSQL_PORT"),
		os.Getenv("MYSQL_DATABASE"),
	)
}

func (d *DBMySQL) Open() (*ent.Client, error) {
	return ent.Open(dialect.MySQL, d.GetDSN())
}

func (d *DBMySQL) RunMigration() error {
	if err := d.client.Schema.Create(
		context.Background(),
		migrate.WithDropIndex(true),
		migrate.WithDropColumn(true),
	); err != nil {
		return NewCreateSchemaError(err)
	}
	return nil
}

func NewMySQLClient() (*ent.Client, error) {
	db := &DBMySQL{}
	client, err := db.Open()
	if err != nil {
		return nil, NewConnectionError(err)
	}
	db.client = client
	if err = db.RunMigration(); err != nil {
		return nil, err
	}
	return db.client, nil
}
