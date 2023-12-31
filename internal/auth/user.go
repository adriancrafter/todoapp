package auth

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/adriancrafter/todoapp/internal/am"
)

const (
	userModel = "user"
)

type (
	// User model
	User struct {
		am.ID
		Username          string
		Password          string
		PasswordDigest    string
		Email             string
		EmailConfirmation string
		LastIP            string
		ConfirmationToken string
		IsConfirmed       bool
		LastGeoLocation   am.GeoPoint
		Since             time.Time
		Until             time.Time
		IsActive          bool
		am.Audit
	}

	UserAuth struct {
		User
		PermissionTags sql.NullString
	}
)

func NewUser(username, password string) *User {
	return &User{
		Username: username,
		Password: password,
	}
}

func (user *User) UpdatePasswordDigest() (digest string, err error) {
	if user.Password == "" {
		return user.PasswordDigest, nil
	}

	hpass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return user.PasswordDigest, err
	}
	user.PasswordDigest = string(hpass)
	return user.PasswordDigest, nil
}

func (user *User) SetCreateValues() (err error) {
	// Set create values only if not set
	if user.ID.Val() == uuid.Nil || user.ID.Slug == "" {
		pfx := user.Username
		user.ID.SetCreateValues(pfx)
		user.Audit.SetCreateValues()
		_, err = user.UpdatePasswordDigest()
	}
	return err
}

func (user *User) SetUpdateValues() (err error) {
	user.Audit.SetUpdateValues()
	_, err = user.UpdatePasswordDigest()
	return err
}

func (user *User) GenConfirmationToken() {
	user.ConfirmationToken = uuid.New().String()
	user.IsConfirmed = false
}

func (user *User) GenAutoConfirmationToken() {
	user.ConfirmationToken = uuid.New().String()
	user.IsConfirmed = true
}

func (user *User) Match(tc *User) bool {
	r := user.ID.Match(tc.ID) &&
		user.Username == tc.Username &&
		user.PasswordDigest == tc.PasswordDigest &&
		user.Email == tc.Email &&
		user.IsConfirmed == tc.IsConfirmed &&
		user.LastGeoLocation == tc.LastGeoLocation &&
		user.Since == tc.Since &&
		user.Until == tc.Until
	return r
}
