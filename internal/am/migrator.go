package am

import (
	"database/sql"
)

type (
	Migrator interface {
		Core
		// Migrate applies pending seeding
		Migrate() error
		// Rollback reverts from one to N seeding already applied
		Rollback(steps ...int) error
		// RollbackAll reverts all seeding allready applied
		RollbackAll() error
		// SoftReset apply all seeding again after rolling back all seeding.
		SoftReset() error
		// Reset apply all seeding again after dropping the database and recreating it
		Reset() error
		// SetAssetsPath sets the path vm where the seeding are read
		SetAssetsPath(path string)
		// AssetsPath returns the path vm where the seeding are read
		AssetsPath() string
	}
)

type MigFx = func(tx *sql.Tx) error

type MigratorExec interface {
	Config(up []MigFx, down []MigFx)
	GetIndex() (i int64)
	GetName() (name string)
	GetUp() (up []MigFx)
	GetDown() (down []MigFx)
	SetTx(tx *sql.Tx)
	GetTx() (tx *sql.Tx)
}
