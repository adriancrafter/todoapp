package auth

import (
	"time"

	"github.com/google/uuid" // Assuming UUID is handled by this package

	"github.com/adriancrafter/todoapp/internal/am"
)

type Profile struct {
	am.ID
	OwnerID        uuid.UUID `json:"owner_id"`
	AccountType    string    `json:"account_type"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	Location       string    `json:"location"`
	Bio            string    `json:"bio"`
	Moto           string    `json:"moto"`
	Website        string    `json:"website"`
	AniversaryDate time.Time `json:"aniversary_date"`
	AvatarSmall    []byte    `json:"avatar_small"`
	HeaderSmall    string    `json:"header_small"`
	AvatarPath     string    `json:"avatar_path"`
	HeaderPath     string    `json:"header_path"`
	ExtendedData   []byte    `json:"extended_data"`
	Geolocation    string    `json:"geolocation"`
	IsActive       bool      `json:"is_active"`
	am.Audit                 // Embedded struct for Audit fields
}

// NewProfile creates a new Profile instance with tenantID and name.
func NewProfile(tenantID uuid.UUID, name string) *Profile {
	return &Profile{
		ID:   am.NewID(tenantID),
		Name: name,
	}
}
