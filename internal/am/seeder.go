package am

import (
	"database/sql"
)

type (
	Seeder interface {
		Core
		// Seed applies pending seeding
		Seed() error
		// SetAssetsPath sets the path vm where the seeding are read
		SetAssetsPath(path string)
		// AssetsPath returns the path vm where the seeding are read
		AssetsPath() string
	}
)

// SeedFx type alias
type SeedFx = func(tx *sql.Tx) error

type SeederExec interface {
	Config(seeds []SeedFx)
	GetIndex() (i int64)
	GetName() (name string)
	GetSeeds() (seedFxs []SeedFx)
	SetTx(tx *sql.Tx)
	GetTx() (tx *sql.Tx)
}
