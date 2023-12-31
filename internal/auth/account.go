package auth

import (
	"github.com/google/uuid"

	"github.com/adriancrafter/todoapp/internal/am"
)

const (
	accountModel = "account"
)

type Account struct {
	am.ID                 // Embedded struct for ID
	OwnerID     uuid.UUID `json:"owner_id"`
	ParentID    uuid.UUID `json:"parent_id"`
	AccountType string    `json:"account_type"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	GivenName   string    `json:"given_name"`
	MiddleNames string    `json:"middle_names"`
	FamilyName  string    `json:"family_names"`
	Locale      string    `json:"locale"`
	BaseTz      string    `json:"base_tz"`
	CurrentTz   string    `json:"current_tz"`
	IsActive    bool      `json:"is_active"`
	am.Audit              // Embedded struct for Audit fields
}

// NewAccount creates a new Account instance.
func NewAccount(tenantID uuid.UUID, ownerID, parentID uuid.UUID) *Account {
	return &Account{
		ID:       am.NewID(tenantID),
		OwnerID:  ownerID,
		ParentID: parentID,
	}
}
