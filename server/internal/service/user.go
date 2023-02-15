package service

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/robertt3kuk/tasks-golang/server/internal/entity"
	"github.com/robertt3kuk/tasks-golang/server/internal/service/repository"
	"github.com/robertt3kuk/tasks-golang/server/pkg/errs"
)

type UserService struct {
	repo repository.User
}

func NewUserService(repo repository.User) *UserService {
	return &UserService{repo: repo}

}

func (u UserService) Upload(ctx context.Context, user entity.User) (entity.User, errs.Errs) {
	salt, err := u.getSalt()
	if err != nil {
		return user, errs.New(http.StatusInternalServerError, err, "getting salt error")
	}
	user.Salt = salt

	err = u.emailValidation(user.Email)
	if err != nil {
		return user, errs.New(http.StatusBadRequest, err, err.Error())
	}
	if user.Password == "" {
		return user, errs.New(http.StatusBadRequest, fmt.Errorf("empty password"), "empty password")
	}

	exist, err := u.exists(ctx, user.Email)
	if err != nil {
		return user, errs.New(http.StatusInternalServerError, err, "error while checking email")
	}
	if exist {
		return user, errs.New(http.StatusBadRequest, fmt.Errorf("user exists already"), "user already registered")
	}
	hash := md5.Sum([]byte(salt + user.Password))
	user.Password = hex.EncodeToString(hash[:])

	_, err = u.repo.Add(ctx, user)
	if err != nil {
		return user, errs.New(http.StatusInternalServerError, err, "error while saving account")

	}

	return user, errs.New(http.StatusCreated, nil, "user saved")
}

func (u UserService) GetByEmail(ctx context.Context, email string) (entity.User, errs.Errs) {
	var user entity.User

	err := u.emailValidation(email)
	if err != nil {
		return user, errs.New(http.StatusBadRequest, err, err.Error())
	}

	exist, err := u.exists(ctx, email)
	if err != nil {
		return user, errs.New(http.StatusInternalServerError, err, "error while checking email")
	}
	if !exist {
		return user, errs.New(http.StatusNotFound, fmt.Errorf("user doesn't exist"), "user doesn't exist with this email")
	}

	user, err = u.repo.GetByEmail(ctx, email)
	if err != nil {
		return user, errs.New(http.StatusInternalServerError, err, "errow while getting email")
	}

	return user, errs.New(http.StatusOK, nil, "user with email")
}
func (u UserService) emailValidation(email string) error {
	pattern := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	matched, err := regexp.MatchString(pattern, email)
	if err != nil {
		return err
	}
	if !matched {
		return fmt.Errorf("not valid email")
	}
	return nil
}
func (u *UserService) exists(ctx context.Context, email string) (bool, error) {
	exist, err := u.repo.Exists(ctx, email)
	if err == mongo.ErrNoDocuments {
		exist = false
		return exist, nil
	}
	if err != nil {
		return exist, err
	}
	return exist, nil
}

func (u *UserService) getSalt() (string, error) {
	type Salt struct {
		S string `json:"s"`
	}

	url := "http://172.16.238.12:3000/generate-salt"

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var salt Salt
	err = json.Unmarshal(body, &salt)
	if err != nil {
		return "", err
	}
	return salt.S, nil
}
