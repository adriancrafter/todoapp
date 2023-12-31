package am

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

type (
	Router struct {
		*SimpleCore
		*mux.Router
		Path       string
		i18nBundle *i18n.Bundle
	}
)

func NewRouter(name string, path string, opts ...Option) *Router {
	return &Router{
		SimpleCore: NewCore(name, opts...),
		Router:     mux.NewRouter(),
		Path:       path,
	}
}

func NewSubRouter(name string, parent *Router, pathPrefix string, opts ...Option) *Router {
	router := parent.PathPrefix(pathPrefix).Subrouter()

	return &Router{
		SimpleCore: NewCore(name, opts...),
		Router:     router,
		Path:       pathPrefix,
	}
}

func (rt *Router) I18N(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		lang := r.FormValue("lang")
		accept := r.Header.Get("Accept-Language")
		bundle := rt.I18NBundle()
		l := i18n.NewLocalizer(bundle, lang, accept)
		ctx := context.WithValue(r.Context(), I18NCtxKey, l)
		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(fn)
}

func (rt *Router) I18NBundle() *i18n.Bundle {
	//if rt.i18nBundle != nil {
	//	return rt.i18nBundle
	//}
	//
	b := i18n.NewBundle(language.English)
	//b.RegisterUnmarshalFunc("json", json.Unmarshal)
	//
	//locales := []string{"en", "pl", "de", "es"}
	//
	//for _, loc := range locales {
	//	path := fmt.Sprintf("/assets/web/embed/locale/%s.json", loc)
	//
	//	// Open pkger filer
	//	f, err := pkger.Open(path)
	//	if err != nil {
	//		rt.Log.InfoMsg("Opening embedded resource", "path", path)
	//		rt.Log.ErrorMsg(err)
	//	}
	//	defer f.Close()
	//
	//	// Read file content
	//	fs, err := f.Stat()
	//	if err != nil {
	//		rt.Log.InfoMsg("Stating embedded resource", "path", path)
	//		rt.Log.ErrorMsg(err)
	//	}
	//	fd := make([]byte, fs.Size())
	//	f.Read(fd)
	//
	//	// Load into bundle
	//	b.ParseMessageFileBytes(fd, path)
	//	//b.LoadEmbeddedMessageFile(fd, path)
	//}
	//
	//// Cache it
	//rt.i18nBundle = b

	return b
}
