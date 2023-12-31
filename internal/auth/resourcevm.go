package auth

type (
	ResourcesVM struct {
		List []PermissionVM `json:"ResourcesVM"`
	}

	ResourceVM struct {
		TenantID    string `json:"tenantID" schema:"tenant-id"`
		Name        string `json:"name" schema:"name"`
		Description string `json:"description" schema:"description"`
		IsActive    bool   `json:"isActive" schema:"is-active"`
	}
)
