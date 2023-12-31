package am

import (
	"fmt"
	"net/http"

	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/schema"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

const (
	secCookieName    = "sec-cookie"
	sessionCookieKey = "session"
)

type (
	Controller interface {
		Core
		http.Handler
		Path() string
		Router() *Router
		Handler() http.Handler
		Localizer(r *http.Request) *Localizer
		Session() *sessions.Session
		Store() *sessions.CookieStore
		SecStore() *securecookie.SecureCookie
		InfoMsg() *InfoMessage
		ErrorMsg() *ErrorMessage
	}
)

type (
	SimpleController struct {
		*SimpleCore
		path     string
		router   *Router
		handler  http.Handler
		i18n     *I18N
		tm       *TemplateManager
		session  *sessions.Session
		store    *sessions.CookieStore
		storeKey string
		secStore *securecookie.SecureCookie
		info     *InfoMessage
		error    *ErrorMessage
	}
)

func NewController(name string, parentRouter *Router, path string, opts ...Option) *SimpleController {
	routerName := name + "-router"
	var router *Router
	if parentRouter != nil {
		router = NewSubRouter(routerName, parentRouter, path, opts...)
	} else {
		router = NewRouter(routerName, path, opts...)
	}

	return &SimpleController{
		SimpleCore: NewCore(name, opts...),
		path:       path,
		router:     router,
	}
}

func (sc *SimpleController) SetI18N(i18n *I18N) {
	sc.i18n = i18n
}

func (sc *SimpleController) SetTemplateManager(tm *TemplateManager) {
	sc.tm = tm
}

func (sc *SimpleController) Template() *TemplateManager {
	return sc.TemplateManager()
}

func (sc *SimpleController) TemplateManager() *TemplateManager {
	return sc.tm
}

func (sc *SimpleController) Path() string {
	return sc.path
}

func (sc *SimpleController) Router() *Router {
	return sc.router
}

func (sc *SimpleController) Handler() http.Handler {
	return sc
}

func (sc *SimpleController) InfoMsg() *InfoMessage {
	return sc.info
}

func (sc *SimpleController) ErrorMsg() *ErrorMessage {
	return sc.error
}

func (sc *SimpleController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sc.Router().ServeHTTP(w, r)
}

func (sc *SimpleController) Localizer(r *http.Request) *Localizer {
	return r.Context().Value(I18NCtxKey).(*Localizer)
}

func (sc *SimpleController) NewResponse(w http.ResponseWriter, r *http.Request, data interface{}, errors ValErrorSet) Response {
	//f := MakeFlashSet()

	// Add pending messages
	//f = f.AddItems(sc.RestoreFlash(r))

	res := Response{
		Data:   data,
		Errors: errors,
		Loc:    sc.Localizer(r),
		//Flash: f,
	}

	//r.addSessionData(r)
	//r.AddCSRF(r)

	//sc.ClearFlash(w, r)

	return res
}

// Redirects

func (sc *SimpleController) Redirect(w http.ResponseWriter, r *http.Request, url string) {
	http.Redirect(w, r, url, http.StatusFound)
}

func (sc *SimpleController) RedirectWithFlash(w http.ResponseWriter, r *http.Request, url string, msg string, msgType MsgType) {
	sc.StoreFlash(w, r, msg, msgType)
	http.Redirect(w, r, url, http.StatusSeeOther)
}

func (sc *SimpleController) ErrorRedirect(w http.ResponseWriter, r *http.Request, redirPath, messageID string, err error) {
	m := sc.Localize(r, messageID)
	sc.RedirectWithFlash(w, r, redirPath, m, ErrorMT)
	sc.Log().Error(err)
}

func (sc *SimpleController) Localize(r *http.Request, msgID string) string {
	l := sc.Localizer(r)
	if l == nil {
		sc.Log().Info("No localizer available")
		return msgID
	}

	t, _, err := l.LocalizeWithTag(&i18n.LocalizeConfig{
		MessageID: msgID,
	})

	if err != nil {
		sc.Log().Error(err)
		return msgID
	}

	return t
}

func (sc *SimpleController) localizeMessageID(l *i18n.Localizer, messageID string) (string, error) {
	return l.Localize(&i18n.LocalizeConfig{
		MessageID: messageID,
	})
}

// Form handling related

func (sc *SimpleController) FormToModel(r *http.Request, model interface{}) error {
	d := schema.NewDecoder()
	d.IgnoreUnknownKeys(true)
	return d.Decode(model, r.Form)
}

// NewDecoder build a schema decoder that put values from a map[string][]string
// into a struct.
func (sc *SimpleController) NewDecoder() *schema.Decoder {
	d := schema.NewDecoder()
	d.IgnoreUnknownKeys(true)
	return d
}

// Session related

func (sc *SimpleController) Session() *sessions.Session {
	return sc.session
}

func (sc *SimpleController) Store() *sessions.CookieStore {
	return sc.store
}

func (sc *SimpleController) SecStore() *securecookie.SecureCookie {
	return sc.secStore
}

func (sc *SimpleController) GetSession(r *http.Request, name ...string) *sessions.Session {
	session := SessionKey
	if len(name) > 0 {
		session = name[0]
	}
	s, err := sc.Store().Get(r, session)
	if err != nil {
		sc.Log().Info("Cannot get sesssion from store", "reqID", "n/a")
	}
	return s
}

func (sc *SimpleController) SetSecCookieVals(w http.ResponseWriter, r *http.Request, values map[string]string) {
	c, err := r.Cookie(secCookieName)
	if err != nil {
		sc.Log().Info("No secure cookie present")
		c = &http.Cookie{
			Name: secCookieName,
			Path: "/",
			// TODO: GetErr domain from sc.Cfg
			// Domain: "127.0.0.1",
			// TODO: Dev -> false, Prod -> true
			Secure: false,
		}
	}

	vals := make(map[string]string)

	// Decode the cookie content
	if c.Value != "" {
		err = sc.secStore.Decode(secCookieName, c.Value, &vals)
		if err != nil {
			sc.Log().Info("Cannot decode current secure cookie")
		}
	}

	for k, v := range values {
		// UpdateErr cookie value
		delete(vals, k)

		if v != "" {
			vals[k] = v
		}
	}

	// Encode values again
	e, err := sc.secStore.Encode(secCookieName, vals)
	if err != nil {
		sc.Log().Info("Cannot encode secure cookie")
		return
	}

	c.Value = e

	sc.Log().Info("Storing secure cookie", "vals", vals, "encrypted", e)

	http.SetCookie(w, c)
}

func (sc *SimpleController) SetSecCookieVal(w http.ResponseWriter, r *http.Request, key, value string) {
	c, err := r.Cookie(secCookieName)
	if err != nil {
		sc.Log().Info("No secure cookie present")
		c = &http.Cookie{
			Name: secCookieName,
			Path: "/",
			// TODO: GetErr domain from sc.Cfg
			// Domain: "127.0.0.1",
			// TODO: Dev -> false, Prod -> true
			Secure: false,
		}
	}

	vals := make(map[string]string)

	// Decode the cookie content
	if c.Value != "" {
		err = sc.secStore.Decode(secCookieName, c.Value, &vals)
		if err != nil {
			sc.Log().Info("Cannot decode current secure cookie")
		}
	}

	// UpdateErr cookie value
	delete(vals, key)

	if value != "" {
		vals[key] = value
	}

	// Encode values again
	e, err := sc.secStore.Encode(secCookieName, vals)
	if err != nil {
		sc.Log().Info("Cannot encode secure cookie")
		return
	}

	c.Value = e

	sc.Log().Info("Storing secure cookie", "vals", vals, "encrypted", e)

	http.SetCookie(w, c)
}

func (sc *SimpleController) GetSecCookieValues(r *http.Request, key string) (userData map[string]string, ok bool) {
	c, err := r.Cookie(SecureCookieKey)
	if err != nil {
		sc.Log().Info("No secure cookie")
		sc.Log().Error(err)
		return nil, false
	}

	var vals map[string]string

	// Decode the cookie content
	err = sc.secStore.Decode(secCookieName, c.Value, &vals)
	if err != nil {
		sc.Log().Error(err)
		return vals, false
	}

	sc.Log().Debug("Retrieved from secure cookie", "key", key, "val", fmt.Sprintf("%+v", vals))

	return vals, true
}

func (sc *SimpleController) ReadSecCookieVal(r *http.Request, key string) (val string, ok bool) {
	c, err := r.Cookie(SecureCookieKey)
	if err != nil {
		sc.Log().Info("No secure cookie")
		sc.Log().Error(err)
		return "", false
	}

	var vals map[string]string

	// Decode the cookie content
	err = sc.secStore.Decode(secCookieName, c.Value, &vals)
	if err != nil {
		sc.Log().Error(err)
		return "", false
	}

	val, ok = vals[key]
	if !ok {
		sc.Log().Debug("No value stored in secure cookie", "key", key)
		return val, false
	}

	sc.Log().Debug("Retrieved from secure cookie", "key", key, "val", val)

	return val, true
}

func (sc *SimpleController) ClearSecCookie(w http.ResponseWriter, r *http.Request) {
	c := &http.Cookie{
		Name: secCookieName,
		Path: "/",
		// TODO: GetErr domain from sc.Cfg
		// Domain: "127.0.0.1",
		// TODO: Dev -> false, Prod -> true
		MaxAge: -1,
		Secure: false,
	}
	http.SetCookie(w, c)
}

func (sc *SimpleController) DumpCookieValues(r *http.Request, cookieName, key string) (values string) {
	c, err := r.Cookie(cookieName)
	if err != nil {
		sc.Log().Error(err)
		return "n/a"
	}

	var vals map[string]string

	// Decode the cookie content
	err = sc.secStore.Decode(cookieName, c.Value, &vals)
	if err != nil {
		sc.Log().Error(err)
		return "n/a"
	}

	return fmt.Sprintf("+v", vals)
}

// Flash messages related

func (sc *SimpleController) StoreFlash(w http.ResponseWriter, r *http.Request, message string, mt MsgType) (ok bool) {
	s := sc.GetSession(r)

	// Append to current ones
	f := sc.RestoreFlash(r)
	f = append(f, NewFlashItem(message, mt))

	s.Values[FlashStoreKey] = f
	err := s.Save(r, w)
	if err != nil {
		sc.Log().Error(err)
		return true
	}

	return false
}

func (sc *SimpleController) RestoreFlash(r *http.Request) FlashSet {
	s := sc.GetSession(r)
	v := s.Values[FlashStoreKey]

	f, ok := v.(FlashSet)
	if ok {
		sc.Log().Debugf("Stored flash: %v", spew.Sdump(f))
		return f
	}

	sc.Log().Infof("No stored flash: %v", FlashStoreKey)

	return NewFlashSet()
}

func (sc *SimpleController) ClearFlash(w http.ResponseWriter, r *http.Request) (ok bool) {
	s := sc.GetSession(r)
	delete(s.Values, FlashStoreKey)
	err := s.Save(r, w)
	if err != nil {
		return true
	}
	return false
}
