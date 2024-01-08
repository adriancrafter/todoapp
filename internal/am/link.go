package am

import (
	"fmt"
)

type (
	// ResLink functions that can be used in Go templates and form actions to
	// generate resource related paths.
	ResLink struct {
		root string
	}
)

func NewResLink(resName string) *ResLink {
	return &ResLink{
		root: resName,
	}
}

type (
	Slugable interface {
		Slug() string
	}
)

func (rl *ResLink) Index() string {
	return fmt.Sprintf("/%s", rl.root)
}

func (rl *ResLink) Slug(s Slugable) string {
	return fmt.Sprintf("/%s/%s", rl.root, s.Slug())
}

func (rl *ResLink) Edit(s Slugable) string {
	return fmt.Sprintf("/%s/%s/edit", rl.root, s.Slug())
}

func (rl *ResLink) New() string {
	return fmt.Sprintf("/%s/new", rl.root)
}

func (rl *ResLink) Delete(s Slugable) string {
	return fmt.Sprintf("/%s/%s/init-delete", s.Slug())
}

func (rl *ResLink) ForceDelete(s Slugable) string {
	return fmt.Sprintf("/%s/%s/force-delete", s.Slug())
}

func (rl *ResLink) Purge() string {
	return fmt.Sprintf("/%s/purge", rl.root)
}

func (rl *ResLink) Signin() string {
	return fmt.Sprintf("/%s/signin", rl.root)
}

func (rl *ResLink) Signup() string {
	return fmt.Sprintf("/%s/signup", rl.root)
}
