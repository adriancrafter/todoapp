package auth

import (
	"github.com/google/uuid"

	"github.com/adriancrafter/todoapp/internal/am"
)

type Resource struct {
	am.ID
	Name        string `json:"name"`
	Description string `json:"description"`
	Tag         string `json:"tag"`
	Path        string `json:"path"`
	IsActive    bool   `json:"is_active"`
	am.Audit
}

// NewResource creates a new Resource instance with tenantID and name.
func NewResource(tenantID uuid.UUID, name, description string) *Resource {
	return &Resource{
		ID:          am.NewID(tenantID),
		Name:        name,
		Description: description,
	}
}
