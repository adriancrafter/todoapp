package am

import (
	"fmt"
)

type (
	ResWebPath struct {
		root string
	}
)

func NewResWebPath(resName string) *ResWebPath {
	return &ResWebPath{
		root: resName,
	}
}

type (
	Slugable interface {
		Slug() string
	}
)

// IndexPath returns index path under resource root path.
func (rwp *ResWebPath) IndexPath() string {
	return ""
}

// EditPath returns edit path under resource root path.
func (rwp *ResWebPath) EditPath() string {
	return "/{slug}/edit"
}

// NewPath returns new path under resource root path.
func (rwp *ResWebPath) NewPath() string {
	return "/new"
}

// ShowPath returns show path under resource root path.
func (rwp *ResWebPath) ShowPath() string {
	return "/{slug}"
}

// CreatePath returns create path under resource root path.
func (rwp *ResWebPath) CreatePath() string {
	return ""
}

// UpdatePath returns update path under resource root path.
func (rwp *ResWebPath) UpdatePath() string {
	return "/{slug}"
}

// InitDeletePath returns init delete path under resource root path.
func (rwp *ResWebPath) InitDeletePath() string {
	return "/{slug}/init-delete"
}

// DeletePath returns delete path under resource root path.
func (rwp *ResWebPath) DeletePath() string {
	return "/{slug}"
}

// SignupPath returns signup path.
func (rwp *ResWebPath) SignupPath() string {
	return "/signup"
}

// SigninPath returns login path.
func (rwp *ResWebPath) SigninPath() string {
	return "/signin"
}

func (rwp *ResWebPath) SignupRoute() string {
	return "/signup"
}

func (rwp *ResWebPath) SigninRoute() string {
	return "/signin"
}

func (rwp *ResWebPath) ResRoute() string {
	return "/" + rwp.root + rwp.IndexPath()
}

func (rwp *ResWebPath) ResSlugRoute() string {
	return "/%s/{slug}" + rwp.root + rwp.IndexPath()
}

func (rwp *ResWebPath) ResEditRoute() string {
	return fmt.Sprintf("/%s/{slug}/edit", rwp.root)
}

func (rwp *ResWebPath) ResNewRoute() string {
	return fmt.Sprintf("/%s/new", rwp.root)
}

func (rwp *ResWebPath) ResInitDeleteRoute() string {
	return fmt.Sprintf("/%s/{slug}/init-delete", rwp.root)
}

func (rwp *ResWebPath) ResForceDeleteRoute() string {
	return fmt.Sprintf("/%s/{slug}/force-delete", rwp.root)
}

func (rwp *ResWebPath) ResPurgeRoute() string {
	return fmt.Sprintf("/%s/purge", rwp.root)
}

func (rwp *ResWebPath) ResSigninRoute() string {
	return fmt.Sprintf("/%s/signin", rwp.root)
}

func (rwp *ResWebPath) ResSignupRoute() string {
	return fmt.Sprintf("/%s/signup", rwp.root)
}
