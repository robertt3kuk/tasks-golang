package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/robertt3kuk/tasks-golang/server/internal/entity"
)

type UserMongo struct {
	db *mongo.Collection
}
type UserM struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Email    string             `bson:"email"`
	Salt     string             `bson:"salt"`
	Password string             `bson:"password"`
}

func NewUser(db *mongo.Database, collection string) *UserMongo {
	return &UserMongo{
		db: db.Collection(collection),
	}
}

func (u *UserMongo) Add(ctx context.Context, user entity.User) (entity.User, error) {
	model := toMongoUser(user)
	_, err := u.db.InsertOne(ctx, model)
	if err != nil {
		return user, err
	}
	return user, err

}
func (u *UserMongo) Exists(ctx context.Context, email string) (bool, error) {
	model := new(UserM)
	err := u.db.FindOne(ctx, bson.M{
		"email": email,
	}).Decode(model)
	if err != nil {
		return false, err
	}
	return true, nil

}
func (u *UserMongo) GetByEmail(ctx context.Context, email string) (entity.User, error) {
	var model UserM
	err := u.db.FindOne(ctx, bson.M{
		"email": email,
	}).Decode(&model)
	result := toEntityUser(model)
	if err != nil {
		return result, err
	}
	return result, nil
}

func toMongoUser(u entity.User) UserM {
	return UserM{
		Email:    u.Email,
		Salt:     u.Salt,
		Password: u.Password,
	}

}
func toEntityUser(u UserM) entity.User {
	return entity.User{
		ID:       u.ID.Hex(),
		Email:    u.Email,
		Salt:     u.Salt,
		Password: u.Password,
	}
}
