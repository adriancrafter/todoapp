package auth

import (
	"context"

	"golang.org/x/crypto/bcrypt"

	"github.com/adriancrafter/todoapp/internal/am"
	"github.com/adriancrafter/todoapp/internal/am/errors"
)

const (
	UserModel = "User"
	AuthModel = "Auth"
)

const (
	SignInQry       = "Signin"
	UpdateSigninQry = "UpdateSignin"
)

type (
	MainRepo struct {
		*am.SimpleRepo
	}
)

func NewRepo(db am.DB, qm *am.QueryManager, opts ...am.Option) *MainRepo {
	name := "auth-repo"
	return &MainRepo{
		SimpleRepo: am.NewRepo(name, db, qm, opts...),
	}
}

func (r *MainRepo) GetUser(ctx context.Context, userID string) (user User, err error) {
	return user, errors.NewError("not implemented")
}

func (r *MainRepo) SignIn(ctx context.Context, si Signin) (ua UserAuth, err error) {
	sda := SigninToDA(si)
	uada := UserAuthDA{}

	st, err := r.Query(AuthModel, SignInQry)
	if err != nil {
		return ua, errors.Wrapf(err, "query not found: %s", SignInQry)
	}

	//r.Log().Debugf("query: %s", st)

	tx, err := r.DB().BeginTxx(ctx, nil)
	if err != nil {
		return ua, err
	}

	err = tx.Get(&uada, st, sda.TenantID, sda.Username, sda.Email)
	if err != nil {
		return ua, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(uada.PasswordDigest.String), []byte(si.Password))
	if err != nil {
		return ua, errors.Wrapf(err, "invalid password")
	}

	lng, lat := sda.GeoData.Lng, sda.GeoData.Lat

	st2, err := r.Query(AuthModel, UpdateSigninQry)
	if err != nil {
		return ua, errors.Wrapf(err, "query not found: %s", SignInQry)
	}

	//r.Log().Debugf("query: %s", st2)

	_, err = tx.ExecContext(ctx, st2, sda.IP, lng, lat, uada.TenantID, uada.Slug)

	if err != nil {
		err2 := tx.Rollback()
		if err2 != nil {
			return ua, errors.Wrapf(err2, "error rolling back transaction caused by err: %s", err.Error())
		}
		return ua, err
	}

	// NOTE: A second query is needed to get the updated data.
	// Probably not really needed but not a big deal for now.
	// A `...RETURNING *` clause in previous could work.
	err = tx.Get(&uada, st, sda.TenantID, sda.Username, sda.Email)
	if err != nil {
		return ua, err
	}

	if err = tx.Commit(); err != nil {
		return ua, errors.Wrapf(err, "error committing transaction")
	}

	ua = UserAuthDAToModel(uada)
	return ua, nil
}
