package auth

import (
	"github.com/google/uuid"

	"github.com/adriancrafter/todoapp/internal/am"
)

type AccountRole struct {
	am.ID
	Name      string    `json:"name"`
	AccountID uuid.UUID `json:"account_id"`
	RoleID    uuid.UUID `json:"role_id"`
	IsActive  bool      `json:"is_active"`
	IsDeleted bool      `json:"is_deleted"`
	am.Audit
}

// NewAccountRole creates a new AccountRole instance with tenantID, name, accountID, and roleID.
func NewAccountRole(tenantID uuid.UUID, accountID uuid.UUID, roleID uuid.UUID) *AccountRole {
	return &AccountRole{
		ID:        am.NewID(tenantID),
		AccountID: accountID,
		RoleID:    roleID,
	}
}
