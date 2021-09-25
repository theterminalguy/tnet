package database

type Database interface {
	Open() error
	GetConnectionString() string
	RunMigration() error
}
