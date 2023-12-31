package auth

import (
	"github.com/google/uuid"

	"github.com/adriancrafter/todoapp/internal/am"
)

type Permission struct {
	am.ID
	Name        string `json:"name"`
	Description string `json:"description"`
	Tag         string `json:"tag"`
	Path        string `json:"path"`
	IsActive    bool   `json:"is_active"`
	am.Audit
}

// NewPermission creates a new Permission instance with tenantID and name.
func NewPermission(tenantID uuid.UUID, name, description string) *Permission {
	return &Permission{
		ID:          am.NewID(tenantID),
		Name:        name,
		Description: description,
	}
}
