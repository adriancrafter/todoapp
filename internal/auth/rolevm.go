package auth

type (
	RolesVM struct {
		List []RoleVM `json:"roles"`
	}

	RoleVM struct {
		TenantID    string `json:"tenantID" schema:"tenant-id"`
		Name        string `json:"name" schema:"name"`
		Description string `json:"description" schema:"description"`
		IsActive    bool   `json:"isActive" schema:"is-active"`
	}
)
