package auth

import (
	"github.com/adriancrafter/todoapp/internal/am"
)

func (c *MainWebController) routes() {
	c.Router().HandleFunc("/users", c.UserIndex).Methods("GET")
	c.Router().HandleFunc("/users/{tableID}", c.UserShow).Methods("GET")
	c.Router().HandleFunc("/users", c.UserCreate).Methods("POST")
	c.Router().HandleFunc("/users/{tableID}", c.UserUpdate).Methods("PUT")
	c.Router().HandleFunc("/users/{tableID}", c.UserPreDelete).Methods("DELETE")
	c.Router().HandleFunc("/users/{tableID}", c.UserSoftDelete).Methods("PATCH")
	c.Router().HandleFunc("/users/{tableID}/purge", c.UserPurge).Methods("DELETE")
	c.Router().HandleFunc("/users/{tableID}/force-delete", c.UserDelete).Methods("DELETE")
	c.Router().HandleFunc("/init-signin", c.UserInitSignin).Methods("GET")
	c.Router().HandleFunc("/signin", c.UserSignin).Methods("POST")
	c.Router().HandleFunc("/init-signup", c.UserInitSignup).Methods("GET")
	c.Router().HandleFunc("/signup", c.UserSignup).Methods("POST")
}

// AuthRoot - Service root path.
var AuthRoot = "auth"

// UserRoot - User resource root path.
var UserRoot = "users"

func AuthPath() string {
	return am.WebPath.ResPath(AuthRoot)
}

// AuthPathSignUp returns the path to the signup page.
func AuthPathSignUp() string {
	return am.WebPath.ResPath(AuthRoot) + "/signup"
}

// AuthPathSignIn returns the path to the signin page.
func AuthPathSignIn() string {
	return am.WebPath.ResPath(AuthRoot) + "/signin"
}

// UserPath returns the path to the user resource.
func UserPath() string {
	return am.WebPath.ResPath(UserRoot)
}

// UserPathEdit returns the path to the user resource edit page.
func UserPathEdit(res am.Slugable) string {
	// TODO: Analize if in a multi-tenant setup this could be
	// a problem.
	return am.WebPath.ResPathEdit(UserRoot, res)
	//return fmt.Sprintf("/%s/%s/edit", UserRoot, res.U)
}

// UserPathNew returns the path to the user resource new page.
func UserPathNew() string {
	return am.WebPath.ResPathNew(UserRoot)
}

// UserPathInitDelete returns the path to the user resource init delete page.
func UserPathInitDelete(res am.Slugable) string {
	return am.WebPath.ResPathInitDelete(UserRoot, res)
}

// UserPathSlug returns the user resource slug path.
func UserPathSlug(res am.Slugable) string {
	return am.WebPath.ResPathSlug(UserRoot, res)
}
