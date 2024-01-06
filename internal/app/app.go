package app

import (
	"context"
	"embed"
	"fmt"
	"sync"

	//"github.com/adriancrafter/todoapp/internal/api"
	"github.com/adriancrafter/todoapp/internal/am"
	"github.com/adriancrafter/todoapp/internal/am/db/pg"
	"github.com/adriancrafter/todoapp/internal/am/errors"
	"github.com/adriancrafter/todoapp/internal/auth"
)

const (
	launchError   = "app launch error"
	runError      = "app run error"
	setupError    = "app setup error"
	startError    = "app start error"
	stopError     = "app stop error"
	shutdownError = "app shutdown error"
)

var (
	// locales are the enabled locales.
	// Will be loaded from config at some point.
	locales = []string{"en", "de", "es", "pl"}
)

type App struct {
	sync.Mutex
	am.Core
	opts       []am.Option
	migFS      embed.FS
	seedFS     embed.FS
	queryFS    embed.FS
	templateFS embed.FS
	webFS      embed.FS
	localeFS   embed.FS
	fs         am.Filesystems
	spv        am.Supervisor
	db         am.DB
	qm         *am.QueryManager
	router     *am.Router
	authRepo   auth.Repo
	authSvc    auth.Service
	authWeb    *auth.WebController // TODO: Use WebController interface
	//authAPI    *auth.APIController
	//todoRepo todo.Repo
	//todoSvc  todo.Service
	//todoWeb  *todo.WebController
	//todoAPI  *todo.APIController
	//todoWeb  *todo.WebController
	migrator am.Migrator
	seeder   am.Seeder
	i18n     *am.I18N
	tm       *am.TemplateManager
	web      *Server
}

// NewApp creates a new app instance
func NewApp(name string, log am.Logger) (app *App, err error) {
	cfg, err := am.NewConfig(name).Load()
	if err != nil {
		return nil, errors.Wrap(err, launchError)
	}

	opts := []am.Option{
		am.WithConfig(cfg),
		am.WithLogger(log),
	}

	app = &App{
		Core: am.NewCore(name, opts...),
		opts: opts,
	}

	return app, nil
}

// Run runs the app
func (app *App) Run() (err error) {
	ctx := context.Background()

	err = app.Setup(ctx)
	if err != nil {
		return errors.Wrap(err, runError)
	}

	return app.Start(ctx)
}

// Setup sets up the app
func (app *App) Setup(ctx context.Context) error {
	app.EnableSupervisor()

	// Migrator
	app.migrator = pg.NewMigrator(app.migFS, app.opts...)

	// Seeder
	app.seeder = pg.NewSeeder(app.seedFS, app.opts...)

	// Databases
	app.db = pg.NewDB(app.opts...)
	err := app.db.Setup(ctx)
	if err != nil {
		return errors.Wrap(err, setupError)
	}

	// QueryManager manager
	app.qm = am.NewQueryManager(app.queryFS, "pg", app.opts...)
	app.qm.Setup(ctx)

	// NewI18N
	app.i18n = am.NewI18N(app.localeFS, locales, app.opts...)

	// Templates
	app.tm = am.NewTemplateManager(app.templateFS, false, app.opts...)
	err = app.tm.Setup(ctx)
	if err != nil {
		return errors.Wrap(err, setupError)
	}

	// Repos
	app.authRepo = auth.NewRepo(app.db, app.qm, app.opts...)
	//app.todoRepo = todo.NewRepo(app.db, app.qm, app.opts...)

	// Services
	app.authSvc = auth.NewService(app.authRepo, app.opts...)
	//app.todoSvc = todo.NewService(app.todoRepo, app.authRepo, app.opts...)

	// Router
	app.router = am.NewRouter("web-router", "/", app.opts...)

	// Controllers
	app.authWeb = auth.NewWebController(app.router, app.authSvc, app.opts...)
	app.authWeb.SetTemplateManager(app.tm)

	// Web server
	app.web = NewServer(app.router, app.fs, app.i18n, app.opts...)
	app.web.SetAuthController(app.authWeb)
	err = app.web.Setup(ctx)
	if err != nil {
		return errors.Wrap(err, setupError)
	}

	//app.web.SetTodoController(app.todoWeb)

	// API server

	// gRPC servers

	// Event bus

	// WIP: to avoid unused var message
	//app.Log().Debugf("MainRepo: %v", repo)
	//app.Log().Debugf("MainService: %v", app.todoSvc)

	return nil
}

// Start starts the app
func (app *App) Start(ctx context.Context) error {
	app.Log().Infof("%s starting...", app.Name())
	defer app.Log().Infof("%s stopped", app.Name())

	var err error

	err = app.migrator.Start(ctx)
	if err != nil {
		return errors.Wrap(err, "app start error")
	}

	err = app.seeder.Start(ctx)
	if err != nil {
		return errors.Wrap(err, "app start error")
	}

	err = app.authSvc.Start(ctx)
	if err != nil {
		return errors.Wrap(err, "app start error")
	}

	app.spv.AddTasks(
		app.web.Start,
		//app.api.Start,
		//app.grpc.Start,
	)

	app.Log().Infof("%s started!", app.Name())

	return app.spv.Wait()
}

// Stop stops the app
func (app *App) Stop(ctx context.Context) error {
	return nil
}

// Shutdown shuts down the app
func (app *App) Shutdown(ctx context.Context) error {
	return nil
}

// EnableSupervisor enables the supervisor
func (app *App) EnableSupervisor() {
	name := fmt.Sprintf("%s-supervisor", app.Name())
	app.spv = am.NewSupervisor(name, true, app.opts)
}

// SetupRoutes sets up the app routes
func (app *App) SetupRoutes(ctx context.Context) (err error) {
	return app.web.Setup(ctx)
}

// SetupProbes WIP implementation
func (app *App) SetupProbes(ctx context.Context) {
	app.web.Router().Handle("/healthz", am.HealthzProbe)
}

// File Systems

func (app *App) SetFileSystems(fs am.Filesystems) {
	app.fs = fs
	for name, f := range fs {
		switch name {
		case am.MigFSKey:
			app.migFS = f
		case am.SeedFSKey:
			app.seedFS = f
		case am.QueryFSKey:
			app.queryFS = f
		case am.TemplateFSKey:
			app.templateFS = f
		case am.WebFSKey:
			app.webFS = f
		case am.LocaleFSKey:
			app.localeFS = f
		}
	}
}

// SetMigratorFS sets a filesystem for migrations
func (app *App) SetMigratorFS(fs embed.FS) {
	app.migFS = fs
}

// SetSeederFS sets a filesystem for seeders
func (app *App) SetSeederFS(fs embed.FS) {
	app.seedFS = fs
}

// SetQueriesFS sets a filesystem for queries
func (app *App) SetQueriesFS(fs embed.FS) {
	app.queryFS = fs
}

// SetTemplatesFS sets a filesystem for templates
func (app *App) SetTemplatesFS(fs embed.FS) {
	app.templateFS = fs
}

// SetStaticFS sets a filesystem for web static files
func (app *App) SetStaticFS(fs embed.FS) {
	app.webFS = fs
}

// SetLocaleFS sets a filesystem for locale files
func (app *App) SetLocaleFS(fs embed.FS) {
	app.localeFS = fs
}
