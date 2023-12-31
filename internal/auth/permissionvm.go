package auth

type (
	PermissionsVM struct {
		List []PermissionVM `json:"permissions"`
	}

	PermissionVM struct {
		TenantID    string `json:"tenantID" schema:"tenant-id"`
		Name        string `json:"name" schema:"name"`
		Description string `json:"description" schema:"description"`
		IsActive    bool   `json:"isActive" schema:"is-active"`
	}
)
