package infrastrucure

import (
	"database/sql"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
)

// DatabaseImpl struct implements Database interface
type Database struct {
	Database *bun.DB
}

func (db *Database) DB() *bun.DB {
	return db.Database
}

func NewDatabase() (*Database, error) {
	sqldb, err := sql.Open(sqliteshim.ShimName, "file::memory:?cache=shared")
	db := bun.NewDB(sqldb, sqlitedialect.New())
	if err != nil {
		panic(err)
	}
	return &Database{
		Database: db,
	}, err
}
