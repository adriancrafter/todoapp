package pg

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"strings"
	"time"

	_ "github.com/lib/pq"

	"github.com/google/uuid"

	"github.com/adriancrafter/todoapp/internal/am"
	"github.com/adriancrafter/todoapp/internal/am/db"
	"github.com/adriancrafter/todoapp/internal/am/errors"
)

const (
	seedTable = "seeds"
	seedPath  = "assets/seeding/pg"
)

type (
	// Seeder struct.
	Seeder struct {
		am.Core
		assetsPath string
		dbName     string
		schema     string
		fs         embed.FS
		db         *sql.DB
		steps      []Seed
	}

	// Seed struct
	Seed struct {
		Order    int
		Executor am.SeederExec
	}

	seedRecord struct {
		ID        uuid.UUID      `dbPath:"id" json:"id"`
		Index     sql.NullInt64  `dbPath:"index" json:"index"`
		Name      sql.NullString `dbPath:"name" json:"name"`
		CreatedAt db.NullTime    `dbPath:"created_at" json:"createdAt"`
	}
)

type (
	SeedRecord struct {
		ID        string
		Index     int
		Name      string
		CreatedAt string
	}
)

func NewSeeder(fs embed.FS, opts ...am.Option) (mig *Seeder) {
	m := &Seeder{
		Core:       am.NewCore("seeder", opts...),
		assetsPath: seedPath,
		fs:         fs,
	}

	return m
}

func (s *Seeder) SetAssetsPath(path string) {
	s.assetsPath = path
}

func (s *Seeder) AssetsPath() string {
	return s.assetsPath
}

func (s *Seeder) Start(ctx context.Context) error {
	s.Log().Infof("%s started", s.Name())

	s.dbName = s.Cfg().GetString(cfgKey.PgDB)
	s.schema = s.Cfg().GetString(cfgKey.PgSchema)

	err := s.Connect()
	if err != nil {
		return errors.Wrapf(err, "%s start error", s.Name())
	}

	err = s.addSteps()
	if err != nil {
		return errors.Wrapf(err, "%s start error", s.Name())
	}

	return s.Seed()
}

func (s *Seeder) Connect() error {
	pgDB, err := sql.Open("postgres", s.connectionString())
	if err != nil {
		return errors.Wrapf(err, "%s connection error", s.Name())
	}

	err = pgDB.Ping()
	if err != nil {
		msg := fmt.Sprintf("%s ping connection error", s.Name())
		return errors.Wrap(err, msg)
	}

	s.Log().Infof("%s database connected", s.Name())

	s.db = pgDB

	return nil
}

// GetTx returns a new transaction from seeder connection
func (m *Seeder) GetTx() (tx *sql.Tx, err error) {
	tx, err = m.db.Begin()
	if err != nil {
		return tx, err
	}

	return tx, nil
}

// PreSetup creates database
// and seeds table if needed.
func (s *Seeder) PreSetup() (err error) {
	if !s.seedingTableExists() {
		err = s.createSeedsTable()
		if err != nil {
			return err
		}
	}

	return nil
}

// migTableExists returns true if migration table exists.
func (s *Seeder) seedingTableExists() bool {
	query := `SELECT EXISTS (SELECT 1 
               FROM information_schema.schemata s
                   JOIN information_schema.tables t
                       ON t.table_schema = s.schema_name
               WHERE s.schema_name = '%s'
                 AND t.table_name = '%s');`

	st := fmt.Sprintf(query, s.schema, seedTable)

	rows, err := s.db.Query(st)
	if err != nil {
		s.Log().Errorf("error checking database: %s", err)
		return false
	}
	defer rows.Close()

	for rows.Next() {
		var exists bool
		err = rows.Scan(&exists)
		if err != nil {
			s.Log().Errorf("cannot read query result: %s\n", err)
			return false
		}

		return exists
	}

	return false
}

func (s *Seeder) CloseAppConns() (string, error) {
	dbName := s.Cfg().GetString(am.Key.SQLiteFilePath)

	err := s.db.Close()
	if err != nil {
		return dbName, err
	}

	adminConn, err := sql.Open("sqlite3", s.Name())
	if err != nil {
		return dbName, err
	}
	defer adminConn.Close()

	// Terminate all connections to the database (SQLite does not support concurrent connections)
	st := fmt.Sprintf(`PRAGMA busy_timeout = 5000;`)
	_, err = adminConn.Exec(st)
	if err != nil {
		return dbName, err
	}

	return dbName, nil
}

func (s *Seeder) createSeedsTable() (err error) {
	tx, err := s.GetTx()
	if err != nil {
		return err
	}

	query := `CREATE TABLE %s (
    id UUID PRIMARY KEY,
    idx INTEGER,
    name VARCHAR(64),
    created_at TEXT
    );`

	st := fmt.Sprintf(query, seedTable)

	_, err = tx.Exec(st)
	if err != nil {
		err2 := tx.Rollback()
		if err2 != nil {
			return err
		}
		return err
	}

	return tx.Commit()
}

func (s *Seeder) AddSeed(o int, e am.SeederExec) {
	mig := Seed{Order: o, Executor: e}
	s.steps = append(s.steps, mig)
}

func (s *Seeder) Seed() (err error) {
	err = s.PreSetup()
	if err != nil {
		return errors.Wrap(err, "seeder error")
	}

	for i, _ := range s.steps {
		mg := s.steps[i]
		exec := mg.Executor
		idx := exec.GetIndex()
		name := exec.GetName()
		seedFxs := exec.GetSeeds()

		// GetErr a new Tx from migrator
		tx, err := s.GetTx()
		if err != nil {
			return errors.Wrap(err, "seeder error")
		}

		//Continue if already applied
		if !s.canApplySeed(idx, name, tx) {
			s.Log().Infof("Seed '%s' already applied", name)
			tx.Commit() // No need to handle eventual error here
			continue
		}

		// Pass Tx to the executor
		for _, sfx := range seedFxs {
			err = sfx(tx)
			if err != nil {
				break
			}
		}

		if err != nil {
			s.Log().Infof("%s seeder not executed", name)
			err2 := tx.Rollback()
			if err2 != nil {
				return errors.Wrap(err2, "seeder rollback error")
			}

			return errors.Wrapf(err, "cannot run seeder '%s'", name)
		}

		// Register seeder
		exec.SetTx(tx)
		err = s.recSeed(exec)

		err = tx.Commit()
		if err != nil {
			msg := fmt.Sprintf("Cannot update seeder table: %s\n", err.Error())
			s.Log().Errorf("seeder commit error: %s", msg)
			err = tx.Rollback()
			if err != nil {
				return errors.Wrap(err, "seeder rollback error")
			}
			return errors.NewError(msg)
		}

		s.Log().Infof("Seed executed: %s", name)
	}

	return nil
}

func (s *Seeder) recSeed(e am.SeederExec) error {
	query := `INSERT INTO %s.%s (id, idx, name, created_at) VALUES ($1, $2, $3, $4);`

	st := fmt.Sprintf(query, s.schema, seedTable)

	uid, err := uuid.NewUUID()
	if err != nil {
		return errors.Wrap(err, "cannot update seeder table")
	}

	_, err = e.GetTx().Exec(st,
		ToNullString(uid.String()),
		ToNullInt64(e.GetIndex()),
		ToNullString(e.GetName()),
		ToNullTime(time.Now()),
	)

	if err != nil {
		s.Log().Error(err)
		return errors.Wrap(err, "cannot update seeder table")
	}

	return nil
}

func (s *Seeder) cancelRollback(index int64, name string, tx *sql.Tx) bool {
	query := `SELECT (COUNT(*) > 0) AS record_exists FROM %s 
                                       WHERE idx = %d 
                                           :AND name = '%s'`

	st := fmt.Sprintf(query, seedTable, index, name)
	r, err := tx.Query(st)

	if err != nil {
		s.Log().Errorf("Cannot determine rollback status: %w", err)
		return true
	}

	for r.Next() {
		var applied sql.NullBool
		err = r.Scan(&applied)
		if err != nil {
			s.Log().Errorf("Cannot determine seeder status: %w", err)
			return true
		}

		return !applied.Bool
	}

	return true
}

func (s *Seeder) canApplySeed(index int64, name string, tx *sql.Tx) bool {
	query := `SELECT (COUNT(*) > 0) AS record_exists FROM %s 
                                       WHERE idx = %d 
                                           AND name = '%s'`

	st := fmt.Sprintf(query, seedTable, index, name)
	r, err := tx.Query(st)
	defer r.Close()

	if err != nil {
		s.Log().Errorf("Cannot determine seeder status: %w", err)
		return false
	}

	for r.Next() {
		var exists sql.NullBool
		err = r.Scan(&exists)
		if err != nil {
			s.Log().Errorf("Cannot determine seeder status: %s", err)
			return false
		}

		return !exists.Bool
	}

	return true
}

func (s *Seeder) addSteps() error {
	qq, err := s.readInsertSets()
	if err != nil {
		return err
	}

	for i, q := range qq {
		var seeds []am.SeedFx
		for _, i := range q.Inserts {
			seeds = append(seeds, s.genTxExecFunc(i))
		}

		step := &am.SeederStep{
			Index: q.Index,
			Name:  q.Name,
			Seeds: seeds,
		}

		s.AddSeed(i, step)
	}

	return nil
}

func (s *Seeder) genTxExecFunc(query string) func(tx *sql.Tx) error {
	return func(tx *sql.Tx) error {
		_, err := tx.Exec(query)
		return err
	}
}

type insertSet struct {
	Index   int64
	Name    string
	Inserts []string
}

func (s *Seeder) readInsertSets() ([]insertSet, error) {
	var iiss []insertSet

	files, err := s.fs.ReadDir(s.assetsPath)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".sql") {
			continue
		}

		filePath := fmt.Sprintf("%s/%s", s.assetsPath, file.Name())
		content, err := s.fs.ReadFile(filePath)
		if err != nil {
			return nil, err
		}

		var statements []string
		insertStmts := strings.Split(string(content), "--SEED")
		if len(insertStmts) < 1 {
			msg := fmt.Sprintf("invalid seeder file format: %s", file.Name())
			return nil, errors.NewError(msg)
		}

		for _, istmt := range insertStmts {
			insertSt := strings.TrimSpace(strings.TrimPrefix(istmt, "--SEED\n"))
			statements = append(statements, insertSt)
		}

		idx, name := stepName(filePath)

		is := insertSet{
			Index:   idx,
			Name:    name,
			Inserts: statements,
		}

		iiss = append(iiss, is)
	}

	return iiss, nil
}

func (s *Seeder) count() (last int) {
	return len(s.steps)
}

func (s *Seeder) last() (last int) {
	return s.count() - 1
}

func (s *Seeder) connectionString() (connString string) {
	cfg := s.Cfg()
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
