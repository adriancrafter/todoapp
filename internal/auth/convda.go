package auth

import (
	"database/sql"

	"github.com/google/uuid"

	"github.com/adriancrafter/todoapp/internal/am"
	"github.com/adriancrafter/todoapp/internal/am/db"
)

func UserToDA(m User) UserDA {
	return UserDA{
		ID:                m.ID.Val(),
		TenantID:          m.ID.TenantID,
		Slug:              sql.NullString{String: m.ID.Slug, Valid: m.ID.Slug != ""},
		Username:          sql.NullString{String: m.Username, Valid: m.Username != ""},
		Password:          sql.NullString{String: m.Password, Valid: m.Password != ""},
		PasswordDigest:    sql.NullString{String: m.PasswordDigest, Valid: m.PasswordDigest != ""},
		Email:             sql.NullString{String: m.Email, Valid: m.Email != ""},
		EmailConfirmation: sql.NullString{String: m.EmailConfirmation, Valid: m.EmailConfirmation != ""},
		LastIP:            sql.NullString{String: m.LastIP, Valid: m.LastIP != ""},
		ConfirmationToken: sql.NullString{String: m.ConfirmationToken, Valid: m.ConfirmationToken != ""},
		IsConfirmed:       sql.NullBool{Bool: m.IsConfirmed, Valid: true},
		Since:             sql.NullTime{Time: m.Since, Valid: true},
		Until:             sql.NullTime{Time: m.Until, Valid: !m.Until.IsZero()},
		IsActive:          sql.NullBool{Bool: m.IsActive, Valid: true},
		CreatedByID:       m.Audit.CreatedByID,
		UpdatedByID:       m.Audit.UpdatedByID,
		DeletedByID:       m.Audit.DeletedByID,
		CreatedAt:         sql.NullTime{Time: m.Audit.CreatedAt, Valid: !m.Audit.CreatedAt.IsZero()},
		UpdatedAt:         sql.NullTime{Time: m.Audit.UpdatedAt, Valid: !m.Audit.UpdatedAt.IsZero()},
		DeletedAt:         sql.NullTime{Time: m.Audit.DeletedAt, Valid: !m.Audit.DeletedAt.IsZero()},
	}
}

func UserDAToModel(da UserDA) User {
	var lastGeoLocation am.GeoPoint
	if da.LastGeoLocation.Valid {
		lastGeoLocation = da.LastGeoLocation.GeoPoint
	}

	return User{
		ID:                am.NewID(da.TenantID, da.ID),
		Username:          da.Username.String,
		Password:          da.Password.String,
		PasswordDigest:    da.PasswordDigest.String,
		Email:             da.Email.String,
		EmailConfirmation: da.EmailConfirmation.String,
		LastIP:            da.LastIP.String,
		ConfirmationToken: da.ConfirmationToken.String,
		IsConfirmed:       da.IsConfirmed.Bool,
		LastGeoLocation:   lastGeoLocation,
		Since:             da.Since.Time,
		Until:             da.Until.Time,
		IsActive:          da.IsActive.Bool,
		Audit: am.Audit{
			CreatedByID: da.CreatedByID,
			UpdatedByID: da.UpdatedByID,
			DeletedByID: da.DeletedByID,
			CreatedAt:   da.CreatedAt.Time,
			UpdatedAt:   da.UpdatedAt.Time,
			DeletedAt:   da.DeletedAt.Time,
		},
	}
}

func UserAuthToDA(m UserAuth) UserAuthDA {
	userDA := UserToDA(m.User)
	return UserAuthDA{
		UserDA:         userDA,
		PermissionTags: m.PermissionTags,
	}
}

func SigninToDA(m Signin) SigninDA {
	return SigninDA{
		TenantID: db.NewNullString(m.TenantID),
		Slug:     db.NewNullString(m.Slug),
		Username: db.NewNullString(m.Username),
		Password: db.NewNullString(m.Password),
		Email:    db.NewNullString(m.Email),
		IP:       db.NewNullString(m.IP),
		GeoData:  db.NewNullGeoPoint(m.GeoData),
	}
}

func SigninDAToModel(da SigninDA) Signin {
	var geoData am.GeoPoint
	if da.GeoData.Valid {
		geoData = da.GeoData.GeoPoint
	}

	return Signin{
		TenantID: da.TenantID.String,
		Slug:     da.Slug.String,
		Username: da.Username.String,
		Password: da.Password.String,
		Email:    da.Email.String,
		IP:       da.IP.String,
		GeoData:  geoData,
	}
}

func UserAuthDAToModel(da UserAuthDA) UserAuth {
	user := UserDAToModel(da.UserDA)
	return UserAuth{
		User:           user,
		PermissionTags: da.PermissionTags,
	}
}

// StringToUUID converts a string to a UUID.
// If the conversion fails, it returns an error.
func StringToUUID(input string) (uuid.UUID, error) {
	return uuid.Parse(input)
}
