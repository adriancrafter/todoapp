package auth

import (
	"github.com/adriancrafter/todoapp/internal/am"
)

const (
	tenantModel = "tenant"
)

type (
	// Tenant model
	Tenant struct {
		am.ID
		Name        string
		FantasyName string
		ProfileID   string
		IsActive    bool
		OwnerID     string
		am.Audit
	}
)
