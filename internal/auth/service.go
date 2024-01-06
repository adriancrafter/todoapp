package auth

import (
	"context"

	"github.com/adriancrafter/todoapp/internal/am"
)

type (
	MainService struct {
		*am.SimpleCore
		repo Repo
	}
)

func NewService(r Repo, opts ...am.Option) *MainService {
	return &MainService{
		SimpleCore: am.NewCore("todo-service", opts...),
		repo:       r,
	}
}

func (as *MainService) SignInUser(ctx context.Context, sivm SigninVM) (user UserVM, err error) {
	signin := ToSigninModel(sivm)

	userAuth, err := as.repo.SignIn(ctx, signin)
	if err != nil {
		as.Log().Errorf("error signing in: %s", err.Error())
		return UserVM{}, err
	}

	userVM := ToUserVM(userAuth.User)

	return userVM, nil
}
