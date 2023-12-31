package auth

import (
	"database/sql"

	"github.com/google/uuid"
)

type (
	UserDA struct {
		Val                uuid.UUID
		TenantID           uuid.UUID
		Slug               sql.NullString
		Username           sql.NullString
		Password           sql.NullString
		PasswordDigest     sql.NullString
		Email              sql.NullString
		EmailConfirmation  sql.NullString
		LastIP             sql.NullString
		ConfirmationToken  sql.NullString
		IsConfirmed        sql.NullBool
		LastGeoLocationLng sql.NullFloat64
		LastGeoLocationLat sql.NullFloat64
		LastGeoLocationAlt sql.NullFloat64
		Since              sql.NullTime
		Until              sql.NullTime
		IsActive           sql.NullBool
		CreatedByID        uuid.UUID
		UpdatedByID        uuid.UUID
		DeletedByID        uuid.UUID
		CreatedAt          sql.NullTime
		UpdatedAt          sql.NullTime
		DeletedAt          sql.NullTime
	}

	UserAuthDA struct {
		UserDA
		PermissionTags sql.NullString
	}
)
