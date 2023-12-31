package auth

type (
	AccountsVM struct {
		List []AccountVM `json:"users"`
	}

	AccountVM struct {
		TenantID    string `json:"tenantID" schema:"tenant-id"`
		UserID      string `json:"userID" schema:"user-id"`
		Slug        string `json:"slug" schema:"slug"`
		Name        string `json:"username" schema:"username"`
		Description string `json:"description" schema:"description"`
	}
)
