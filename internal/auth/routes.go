package auth

// There are convenience helpers to generate app standard paths
// but nothing prevents you from using your custom strings to define them.
// c.[Verb]("/path/to/resource"), HandlerFunc)
func (c *WebController) routes() {
	ur := c.userRoute
	ar := c.authroute

	c.Get(ur.Index(), c.UserIndex)
	c.Get(ur.Slug(), c.UserShow)
	c.Post(ur.Index(), c.UserCreate)
	c.Put(ur.Slug(), c.UserUpdate)
	c.Get(ur.Slug(), c.UserPreDelete)
	c.Get(ur.Slug(), c.UserSoftDelete)
	c.Delete(ur.ForceDelete(), c.UserDelete)
	c.Delete(ur.Purge(), c.UserPurge)
	c.Get(ar.Signup(), c.UserInitSignup)
	c.Post(ar.Signup(), c.UserSignup)
	c.Get(ar.Signin(), c.UserInitSignin)
	c.Post(ar.Signin(), c.UserSignin)
}
