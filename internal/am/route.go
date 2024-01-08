package am

import (
	"fmt"
)

type (
	// Route functions that can be used to generate app standard paths to be
	// used in controller handlers routes definition.
	Route struct {
		root string
	}
)

func NewRoute(resName string) *Route {
	return &Route{
		root: resName,
	}
}

func (rwp *Route) Index() string {
	return fmt.Sprintf("/%s", rwp.root)

}

func (rwp *Route) Slug() string {
	return fmt.Sprintf("/%s/{slug}", rwp.root)
}

func (rwp *Route) Edit() string {
	return fmt.Sprintf("/%s/{slug}/edit", rwp.root)
}

func (rwp *Route) New() string {
	return fmt.Sprintf("/%s/new", rwp.root)
}

func (rwp *Route) InitDelete() string {
	return fmt.Sprintf("/%s/{slug}/init-delete", rwp.root)
}

func (rwp *Route) ForceDelete() string {
	return fmt.Sprintf("/%s/{slug}/force-delete", rwp.root)
}

func (rwp *Route) Purge() string {
	return fmt.Sprintf("/%s/purge", rwp.root)
}

func (rwp *Route) Signup() string {
	return "/signup"
}

func (rwp *Route) Signin() string {
	return "/signin"
}
