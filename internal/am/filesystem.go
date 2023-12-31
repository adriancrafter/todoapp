package am

import (
	"embed"
)

const (
	MigFSKey      = "migFS"
	SeedFSKey     = "seedFS"
	QueryFSKey    = "queryFS"
	TemplateFSKey = "templateFS"
	WebFSKey      = "webFS"
	LocaleFSKey   = "localeFS"
)

type Filesystems map[string]embed.FS

func NewFilesystems() Filesystems {
	return Filesystems{}
}

func (fs Filesystems) Get(name string) (f embed.FS, ok bool) {
	f, ok = fs[name]
	return f, ok
}

func (fs Filesystems) Add(name string, f embed.FS) {
	fs[name] = f
}
