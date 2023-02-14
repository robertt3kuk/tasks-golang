package v1

import (
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/robertt3kuk/tasks-golang/server/docs"
	"github.com/robertt3kuk/tasks-golang/server/internal/service"
	"github.com/robertt3kuk/tasks-golang/server/pkg/logger"
)

// @title Task Golang
// @version 0.0.1
// @description This is a sample server to create and get user by email

// @contact.name API Support
// @contact.url https://t.me/biqontie
// @contact.email awesome.abaildaev@yandex.kz

// @license.name GPL-3
// @license.url https://www.gnu.org/licenses/gpl-3.0.en.html

// @host localhost:8000
// @BasePath /
func NewRouter(handler *chi.Mux, l logger.Interface, t *service.Service) {

	h := handler
	h.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8000/swagger/doc.json"),
	))
	{
		newUserRoutes(h, t.User, l)
	}

}
