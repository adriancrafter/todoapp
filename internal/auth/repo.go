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
	SignInQry = "SignIn"
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

// SignIn handles user sign-in in MainRepo
func (r *MainRepo) SignIn(ctx context.Context, si Signin) (ua UserAuth, err error) {
	st, err := r.Query(AuthModel, SignInQry)
	if err != nil {
		return ua, errors.Wrap(err, "query now found")
	}

	//r.Log().Debugf("SQL:\n%s", st)

	uada := &UserAuthDA{}
	err = r.DB().Get(uada, st, si.Username, si.Email)
	if err != nil {
		return ua, err
	}

	// Validate password
	// NOTE: This should be done in service layer if does not create too much burden.
	err = bcrypt.CompareHashAndPassword([]byte(uada.PasswordDigest.String), []byte(si.Password))
	if err != nil {
		return ua, errors.NewError("invalid password")
	}

	ua = UserAuthDAToModel(*uada)

	return ua, nil
}
