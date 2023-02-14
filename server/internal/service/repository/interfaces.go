package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/robertt3kuk/tasks-golang/server/internal/entity"
)

type Repo struct {
	User
}

type User interface {
	Add(ctx context.Context, user entity.User) (entity.User, error)
	Exists(ctx context.Context, email string) (bool, error)
	GetByEmail(ctx context.Context, email string) (entity.User, error)
}

var (
	UserCollection = "users"
)

func New(db *mongo.Database) *Repo {
	return &Repo{User: NewUser(db, UserCollection)}
}
