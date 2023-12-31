package auth

type (
	AffiliationsVM struct {
		List []AffiliationVM `json:"affiliations"`
	}

	AffiliationVM struct {
		TenantID string `json:"tenantID" schema:"tenant-id"`
		UserID   string `json:"userID" schema:"user-id"`
		Type     string `json:"type" schema:"type"`
	}
)
