package auth

type (
	TenantsVM struct {
		List []TenantVM `json:"users"`
	}

	TenantVM struct {
		TenantID    string `json:"tenantID" schema:"tenant-id"`
		Slug        string `json:"slug" schema:"slug"`
		Name        string `json:"name" schema:"name"`
		FantasyName string `json:"fantasyName" schema:"fantasy-name"`
		ProfileID   string `json:"profileID" schema:"profile-id"`
		IsActive    bool   `json:"isActive" schema:"is-active"`
		OwnerID     string `json:"ownerID" schema:"owner-id"`
	}
)
