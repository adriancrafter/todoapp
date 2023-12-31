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
)

type (
	MainWebController struct {
		*am.SimpleController
		svc Service
	}
)

func NewWebController(parent *am.Router, svc Service, opts ...am.Option) *MainWebController {
	name := fmt.Sprintf("%s-web-controller", authController)
	c := &MainWebController{
		SimpleController: am.NewController(name, parent, authPath, opts...),
		svc:              svc,
	}

	return c
}

func (c *MainWebController) Setup(ctx context.Context) error {
	c.routes()
	return nil
}

func (c *MainWebController) UserIndex(w http.ResponseWriter, r *http.Request)      {}
func (c *MainWebController) UserShow(w http.ResponseWriter, r *http.Request)       {}
func (c *MainWebController) UserCreate(w http.ResponseWriter, r *http.Request)     {}
func (c *MainWebController) UserUpdate(w http.ResponseWriter, r *http.Request)     {}
func (c *MainWebController) UserPreDelete(w http.ResponseWriter, r *http.Request)  {}
func (c *MainWebController) UserSoftDelete(w http.ResponseWriter, r *http.Request) {}
func (c *MainWebController) UserDelete(w http.ResponseWriter, r *http.Request)     {}
func (c *MainWebController) UserPurge(w http.ResponseWriter, r *http.Request)      {}

func (c *MainWebController) UserInitSignin(w http.ResponseWriter, r *http.Request) {
	userVM := UserVM{}

	res := c.NewResponse(w, r, userVM, nil)
	res.SetAction(userSignInAction())

	t, err := c.Template().Get(authController, authSignin)
	if err != nil {
		c.ErrorRedirect(w, r, AuthPath(), c.ErrorMsg().ProcessErr, err)
		return
	}

	err = t.Execute(w, res)
	if err != nil {
		c.ErrorRedirect(w, r, AuthPath(), c.ErrorMsg().ProcessErr, err)
		return
	}
}

func (c *MainWebController) UserSignin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Decode request data into a form.
	userVM := UserVM{}
	err := c.FormToModel(r, &userVM)
	if err != nil {
		c.ErrorRedirect(w, r, AuthPath(), c.ErrorMsg().ProcessErr, err)
		return
	}

	//// GetErr IP from user request
	//// ip := "0.0.0.0/24"
	//// TODO: Provide IP to the service in order to register last IP
	//// Can be used to detect spurious logins.
	//// user, err := c.MainService().SignInUser(userVM.Username, userVM.Password, ip)
	user, err := c.Service().SignInUser(ctx, userVM)
	if err != nil {
		c.ErrorRedirect(w, r, AuthPath(), c.ErrorMsg().ProcessErr, err)
		return
	}

	c.Log().Debug("UserDA signed in", "user", user.Username)
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

func (c *MainWebController) rerenderUserForm(w http.ResponseWriter, r *http.Request,
	data interface{},
	valErrors am.ValErrorSet,
	handlerTemplate string,
	action am.FormAction) {

	res := c.NewResponse(w, r, data, valErrors)
	res.AddErrorFlash(c.ErrorMsg().InputValuesErr)
	res.SetAction(action)

	ts, err := c.Template().Get(authController, handlerTemplate)
	if err != nil {
		c.ErrorRedirect(w, r, UserPath(), c.ErrorMsg().InputValuesErr, err)
		return
	}

	err = ts.Execute(w, res)
	if err != nil {
		c.ErrorRedirect(w, r, UserPath(), c.ErrorMsg().ProcessErr, err)
		return
	}

	return
}

func (c *MainWebController) UserInitSignup(w http.ResponseWriter, r *http.Request) {
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

func (c *MainWebController) UserSignup(w http.ResponseWriter, r *http.Request) {
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

func (c *MainWebController) Service() Service {
	return c.svc
}

// Helpers

func (c *MainWebController) User(r *http.Request) (userID uuid.UUID, err error) {
	panic("not implemented yet")
}

func (c *MainWebController) closeBody(body io.ReadCloser) {
	if err := body.Close(); err != nil {
		c.Log().Error(errors.Wrap(err, "failed to close body"))
	}
}

// Form Actions

// userCreateAction returns a new form action for user creatiom.
func userCreateAction() am.FormAction {
	return am.NewFormAction(UserPath(), am.POST)
}

// userUpdateAction returns a new form action for user update.
func userUpdateAction(model am.Slugable) am.FormAction {
	return am.NewFormAction(UserPathSlug(model), am.PUT)
}

// userDeleteAction returns a new form action for user deletion.
func userDeleteAction(model am.Slugable) am.FormAction {
	return am.NewFormAction(UserPathSlug(model), am.DELETE)
}

// userSignUpAction returns a new form action for user sign up.
func userSignUpAction() am.FormAction {
	return am.NewFormAction(AuthPathSignUp(), am.POST)
}

// userSignInAction returns a new form action for user sign in.
func userSignInAction() am.FormAction {
	return am.NewFormAction(AuthPathSignIn(), am.POST)
}
