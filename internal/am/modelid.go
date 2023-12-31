package am

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/google/uuid"
)

type (
	ID struct {
		val      uuid.UUID
		TenantID uuid.UUID
		Slug     string
	}
)

func NewID(tenantID uuid.UUID, uid ...uuid.UUID) ID {
	var id uuid.UUID
	if len(uid) > 0 {
		id = uid[0]
	}

	return ID{
		val:      id,
		TenantID: tenantID,
	}
}

func NewCustomID(tenantID uuid.UUID, id uuid.UUID, slug string) ID {
	return ID{
		val:      id,
		TenantID: tenantID,
		Slug:     slug,
	}
}

func (id *ID) GenID(provided ...uuid.UUID) {
	if id.val != uuid.Nil {
		return // A value has been previously assigned
	}

	if len(provided) > 0 {
		id.val = provided[0] // If value is provided, use it
		return
	}

	id.val = uuid.New()
}

func (id *ID) Val() uuid.UUID {
	return id.val
}

func (id *ID) String() string {
	return id.val.String()
}

func (id *ID) IsNew() bool {
	return id.Val() == ZeroUUID
}

func (id *ID) UpdateSlug(prefix string) (slug string) {
	if strings.Trim(id.Slug, " ") == "" {
		s, err := id.genSlug(prefix)
		if err != nil {
			return ""
		}

		id.Slug = s
	}

	return id.Slug
}

func (id *ID) genSlug(prefix string) (slug string, err error) {
	if strings.TrimSpace(prefix) == "" {
		return "", NoSlugPrefixErr
	}

	prefix = strings.Replace(prefix, "-", "", -1)
	prefix = strings.Replace(prefix, "_", "", -1)

	if !utf8.ValidString(prefix) {
		v := make([]rune, 0, len(prefix))
		for i, r := range prefix {
			if r == utf8.RuneError {
				_, size := utf8.DecodeRuneInString(prefix[i:])
				if size == 1 {
					continue
				}
			}
			v = append(v, r)
		}
		prefix = string(v)
	}

	prefix = strings.ToLower(prefix)

	s := strings.Split(uuid.New().String(), "-")
	l := s[len(s)-1]

	return strings.ToLower(fmt.Sprintf("%s-%s", prefix, l)), nil
}

func (id *ID) SetCreateValues(slugPrefix string) {
	id.GenID()
	id.UpdateSlug(slugPrefix)
}

func GenTag() string {
	s := strings.Split(uuid.New().String(), "-")
	l := s[len(s)-1]

	return strings.ToUpper(l[4:12])
}

func (id *ID) Match(tc ID) (ok bool) {
	ok = id.Val() == tc.Val() &&
		id.TenantID == tc.TenantID &&
		id.Slug == tc.Slug
	return ok
}
