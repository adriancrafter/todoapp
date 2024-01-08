package auth

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"

	"github.com/adriancrafter/todoapp/internal/am"
	"github.com/adriancrafter/todoapp/internal/am/errors"
)

const (
	authPath       = "/auth"
	authController = "auth"
	authSignin     = "signin"
	authSignup     = "signup"
)

const (
	userRes = "user"
	authRes = "auth"
)

type (
	WebController struct {
		*am.SimpleController
		up  *am.ResWebPath // user resource web path
		ap  *am.ResWebPath // auth resource web path
		svc Service
	}
)

func NewWebController(parent *am.Router, svc Service, opts ...am.Option) *WebController {
	name := fmt.Sprintf("%s-web-controller", authController)
	c := &WebController{
		SimpleController: am.NewController(name, parent, authPath, opts...),
		up:               am.NewResWebPath(userRes),
		ap:               am.NewResWebPath(authRes),
		svc:              svc,
	}

	return c
}

func (c *WebController) Setup(ctx context.Context) error {
	c.routes()
	return nil
}

func (c *WebController) UserIndex(w http.ResponseWriter, r *http.Request)      {}
func (c *WebController) UserShow(w http.ResponseWriter, r *http.Request)       {}
func (c *WebController) UserCreate(w http.ResponseWriter, r *http.Request)     {}
func (c *WebController) UserUpdate(w http.ResponseWriter, r *http.Request)     {}
func (c *WebController) UserPreDelete(w http.ResponseWriter, r *http.Request)  {}
func (c *WebController) UserSoftDelete(w http.ResponseWriter, r *http.Request) {}
func (c *WebController) UserDelete(w http.ResponseWriter, r *http.Request)     {}
func (c *WebController) UserPurge(w http.ResponseWriter, r *http.Request)      {}

func (c *WebController) UserInitSignin(w http.ResponseWriter, r *http.Request) {
	userVM := UserVM{}

	res := c.NewResponse(w, r, userVM, nil)
	res.SetAction(c.userSignInAction())

	t, err := c.Template().Get(authController, authSignin)
	if err != nil {
		c.ErrorRedirect(w, r, c.AuthPath(), c.ErrorMsg().ProcessErr, err)
		return
	}

	err = t.Execute(w, res)
	if err != nil {
		c.ErrorRedirect(w, r, c.AuthPath(), c.ErrorMsg().ProcessErr, err)
		return
	}
}

func (c *WebController) UserSignin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Decode request data into a form.
	signinVM := SigninVM{}
	err := c.FormToModel(r, &signinVM)
	if err != nil {
		c.ErrorRedirect(w, r, c.AuthPath(), c.ErrorMsg().ProcessErr, err)
		return
	}

	//// GetErr IP from user request
	//// ip := "0.0.0.0/24"
	//// TODO: Provide IP to the service in order to register last IP
	//// Can be used to detect spurious logins.
	//// user, err := c.MainService().SignInUser(signinVM.Username, signinVM.Password, ip)
	user, err := c.Service().SignInUser(ctx, signinVM)
	if err != nil {
		c.ErrorRedirect(w, r, c.AuthPath(), c.ErrorMsg().ProcessErr, err)
		return
	}

	c.Log().Debugf("user %s signed in", user.Username)
	//if err != nil {
	//	msgID := c.ErrorMsg().SigninErr
	//
	//	// Give a hint to user about kind of error.
	//	if err == svc.CredentialsErr {
	//		msgID = (err.(svc.Err)).MsgID()
	//		c.rerenderUserForm(w, r, user.ToForm(), nil, kbs.SignInTmpl, userSignInAction())
	//		return
	//	}
	//
	//	c.ErrorRedirect(w, r, UserPath(), msgID, err)
	//	return
	//}
	//
	//// Register user data in secure session cookie.
	//userData := map[string]string{
	//	"slug":        user.Slug.String,
	//	"username":    user.Username.String,
	//	"permissions": user.PermissionTags.String,
	//}
	//
	//c.SignIn(w, r, userData)
	//c.Log.Debug("UserDA signed in", "user", user.Username.String)
	//
	//// Localize Ok info message, put it into a flash message
	//// and redirect to index.
	//m := c.Localize(r, c.InfoMsg().SignedInMsg)
	//c.RedirectWithFlash(w, r, UserPath(), m, am.InfoMT)
}

func (c *WebController) rerenderUserForm(w http.ResponseWriter, r *http.Request,
	data interface{},
	valErrors am.ValErrorSet,
	handlerTemplate string,
	action am.FormAction) {

	res := c.NewResponse(w, r, data, valErrors)
	res.AddErrorFlash(c.ErrorMsg().InputValuesErr)
	res.SetAction(action)

	ts, err := c.Template().Get(authController, handlerTemplate)
	if err != nil {
		c.ErrorRedirect(w, r, c.UserPath(), c.ErrorMsg().InputValuesErr, err)
		return
	}

	err = ts.Execute(w, res)
	if err != nil {
		c.ErrorRedirect(w, r, c.UserPath(), c.ErrorMsg().ProcessErr, err)
		return
	}

	return
}

func (c *WebController) UserInitSignup(w http.ResponseWriter, r *http.Request) {
	//userVM := vm.UserDA{}
	//
	//res := c.NewResponse(w, r, userVM, nil)
	//res.SetAction(userSignUpAction())
	//
	//err := c.templates.Execute(w, res)
	//if err != nil {
	//	fmt.Println(err)
	//	http.ErrorMsg(w, "ErrorMsg rendering sign-in page", http.StatusInternalServerError)
	//}
}

func (c *WebController) UserSignup(w http.ResponseWriter, r *http.Request) {
	// Parse form values
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//username := r.FormValue("username")
	//password := r.FormValue("password")

	// CreateErr user
	//err = c.CreateUser(username, password)
	//if err != nil {
	//	http.ErrorMsg(w, "Failed to create user", http.StatusInternalServerError)
	//	return
	//}

	// Automatically sign in the user and create a session

	http.Redirect(w, r, "/signin", http.StatusSeeOther)
}

// Handler interface

func (c *WebController) Service() Service {
	return c.svc
}

// Helpers

func (c *WebController) User(r *http.Request) (userID uuid.UUID, err error) {
	panic("not implemented yet")
}

func (c *WebController) closeBody(body io.ReadCloser) {
	if err := body.Close(); err != nil {
		c.Log().Error(errors.Wrap(err, "failed to close body"))
	}
}

// Paths
func (c *WebController) UserPath() string {
	return c.UserPath()
}

func (c *WebController) AuthPath() string {
	return c.AuthPath()
}

// Form Actions

// userCreateAction returns a new form action for user creatiom.
func (c *WebController) userCreateAction() am.FormAction {
	return am.NewFormAction(c.UserPath(), am.POST)
}

// userUpdateAction returns a new form action for user update.
func (c *WebController) userUpdateAction(model am.Slugable) am.FormAction {
	// TODO: This should return a path with identifier (slug) in it not a route.
	// Implement ResSlugPath(slug)
	// c.up.ResSlugRoute() => c.up.ResSlugPath(slug)
	return am.NewFormAction(c.up.ResSlugRoute(), am.PUT)
}

// userDeleteAction returns a new form action for user deletion.
func (c *WebController) userDeleteAction(model am.Slugable) am.FormAction {
	// TODO: This should return a path with identifier (slug) in it not a route.
	// Implement ResSlugPath(slug)
	// c.up.ResSlugRoute() => c.up.ResSlugPath(slug)
	return am.NewFormAction(c.up.ResSlugRoute(), am.DELETE)
}

// userSignUpAction returns a new form action for user sign up.
func (c *WebController) userSignUpAction() am.FormAction {
	return am.NewFormAction(c.ap.ResSignupRoute(), am.POST)
}

// userSignInAction returns a new form action for user sign in.
func (c *WebController) userSignInAction() am.FormAction {
	return am.NewFormAction(c.ap.ResSigninRoute(), am.POST)
}
