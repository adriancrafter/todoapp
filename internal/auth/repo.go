package auth

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/adriancrafter/todoapp/internal/am"
	"github.com/adriancrafter/todoapp/internal/am/errors"
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
	//err = r.db.GetErr(ctx, &user, "SELECT * FROM users WHERE id = $1", userID)
	//if err != nil {
	//	return user, err
	//}
	return user, errors.NewError("not implemented")
}

func (r *MainRepo) SignIn(username, password string) (ua *UserAuth, err error) {
	st := `SELECT users.*, array_to_string(array_agg(DISTINCT permissions.tag), ',') as permission_tags FROM users
						INNER JOIN accounts ON accounts.owner_id = users.id
						INNER JOIN account_roles ON account_roles.account_id = accounts.id
						INNER JOIN roles ON roles.id = account_roles.role_id
						INNER JOIN role_permissions ON role_permissions.role_id = roles.id
						INNER JOIN permissions ON permissions.id = role_permissions.permission_id
						WHERE users.username = 'superadmin' OR users.email = 'superadmin'
						AND (users.is_deleted IS NULL OR NOT users.is_deleted)
						AND (accounts.is_deleted IS NULL OR NOT accounts.is_deleted)
						AND (account_roles.is_deleted IS NULL OR NOT account_roles.is_deleted)
						AND (roles.is_deleted IS NULL OR NOT roles.is_deleted)
						AND (role_permissions.is_deleted IS NULL OR NOT role_permissions.is_deleted)
						AND (permissions.is_deleted IS NULL OR NOT permissions.is_deleted)
						AND (users.is_active IS NULL OR users.is_active)
						AND (accounts.is_active IS NULL OR accounts.is_active)
						AND (account_roles.is_active IS NULL OR account_roles.is_active)
						AND (roles.is_active IS NULL OR roles.is_active)
						AND (role_permissions.is_active IS NULL OR role_permissions.is_active)
						AND (permissions.is_active IS NULL OR permissions.is_active)
						AND users.username = '%s' or users.email = '%s'
						GROUP BY users.ID;`

	st = fmt.Sprintf(st, username, username)

	//r.Log.Info(st)

	err = r.DB().Get(&ua, st)

	// Validate password
	err = bcrypt.CompareHashAndPassword([]byte(ua.PasswordDigest), []byte(password))
	if err != nil {
		r.Log().Errorf("error validating password: %s", err.Error())
		return ua, err
	}

	return ua, nil
}
