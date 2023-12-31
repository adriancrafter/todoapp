package am

import "database/sql"

type (
	MigrationStep struct {
		Index int64
		Name  string
		Up    []MigFx
		Down  []MigFx
		tx    *sql.Tx
	}
)

func (s *MigrationStep) Config(up []MigFx, down []MigFx) {
	s.Up = up
	s.Down = down
}

func (s *MigrationStep) GetIndex() (idx int64) {
	return s.Index
}

func (s *MigrationStep) GetName() (name string) {
	return s.Name
}

func (s *MigrationStep) GetUp() (up []MigFx) {
	return s.Up
}

func (s *MigrationStep) GetDown() (down []MigFx) {
	return s.Down
}

func (s *MigrationStep) SetTx(tx *sql.Tx) {
	s.tx = tx
}

func (s *MigrationStep) GetTx() (tx *sql.Tx) {
	return s.tx
}
