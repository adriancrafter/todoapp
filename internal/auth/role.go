package auth

import (
	"github.com/google/uuid"

	"github.com/adriancrafter/todoapp/internal/am"
)

const (
	roleModel = "role"
)

type (
	// Role model
	Role struct {
		am.ID
		Name        string
		Description string
		IsActive    bool
		am.Audit
	}
)

// NewRole creates a new Role instance with tenantID, name, and description.
func NewRole(tenantID uuid.UUID, name, description string) *Role {
	return &Role{
		ID:          am.NewID(tenantID),
		Name:        name,
		Description: description,
	}
}
