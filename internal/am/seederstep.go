package am

import (
	"database/sql"
)

type (
	SeederStep struct {
		Index int64
		Name  string
		Seeds []SeedFx
		tx    *sql.Tx
	}
)

func (s *SeederStep) Config(seed []SeedFx) {
	s.Seeds = seed
}

func (s *SeederStep) GetIndex() (idx int64) {
	return s.Index
}

func (s *SeederStep) GetName() (name string) {
	return s.Name
}

func (s *SeederStep) GetSeeds() (seeds []SeedFx) {
	return s.Seeds
}

func (s *SeederStep) SetTx(tx *sql.Tx) {
	s.tx = tx
}

func (s *SeederStep) GetTx() (tx *sql.Tx) {
	return s.tx
}
