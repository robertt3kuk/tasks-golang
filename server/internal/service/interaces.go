package service

import (
	"context"

	"github.com/robertt3kuk/tasks-golang/server/internal/entity"
	"github.com/robertt3kuk/tasks-golang/server/internal/service/repository"
	"github.com/robertt3kuk/tasks-golang/server/pkg/errs"
)

type Service struct {
	User
}

type (
	User interface {
		Upload(ctx context.Context, user entity.User) (entity.User, errs.Errs)
		GetByEmail(ctx context.Context, email string) (entity.User, errs.Errs)
		exists(ctx context.Context, email string) (bool, error)
		emailValidation(email string) error
		getSalt() (string, error)
	}
)

func New(repo *repository.Repo) *Service {
	return &Service{User: NewUserService(repo.User)}
}
