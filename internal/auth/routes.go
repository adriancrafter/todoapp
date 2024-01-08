package auth

func (c *WebController) routes() {
	c.Get(c.up.ResRoute(), c.UserIndex)
	c.Get(c.up.ResSlugRoute(), c.UserShow)
	c.Post(c.up.ResRoute(), c.UserCreate)
	c.Put(c.up.ResSlugRoute(), c.UserUpdate)
	c.Get(c.up.ResSlugRoute(), c.UserPreDelete)
	c.Get(c.up.ResSlugRoute(), c.UserSoftDelete)
	c.Delete(c.up.ResForceDeleteRoute(), c.UserDelete)
	c.Delete(c.up.ResPurgeRoute(), c.UserPurge)
	c.Get(c.ap.SignupRoute(), c.UserInitSignup)
	c.Post(c.ap.SignupRoute(), c.UserSignup)
	c.Get(c.ap.SigninPath(), c.UserInitSignin)
	c.Post(c.ap.SigninPath(), c.UserSignin)
}
