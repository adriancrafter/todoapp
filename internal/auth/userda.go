package auth

import (
	"database/sql"

	"github.com/google/uuid"

	"github.com/adriancrafter/todoapp/internal/am/db/pg"
)

type UserDA struct {
	ID                uuid.UUID       `db:"id"`
	TenantID          uuid.UUID       `db:"tenant_id"`
	Slug              sql.NullString  `db:"slug"`
	Username          sql.NullString  `db:"username"`
	Password          sql.NullString  `db:"password"`
	PasswordDigest    sql.NullString  `db:"password_digest"`
	Email             sql.NullString  `db:"email"`
	EmailConfirmation sql.NullString  `db:"email_confirmation"`
	LastIP            sql.NullString  `db:"last_ip"`
	ConfirmationToken sql.NullString  `db:"confirmation_token"`
	IsConfirmed       sql.NullBool    `db:"is_confirmed"`
	LastGeoLocation   pg.NullGeoPoint `db:"last_geolocation"`
	Since             sql.NullTime    `db:"since"`
	Until             sql.NullTime    `db:"until"`
	IsActive          sql.NullBool    `db:"is_active"`
	CreatedByID       uuid.UUID       `db:"created_by_id"`
	UpdatedByID       uuid.UUID       `db:"updated_by_id"`
	DeletedByID       uuid.UUID       `db:"deleted_by_id"`
	CreatedAt         sql.NullTime    `db:"created_at"`
	UpdatedAt         sql.NullTime    `db:"updated_at"`
	DeletedAt         sql.NullTime    `db:"deleted_at"`
}

type UserAuthDA struct {
	UserDA
	PermissionTags sql.NullString `db:"permission_tags"`
}

type SigninDA struct {
	Username sql.NullString `db:"username"`
	Password sql.NullString `db:"password"`
	Email    sql.NullString `db:"email"`
}
