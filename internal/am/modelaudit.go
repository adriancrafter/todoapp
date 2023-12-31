package am

import (
	"time"

	"github.com/google/uuid"
)

type Audit struct {
	CreatedByID uuid.UUID
	UpdatedByID uuid.UUID
	DeletedByID uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   time.Time
}

func (a *Audit) SetCreateValues(userID ...uuid.UUID) {
	if len(userID) > 0 {
		a.CreatedByID = userID[0]
	}

	now := time.Now()
	a.CreatedAt = now
}

func (a *Audit) SetUpdateValues(userID ...uuid.UUID) {
	if len(userID) > 0 {
		a.UpdatedByID = userID[0]
	}

	now := time.Now()
	a.UpdatedAt = now
}

func (a *Audit) SetDeleteValues(userID ...uuid.UUID) {
	if len(userID) > 0 {
		a.DeletedByID = userID[0]
	}

	now := time.Now()
	a.DeletedAt = now
}
