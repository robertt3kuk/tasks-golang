package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/robertt3kuk/tasks-golang/server/internal/entity"
	"github.com/robertt3kuk/tasks-golang/server/internal/service"
	"github.com/robertt3kuk/tasks-golang/server/pkg/logger"
)

type userRoutes struct {
	t service.User
	l logger.Interface
}

func newUserRoutes(handler *chi.Mux, t service.User, l logger.Interface) {
	r := &userRoutes{t, l}
	h := handler
	h.Post("/create-user", r.CreateUser)
	h.Get("/get-user/{email}", r.GetUser)
}

// @Summary CreateUser
// @Description create account
// @ID create-account
// @Accept  json
// @Produce  json
// @Param input body RequestUser true "account info"
// @Success 201 {string} string userCreated
// @Failure 400,404 {string} string
// @Failure 500 {string} string
// @Failure default {string} string
// @Router /create-user [post]
func (u *userRoutes) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req RequestUser
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		u.l.Error(err, "http - v1 - CreateUser")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("bad request body"))
		return
	}
	_, ok := u.t.Upload(r.Context(), entity.User{Email: req.Email, Password: req.Password})
	if ok.NotOK() {
		u.l.Error(ok.Err(), "http -v1 - CreateUser")
		w.WriteHeader(ok.Code())
		w.Write([]byte(ok.Msg()))
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("user %s created", req.Email)))
}

// @Summary Get User By ID
// @Description get user by id
// @ID get-user-by-id
// @Produce  json
// @Param email path string true "Email address of the user"  swaggoType
// @Success 200 {object} ResultUser
// @Failure 400,404 {string} string
// @Failure 500 {string} string
// @Failure default {string} string
// @Router /get-user/{email} [get]
func (u *userRoutes) GetUser(w http.ResponseWriter, r *http.Request) {
	email := chi.URLParam(r, "email")
	correct(&email)
	user, ok := u.t.GetByEmail(r.Context(), email)
	if ok.NotOK() {
		u.l.Error(ok.Err(), "http -v1 - GetUserByEmail")
		w.WriteHeader(ok.Code())
		w.Write([]byte(ok.Msg()))
		return
	}
	result, _ := json.Marshal(ResultUser{
		Email:    user.Email,
		Salt:     user.Salt,
		Password: user.Password,
	})
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func correct(email *string) {
	if strings.Contains(*email, "%40") {
		*email = strings.ReplaceAll(*email, "%40", "@")
	}
}
