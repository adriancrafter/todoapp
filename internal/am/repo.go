package am

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type (
	Repo interface {
		Core
		// DB returns the underlying sqlx database connection
		DB() (db *sqlx.DB)
		// Tx returns a sqlx transaction from the context if exists, otherwise it
		// creates a new one
		Tx(ctx context.Context) (tx *sqlx.Tx, isLocal bool)
		// Query returns a query string from the query manager for the given model
		Query(model, queryName string) (query string, err error)
	}
)

type (
	SimpleRepo struct {
		*SimpleCore
		db DB
		qm *QueryManager
	}
)

type (
	CtxKey string
)

const (
	TxKey CtxKey = "tx"
)

func NewRepo(name string, database DB, qm *QueryManager, opts ...Option) *SimpleRepo {
	return &SimpleRepo{
		SimpleCore: NewCore(name, opts...),
		db:         database,
		qm:         qm,
	}
}

func (r *SimpleRepo) DB() (db *sqlx.DB) {
	return r.db.DB()
}

func (r *SimpleRepo) Tx(ctx context.Context) (tx *sqlx.Tx, isLocal bool) {
	tx, ok := r.TxFromCtx(ctx)
	if !ok {
		tx, _ = r.DB().Beginx()
		isLocal = true
	}
	return tx, ok
}

func (r *SimpleRepo) Query(model, name string) (query string, err error) {
	return r.qm.Get(model, name)
}

func (r *SimpleRepo) Schema() string {
	return r.db.Schema()
}

func (r *SimpleRepo) TxFromCtx(ctx context.Context) (tx *sqlx.Tx, ok bool) {
	tx, ok = ctx.Value(TxKey).(*sqlx.Tx)
	return tx, ok
}

func (r *SimpleRepo) TxToCtx(ctx context.Context, tx *sqlx.Tx) context.Context {
	return context.WithValue(ctx, TxKey, tx)
}
