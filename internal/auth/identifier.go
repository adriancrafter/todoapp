package auth

const (
	// There is one to one mapping between resource name and resource controller
	// names and paths. A (future) code generator will take advantage of this
	// but nothings prevents you from using your custom strings to define them
	// in the places where they are used.
	userRes = "user"
	authRes = "auth"
)
