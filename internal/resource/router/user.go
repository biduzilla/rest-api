package router

import (
	"database/sql"
	"rest-api/internal/resource/handler"
	"rest-api/internal/resource/repository"
	e "rest-api/utils/errors"
	"rest-api/utils/validator"

	"github.com/go-chi/chi"
)

func UserRouter(r chi.Router, db *sql.DB, v *validator.Validator, errResp *e.ErrorResponse) {
	r.Route("/users", func(r chi.Router) {
		repository := repository.NewRepositories(db)
		userAPI := handler.NewUserHandler(repository.UserRepository, v, errResp)

		r.Post("/", userAPI.RegisterUserHandler)
		r.Put("/activate", userAPI.ActivateUserHandler)
	})
}
