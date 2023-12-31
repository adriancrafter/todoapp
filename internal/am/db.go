package am

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type (
	DB interface {
		Core
		DB() *sqlx.DB
		Tx() *sqlx.Tx
		Path() string
		Schema() string
		Name() string
		Connect(ctx context.Context) error
	}
)
