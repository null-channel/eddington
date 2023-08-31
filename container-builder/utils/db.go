package utils

import (
	"context"
	"database/sql"

	"github.com/null-channel/eddington/container-builder/models"
	"github.com/pkg/errors"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
)

func NewDB(dbUrl string) (*bun.DB, error) {
	if dbUrl == "" {
		dbUrl = "file::memory:?cache=shared"
	}
	sqldb, err := sql.Open(sqliteshim.ShimName, dbUrl)
	if err != nil {
		return &bun.DB{}, err
	}
	db := bun.NewDB(sqldb, sqlitedialect.New())
	if err != nil {
		return &bun.DB{}, errors.Wrap(err, "unable to create db")
	}
	_, err = db.NewCreateTable().Model((*models.Build)(nil)).Exec(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "unable to migrate table")
	}
	return db, nil

}
