package main

import (
	"embed"
	"os"

	"github.com/adriancrafter/todoapp/internal/am"
	a "github.com/adriancrafter/todoapp/internal/app"
)

const (
	name     = "todoapp"
	env      = "todoapp"
	logLevel = "debug"
)

var (
	log am.Logger = am.NewLogger(logLevel)
)

var (
	fs am.Filesystems = am.NewFilesystems()

	//go:embed all:assets/migrations/pg/*.sql
	migFS embed.FS

	//go:embed all:assets/seeding/pg/*.sql
	seedFS embed.FS

	//go:embed all:assets/queries/pg/*.sql
	queriesFS embed.FS

	//go:embed all:assets/templates/web/*
	templatesFS embed.FS

	//go:embed all:assets/web/static/**/*.css all:assets/web/static/**/*.png
	staticFS embed.FS

	//go:embed all:assets/locale/*.json
	localeFS embed.FS
)

func main() {
	app, err := a.NewApp(name, log)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	fs.Add(am.MigFSKey, migFS)
	fs.Add(am.SeedFSKey, seedFS)
	fs.Add(am.QueryFSKey, queriesFS)
	fs.Add(am.TemplateFSKey, templatesFS)
	fs.Add(am.WebFSKey, staticFS)
	fs.Add(am.LocaleFSKey, localeFS)
	app.SetFileSystems(fs)

	err = app.Run()
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
}
