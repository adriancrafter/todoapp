package auth

func (c *WebController) routes() {
	// TODO: Improve ergonomics of this.
	// This will be automatically generated but still could look better.
	c.Router().HandleFunc(c.up.ResRoute(), c.UserIndex).Methods("GET")
	c.Router().HandleFunc(c.up.ResSlugRoute(), c.UserShow).Methods("GET")
	c.Router().HandleFunc(c.up.ResRoute(), c.UserCreate).Methods("POST")
	c.Router().HandleFunc(c.up.ResSlugRoute(), c.UserUpdate).Methods("PUT")
	c.Router().HandleFunc(c.up.ResSlugRoute(), c.UserPreDelete).Methods("DELETE")
	c.Router().HandleFunc(c.up.ResSlugRoute(), c.UserSoftDelete).Methods("PATCH")
	c.Router().HandleFunc(c.up.ResForceDeleteRoute(), c.UserDelete).Methods("DELETE")
	c.Router().HandleFunc(c.up.ResPurgeRoute(), c.UserPurge).Methods("DELETE")
	c.Router().HandleFunc(c.ap.SignupRoute(), c.UserInitSignup).Methods("GET")
	c.Router().HandleFunc(c.ap.SignupRoute(), c.UserSignup).Methods("POST")
	c.Router().HandleFunc(c.ap.SigninPath(), c.UserInitSignin).Methods("GET")
	c.Router().HandleFunc(c.ap.SigninPath(), c.UserSignin).Methods("POST")
}
