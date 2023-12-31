package auth

import (
	"context"

	"github.com/adriancrafter/todoapp/internal/am"
)

type (
	MainService struct {
		*am.SimpleCore
		repo MainRepo
	}
)

func NewService(r MainRepo, opts ...am.Option) *MainService {
	return &MainService{
		SimpleCore: am.NewCore("todo-service", opts...),
		repo:       r,
	}
}

func (as *MainService) SignInUser(ctx context.Context, credentials UserVM) (user UserVM, err error) {
	return UserVM{}, nil
}
