package auth

import (
	"github.com/adriancrafter/todoapp/internal/am"
)

type (
	UsersVM struct {
		List []UserVM `json:"users"`
	}

	UserVM struct {
		TenantID          string `json:"tenantID" schema:"tenant-id"`
		Slug              string `json:"slug" schema:"slug"`
		Username          string `json:"username" schema:"username"`
		Password          string `json:"password" schema:"password"`
		Email             string `json:"email" schema:"email"`
		EmailConfirmation string `json:"emailConfirmation" schema:"email-confirmation"`
		IsNew             bool   `json:"-" schema:"-"`
	}

	SigninVM struct {
		TenantID string
		Username string `json:"username" schema:"username"`
		Password string `json:"password" schema:"password"`
		IP       string
		GeoData  am.GeoPoint
	}
)

func NewSigninVM(tid string, ip string, gd am.GeoPoint) SigninVM {
	return SigninVM{
		TenantID: tid,
		IP:       ip,
		GeoData:  gd,
	}
}
