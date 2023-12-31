package am

var Key = newCfgKeyReg()

func newCfgKeyReg() *cfgKeyReg {
	return &cfgKeyReg{
		// Logging
		LogOutputCompact: "log.output.compact",

		// Web Server

		WebServerHost:     "http.web.server.host",
		WebServerPort:     "http.web.server.port",
		WebServerTimeout:  "http.web.server.shutdown.timeout.secs",
		WebErrorExposeInt: "http.errors.expose.internal",

		// API Server

		APIServerHost:     "http.api.server.host",
		APIServerPort:     "http.api.server.port",
		APIServerTimeout:  "http.api.server.shutdown.timeout.secs",
		APIErrorExposeInt: "http.api.errors.expose.internal",

		// Postgres

		SQLiteUser:     "db.sqlite.user",
		SQLitePass:     "db.sqlite.pass",
		SQLiteDB:       "db.sqlite.database",
		SQLiteHost:     "db.sqlite.host",
		SQLiteSchema:   "db.sqlite.schema",
		SQLiteFilePath: "db.sqlite.filepath",

		// Postgres

		PgUser:   "db.pg.user",
		PgPass:   "db.pg.pass",
		PgDB:     "db.pg.database",
		PgHost:   "db.pg.host",
		PgPort:   "db.pg.port",
		PgSchema: "db.pg.schema",
		PgSSL:    "db.pg.sslmode",
	}
}

type cfgKeyReg struct {
	// Logging
	LogOutputCompact string

	// Web Server

	WebServerHost     string
	WebServerPort     string
	WebServerTimeout  string
	WebErrorExposeInt string

	// API Server

	APIServerHost     string
	APIServerPort     string
	APIServerTimeout  string
	APIErrorExposeInt string

	// SQLite

	SQLiteUser     string
	SQLitePass     string
	SQLiteDB       string
	SQLiteHost     string
	SQLitePort     string
	SQLiteSchema   string
	SQLiteSSL      string
	SQLiteFilePath string

	// Postgres

	PgUser   string
	PgPass   string
	PgDB     string
	PgHost   string
	PgPort   string
	PgSchema string
	PgSSL    string
}
