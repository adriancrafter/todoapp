package pg

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/adriancrafter/todoapp/internal/am"
	"github.com/adriancrafter/todoapp/internal/am/errors"
)

type DB struct {
	*am.SimpleCore
	db *sqlx.DB
}

var (
	cfgKey = am.Key
)

func NewDB(opts ...am.Option) *DB {
	name := "pg-db"

	return &DB{
		SimpleCore: am.NewCore(name, opts...),
	}
}

func (db *DB) Start(ctx context.Context) error {
	return db.Connect(ctx)
}

func (db *DB) Connect(ctx context.Context) error {
	pgDB, err := sql.Open("postgres", db.connString())
	if err != nil {
		msg := fmt.Sprintf("%s connection error", db.Name())
		return errors.Wrap(err, msg)
	}

	err = pgDB.Ping()
	if err != nil {
		msg := fmt.Sprintf("%s ping connection error", db.Name())
		return errors.Wrap(err, msg)
	}

	db.Log().Infof("%s database connected!", db.Name())
	return nil
}

func (db *DB) DBConn(ctx context.Context) (*sqlx.DB, error) {
	return db.db, nil
}

func (db *DB) DB() *sqlx.DB {
	return db.db
}

func (db *DB) Tx() *sqlx.Tx {
	return db.db.MustBegin()
}

func (db *DB) Schema() string {
	cfg := db.Cfg()
	dbPath := cfg.GetString(cfgKey.PgSchema)
	return dbPath
}

func (db *DB) Name() string {
	cfg := db.Cfg()
	dbPath := cfg.GetString(cfgKey.PgDB)
	return dbPath
}

func (db *DB) Path() string {
	return ""
}

func (db *DB) connString() (connString string) {
	cfg := db.Cfg()
	user := cfg.GetString(cfgKey.PgUser)
	pass := cfg.GetString(cfgKey.PgPass)
	name := cfg.GetString(cfgKey.PgDB)
	host := cfg.GetString(cfgKey.PgHost)
	port := cfg.GetInt(cfgKey.PgPort)
	schema := cfg.GetString(cfgKey.PgSchema)
	sslMode := cfg.GetBool(cfgKey.PgSSL)

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d search_path=%s", user, pass, name, host, port, schema)

	if sslMode {
		connStr = connStr + " sslmode=enable"
	} else {
		connStr = connStr + " sslmode=disable"
	}

	return connStr
}
