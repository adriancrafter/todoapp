package am

import (
	"fmt"
)

type (
	ResourceWebPath struct{}
)

var (
	WebPath = ResourceWebPath{}
)

type (
	Slugable interface {
		Slug() string
	}
)

// IndexPath returns index path under resource root path.
func (rwp ResourceWebPath) IndexPath() string {
	return ""
}

// EditPath returns edit path under resource root path.
func (rwp ResourceWebPath) EditPath() string {
	return "/{id}/edit"
}

// NewPath returns new path under resource root path.
func (rwp ResourceWebPath) NewPath() string {
	return "/new"
}

// ShowPath returns show path under resource root path.
func (rwp ResourceWebPath) ShowPath() string {
	return "/{id}"
}

// CreatePath returns create path under resource root path.
func (rwp ResourceWebPath) CreatePath() string {
	return ""
}

// UpdatePath returns update path under resource root path.
func (rwp ResourceWebPath) UpdatePath() string {
	return "/{id}"
}

// InitDeletePath returns init delete path under resource root path.
func (rwp ResourceWebPath) InitDeletePath() string {
	return "/{id}/init-delete"
}

// DeletePath returns delete path under resource root path.
func (rwp ResourceWebPath) DeletePath() string {
	return "/{id}"
}

// SignupPath returns signup path.
func (rwp ResourceWebPath) SignupPath() string {
	return "/signup"
}

// LoginPath returns login path.
func (rwp ResourceWebPath) LoginPath() string {
	return "/login"
}

// ResPath returns resource path.
func (rwp ResourceWebPath) ResPath(rootPath string) string {
	return "/" + rootPath + rwp.IndexPath()
}

// ResPathEdit returns resource path edit link
func (rwp ResourceWebPath) ResPathEdit(rootPath string, s Slugable) string {
	return fmt.Sprintf("/%s/%s/edit", rootPath, s.Slug())
}

// ResPathNew returns resource path new link
func (rwp ResourceWebPath) ResPathNew(rootPath string) string {
	return fmt.Sprintf("/%s/new", rootPath)
}

// ResPathInitDelete returns resource path init delete link
func (rwp ResourceWebPath) ResPathInitDelete(rootPath string, s Slugable) string {
	return fmt.Sprintf("/%s/%s/init-delete", rootPath, s.Slug())
}

// ResPathSlug returns resource path  slug link
func (rwp ResourceWebPath) ResPathSlug(rootPath string, s Slugable) string {
	return fmt.Sprintf("/%s/%s", rootPath, s.Slug())
}

// ResAdmin returns resource path under admin path.
func (rwp ResourceWebPath) ResAdmin(path, adminPathPfx string) string {
	return fmt.Sprintf("/%s/%s", adminPathPfx, path)
}
