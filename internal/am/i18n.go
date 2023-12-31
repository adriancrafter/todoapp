package am

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

const (
	localesPath = "assets/locale/%s.json"
)

var (
	lock         = &sync.Mutex{}
	i18nInstance *I18N
)

type (
	I18N struct {
		*SimpleCore
		fs      embed.FS
		locales []string // list of locales to load
		bundle  *i18n.Bundle
	}
)

func NewI18N(localeFS embed.FS, locales []string, opts ...Option) *I18N {
	if i18nInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if i18nInstance == nil {
			i18nInstance = &I18N{
				SimpleCore: NewCore("i18n", opts...),
				fs:         localeFS,
				locales:    locales,
			}

			i18nInstance.load()
		}
	}

	return i18nInstance
}

func NewLocalizer(bundle *i18n.Bundle, lang ...string) (loc *Localizer) {
	return &Localizer{
		Localizer: i18n.NewLocalizer(bundle, lang...),
	}
}

func (l *Localizer) Localize(textID string) string {
	if l.Localizer != nil {
		t, _, err := l.LocalizeWithTag(&i18n.LocalizeConfig{
			MessageID: textID,
		})

		if err != nil {
			return fmt.Sprintf("%s [msg. ID]", textID)
		}

		return t
	}

	return fmt.Sprintf("%s [msg. ID]", textID)
}

func (i *I18N) Bundle() *i18n.Bundle {
	if i.bundle != nil {
		return i.bundle
	}

	i.load()

	return i.bundle
}

func (i *I18N) load() {
	b := i18n.NewBundle(language.English)
	b.RegisterUnmarshalFunc("json", json.Unmarshal)

	for _, loc := range i.locales {
		path := fmt.Sprintf(localesPath, loc)

		_, err := b.LoadMessageFileFS(i.fs, path)
		if err != nil {
			i.Log().Errorf("error loading bundle file %s: %v", path, err)
		}
	}

	i.bundle = b

}

func (i *I18N) Middleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		lang := r.FormValue("lang")
		accept := r.Header.Get("Accept-Language")
		bundle := i.bundle
		l := NewLocalizer(bundle, lang, accept)
		ctx := context.WithValue(r.Context(), I18NCtxKey, l)
		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(fn)
}

type (
	Localizer struct {
		*i18n.Localizer
	}
)
