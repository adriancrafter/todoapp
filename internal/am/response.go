package am

import (
	"strings"
)

const (
	GET    = "GET"
	POST   = "POST"
	PATCH  = "PATCH"
	PUT    = "PUT"
	DELETE = "DELETE"
	HEAD   = "HEAD"
)

type (
	// Response stands for wrapped response
	Response struct {
		// Data stores model data
		Data interface{}
		// Stores model detailed validation errors.
		Errors ValErrorSet
		// Action can be used to reuse form templates letting change target and method from controller.
		Action FormAction
		// Loc can be used to show msgID in different languages
		Loc *Localizer
		// Flash messages
		Flash FlashSet
		// Session data
		SessionData map[string]string
		// Cross-site request forgery protection
		CSRF map[string]interface{}
	}
)

const (
	InfoMT  MsgType = "info"
	WarnMT  MsgType = "warn"
	ErrorMT MsgType = "error"
	DebugMT MsgType = "debug"
)

func (res *Response) SetAction(fa FormAction) {
	res.Action = fa
}

func (r *Response) AddInfoFlash(infoMsg string) {
	f := r.Flash

	m := strings.Trim(infoMsg, " ")
	if m != "" {
		f = f.AddItem(NewFlashItem(m, InfoMT))
	}

	r.Flash = f
}

func (r *Response) AddWarnFlash(warnMsgs ...string) {
	f := r.Flash

	if len(warnMsgs) > 0 {
		for _, m := range warnMsgs {
			f = f.AddItem(NewFlashItem(m, WarnMT))
		}
	}

	r.Flash = f
}

func (r *Response) AddErrorFlash(errorMsg string) {
	f := r.Flash

	m := strings.Trim(errorMsg, " ")
	if m != "" {
		f = f.AddItem(NewFlashItem(m, ErrorMT))
	}

	r.Flash = f
}

type FormAction struct {
	Target string
	Method string
}

func NewFormAction(target, method string) FormAction {
	return FormAction{Target: target, Method: method}
}
