package pg

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"

	"github.com/google/uuid"

	"github.com/adriancrafter/todoapp/internal/am"
	"github.com/adriancrafter/todoapp/internal/am/errors"
)

const (
	migTable = "migrations"
	migPath  = "assets/migrations/pg"
)

type (
	// MigFx type alias
	MigFx = func(tx *sql.Tx) error

	// Migrator struct.
	Migrator struct {
		am.Core
		assetsPath string
		dbName     string
		schema     string
		fs         embed.FS
		db         *sql.DB
		steps      []Migration
	}

	// Migration struct.
	Migration struct {
		Order    int
		Executor am.MigratorExec
	}

	migRecord struct {
		ID        uuid.UUID      `dbPath:"id" json:"id"`
		Index     sql.NullInt64  `dbPath:"index" json:"index"`
		Name      sql.NullString `dbPath:"name" json:"name"`
		CreatedAt am.NullTime    `dbPath:"created_at" json:"createdAt"`
	}
)

type (
	MigrationRecord struct {
		ID        string
		Index     int
		Name      string
		CreatedAt string
	}
)

// NewMigrator returns a new Migrator instance.
func NewMigrator(fs embed.FS, opts ...am.Option) (mig *Migrator) {
	m := &Migrator{
		Core:       am.NewCore("migrator", opts...),
		assetsPath: migPath,
		fs:         fs,
	}

	return m
}

// SetAssetsPath sets the assets path.
func (m *Migrator) SetAssetsPath(path string) {
	m.assetsPath = path
}

// AssetsPath returns the assets path.
func (m *Migrator) AssetsPath() string {
	return m.assetsPath
}

// Start migrator.
func (m *Migrator) Start(ctx context.Context) error {
	m.Log().Infof("%s started", m.Name())

	m.dbName = m.Cfg().GetString(cfgKey.PgDB)
	m.schema = m.Cfg().GetString(cfgKey.PgSchema)

	err := m.Connect()
	if err != nil {
		return errors.Wrapf(err, "%s start error", m.Name())
	}

	err = m.addSteps()
	if err != nil {
		return errors.Wrapf(err, "%s start error", m.Name())
	}

	return m.Migrate()
}

// Close migrator to database.
func (m *Migrator) Connect() error {
	pgDB, err := sql.Open("postgres", m.connectionString())

	if err != nil {
		return errors.Wrapf(err, "%s connection error", m.Name())
	}

	err = pgDB.Ping()
	if err != nil {
		msg := fmt.Sprintf("%s ping connection error", m.Name())
		return errors.Wrap(err, msg)
	}

	m.Log().Infof("%s database connected", m.Name())

	m.db = pgDB

	return nil
}

// GetTx returns a new transaction from migrator db.
func (m *Migrator) GetTx() (tx *sql.Tx, err error) {
	tx, err = m.db.Begin()
	if err != nil {
		return tx, err
	}

	return tx, nil
}

// PreSetup migrator.
func (m *Migrator) PreSetup() (err error) {
	if !m.dbExists() {
		err := m.createDB()
		if err != nil {
			return err
		}
	}

	if !m.migTableExists() {
		err = m.createMigrationsTable()
		if err != nil {
			return err
		}
	}

	return nil
}

// dbExists returns true if database exists.
func (m *Migrator) dbExists() bool {
	query := `SELECT datname FROM pg_catalog.pg_database WHERE lower(datname) = lower('%s')`
	st := fmt.Sprintf(query, m.dbName)

	rows, err := m.db.Query(st)
	if err != nil {
		m.Log().Infof("error checking database: %w", err)
		return false
	}
	defer rows.Close()
	for rows.Next() {
		var dbName string
		err = rows.Scan(&dbName)
		if err != nil {
			m.Log().Errorf("Cannot read query result: %w", err)
			return false
		}
		return true
	}

	return false
}

// migTableExists returns true if migration table exists.
func (m *Migrator) migTableExists() bool {
	query := `SELECT EXISTS (SELECT 1 
               FROM information_schema.schemata s
                   JOIN information_schema.tables t
                       ON t.table_schema = s.schema_name
               WHERE s.schema_name = '%s'
                 AND t.table_name = '%s');`

	st := fmt.Sprintf(query, m.schema, migTable)

	rows, err := m.db.Query(st)
	if err != nil {
		m.Log().Errorf("error checking database: %s", err)
		return false
	}
	defer rows.Close()

	for rows.Next() {
		var exists bool
		err = rows.Scan(&exists)
		if err != nil {
			m.Log().Errorf("cannot read query result: %s\n", err)
			return false
		}

		return exists
	}

	return false
}

// createDB creates a new database.
func (m *Migrator) createDB() (err error) {
	err = m.closeAllConnections()
	if err != nil {
		return err
	}

	query := `CREATE DATABASE %s;`
	st := fmt.Sprintf(query, m.dbName)

	_, err = m.db.Exec(st)
	if err != nil {
		return err
	}

	return nil
}

// DropDB drops the database.
func (m *Migrator) DropDB() (dbPath string, err error) {
	err = m.closeAllConnections()
	if err != nil {
		return dbPath, errors.Wrap(err, "drop db error")
	}

	err = m.db.Close()
	if err != nil {
		m.Log().Errorf("drop dbPath error: %w", err)
	}

	err = os.Remove(dbPath)
	if err != nil {
		return dbPath, err
	}

	return dbPath, nil
}

func (m *Migrator) closeAllConnections() error {
	query := `SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname = '%s' 
                                                         AND pid <> pg_backend_pid()
                                                         AND EXISTS (
                                                         SELECT 1
                                                         FROM pg_namespace n
                                                             JOIN pg_database d ON d.datnamespace = n.oid
                                                         WHERE n.nspname = '%s'
                                                           AND d.datname = '%s');`

	st := fmt.Sprintf(query, m.dbName, m.schema, m.dbName)

	_, err := m.db.Exec(st)
	if err != nil {
		return err
	}

	return nil
}

// createMigrationsTable creates the migrations table.
func (m *Migrator) createMigrationsTable() (err error) {
	tx, err := m.GetTx()
	if err != nil {
		return err
	}

	query := `CREATE TABLE %s.%s (id UUID PRIMARY KEY,
    idx INTEGER,
    name VARCHAR(64),
    created_at TEXT
    );`

	st := fmt.Sprintf(query, m.schema, migTable)

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

// addMigration adds a new migration.
func (m *Migrator) AddMigration(o int, e am.MigratorExec) {
	mig := Migration{Order: o, Executor: e}
	m.steps = append(m.steps, mig)
}

// Migrate database.
func (m *Migrator) Migrate() (err error) {
	err = m.PreSetup()
	if err != nil {
		return errors.Wrap(err, "migrate error")
	}

	for i := range m.steps {
		mg := m.steps[i]
		exec := mg.Executor
		idx := exec.GetIndex()
		name := exec.GetName()
		upFx := exec.GetUp()

		// GetErr a new Tx from migrator
		tx, err := m.GetTx()
		if err != nil {
			return errors.Wrap(err, "migrate error")
		}

		//Continue if already applied
		if !m.canApplyMigration(idx, name, tx) {
			m.Log().Infof("migration '%s' already applied", name)
			tx.Commit() // No need to handle eventual error here
			continue
		}

		for _, mfx := range upFx {
			err = mfx(tx)
			if err != nil {
				break
			}
		}

		if err != nil {
			m.Log().Infof("%s migration not executed", name)
			err2 := tx.Rollback()
			if err2 != nil {
				return errors.Wrap(err2, "migrate rollback error")
			}

			return errors.Wrapf(err, "cannot run migration '%s'", name)
		}

		exec.SetTx(tx)
		err = m.recMigration(exec)

		err = tx.Commit()
		if err != nil {
			msg := fmt.Sprintf("Cannot update migration table: %s\n", err.Error())
			m.Log().Errorf("migrate commit error: %s", msg)
			err = tx.Rollback()
			if err != nil {
				return errors.Wrap(err, "migrate rollback error")
			}
			return errors.NewError(msg)
		}

		m.Log().Infof("Migration executed: %s", name)
	}

	return nil
}

// Rollback database.
func (m *Migrator) Rollback(steps ...int) error {
	// Default to 1 step if no value is provided
	s := 1
	if len(steps) > 0 && steps[0] > 1 {
		s = steps[0]
	}

	// Default to max n° migration if steps is higher
	c := m.count()
	if s > c {
		s = c
	}

	m.rollback(s)
	return nil
}

// RollbackAll migrations.
func (m *Migrator) RollbackAll() error {
	return m.rollback(m.count())
}

func (m *Migrator) rollback(steps int) error {
	processed := 0
	count := m.count()

	for i := count - 1; i >= 0; i-- {
		mg := m.steps[i]
		exec := mg.Executor
		idx := exec.GetIndex()
		name := exec.GetName()
		downFx := exec.GetDown()

		// GetErr a new Tx from migrator
		tx, err := m.GetTx()
		if err != nil {
			return errors.Wrap(err, "rollback error")
		}

		// Continue if already applied
		if !m.canApplyRollback(idx, name, tx) {
			m.Log().Infof("Rollback '%s' cannot be executed", name)
			tx.Commit() // No need to handle eventual error here
			continue
		}

		// Pass Tx to the executor
		for _, rfx := range downFx {
			err = rfx(tx)
			if err != nil {
				break
			}
		}

		if err != nil {
			m.Log().Infof("%s rollback not executed", name)
			err2 := tx.Rollback()
			if err2 != nil {
				return errors.Wrap(err2, "rollback rollback error")
			}
			return errors.Wrapf(err, "cannot run rollback '%s'", name)
		}

		exec.SetTx(tx)
		err = m.delMigration(exec)

		err = tx.Commit()
		if err != nil {
			msg := fmt.Sprintf("Cannot delete migration table: %s\n", err.Error())
			m.Log().Errorf("rollback commit error: %s", msg)
			err = tx.Rollback()
			if err != nil {
				return errors.Wrap(err, "rollback rollback error")
			}
			return errors.NewError(msg)
		}

		processed++
		if processed == steps {
			m.Log().Infof("Rollback executed: %s", name)
			return nil
		}
	}

	return nil
}

// SoftReset database.
func (m *Migrator) SoftReset() error {
	err := m.RollbackAll()
	if err != nil {
		log.Printf("Cannot rollback database: %s", err.Error())
		return err
	}

	err = m.Migrate()
	if err != nil {
		log.Printf("Cannot migrate database: %s", err.Error())
		return err
	}

	return nil
}

func (m *Migrator) Reset() error {
	_, err := m.DropDB()
	if err != nil {
		m.Log().Errorf("Drop database error: %w", err)
		// Don't return maybe it was not created before.
	}

	err = m.Migrate()
	if err != nil {
		return errors.Wrap(err, "drop database error")
	}

	return nil
}

// readMigrations records a migration step.
func (m *Migrator) recMigration(e am.MigratorExec) error {
	query := `INSERT INTO %s (id, idx, name, created_at) VALUES ($1, $2, $3, $4);`
	st := fmt.Sprintf(query, migTable)

	uid, err := uuid.NewRandom()
	fmt.Println(uid.String())
	if err != nil {
		return errors.Wrap(err, "cannot update migration table")
	}

	_, err = e.GetTx().Exec(st,
		ToNullString(uid.String()),
		ToNullInt64(e.GetIndex()),
		ToNullString(e.GetName()),
		ToNullTime(time.Now()),
	)

	if err != nil {
		m.Log().Error(err)
		return errors.Wrap(err, "cannot update migration table")
	}

	return nil
}

// cancelRollback returns true if rollback should be cancelled.
func (m *Migrator) cancelRollback(index int64, name string, tx *sql.Tx) bool {
	query := `SELECT (COUNT(*) > 0) AS record_exists FROM %s WHERE idx = %d AND name = '%s'`
	st := fmt.Sprintf(query, migTable, index, name)
	r, err := tx.Query(st)

	if err != nil {
		m.Log().Errorf("Cannot determine rollback status: %w", err)
		return true
	}

	for r.Next() {
		var applied sql.NullBool
		err = r.Scan(&applied)
		if err != nil {
			m.Log().Errorf("Cannot determine migration status: %w", err)
			return true
		}

		return !applied.Bool
	}

	return true
}

// canApplyMigration returns true if migration can be applied.
func (m *Migrator) canApplyMigration(index int64, name string, tx *sql.Tx) bool {
	query := `SELECT (COUNT(*) > 0) AS record_exists FROM %s WHERE idx = %d AND name = '%s'`
	st := fmt.Sprintf(query, migTable, index, name)

	r, err := tx.Query(st)
	defer r.Close()

	if err != nil {
		m.Log().Errorf("Cannot determine migration status: %w", err)
		return false
	}

	for r.Next() {
		var exists sql.NullBool
		err = r.Scan(&exists)
		if err != nil {
			m.Log().Errorf("Cannot determine migration status: %s", err)
			return false
		}

		return !exists.Bool
	}

	return true
}

// canApplyRollback returns true if rollback can be applied.
func (m *Migrator) canApplyRollback(index int64, name string, tx *sql.Tx) bool {
	return !m.canApplyMigration(index, name, tx)
}

func (m *Migrator) delMigration(e am.MigratorExec) error {
	idx := e.GetIndex()
	name := e.GetName()

	query := `DELETE FROM %s WHERE idx = %d AND name = '%s'`
	st := fmt.Sprintf(query, migTable, idx, name)

	_, err := e.GetTx().Exec(st)
	if err != nil {
		return errors.Wrap(err, "cannot delete migration table record")
	}

	return nil
}

// addSteps adds migration steps.
func (m *Migrator) addSteps() error {
	qq, err := m.readMigQueries()
	if err != nil {
		return err
	}

	for i, q := range qq {
		var ups []am.MigFx
		for _, u := range q.Ups {
			ups = append(ups, m.genTxExecFunc(u))
		}

		var downs []am.MigFx
		for _, d := range q.Downs {
			downs = append(downs, m.genTxExecFunc(d))
		}

		step := &am.MigrationStep{
			Index: q.Index,
			Name:  q.Name,
			Up:    ups,
			Down:  downs,
		}

		m.AddMigration(i, step)
	}

	return nil
}

// genTxExecFunc generates a function that let's you execute a query in a transaction.
func (m *Migrator) genTxExecFunc(query string) func(tx *sql.Tx) error {
	return func(tx *sql.Tx) error {
		_, err := tx.Exec(query)
		return err
	}
}

type queries struct {
	Index int64
	Name  string
	Ups   []string
	Downs []string
}

// readMigQueries reads migration queries.
func (m *Migrator) readMigQueries() ([]queries, error) {
	var qq []queries

	files, err := m.fs.ReadDir(m.assetsPath)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".sql") {
			continue
		}

		filePath := fmt.Sprintf("%s/%s", m.assetsPath, file.Name())
		content, err := m.fs.ReadFile(filePath)
		if err != nil {
			return nil, err
		}

		var upStatements []string
		var downStatements []string

		migrationSections := strings.Split(string(content), "--UP")
		if len(migrationSections) < 1 {
			msg := fmt.Sprintf("invalid migrator file format: %s", file.Name())
			return nil, errors.NewError(msg)
		}

		for i, section := range migrationSections {
			section = strings.TrimSpace(section)
			if i == 0 {
				continue
			}

			var up string
			var down string

			parts := strings.Split(section, "--DOWN")
			if len(parts) > 0 {
				up = strings.TrimSpace(parts[0])
				upStatements = append(upStatements, up)
			}

			for i := 1; i < len(parts); i++ {
				down = strings.TrimSpace(parts[i])
				downStatements = append(downStatements, down)
			}
		}

		idx, name := stepName(filePath)

		q := queries{
			Index: idx,
			Name:  name,
			Ups:   upStatements,
			Downs: downStatements,
		}

		qq = append(qq, q)
	}

	return qq, nil
}

func (m *Migrator) count() (last int) {
	return len(m.steps)
}

func (m *Migrator) last() (last int) {
	return m.count() - 1
}

// connectionString returns the connection string.
func (m *Migrator) connectionString() (connString string) {
	cfg := m.Cfg()
	user := cfg.GetString(cfgKey.PgUser)
	pass := cfg.GetString(cfgKey.PgPass)
	name := cfg.GetString(cfgKey.PgDB)
	host := cfg.GetString(cfgKey.PgHost)
	port := cfg.GetInt(cfgKey.PgPort)
	schema := cfg.GetString(cfgKey.PgSchema)
	sslMode := cfg.GetBool(cfgKey.PgSSL)

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d search_path=%s,public", user, pass, name, host, port, schema)

	if sslMode {
		connStr = connStr + " sslmode=enable"
	} else {
		connStr = connStr + " sslmode=disable"
	}

	return connStr
}
