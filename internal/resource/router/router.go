package router

import (
	"database/sql"
	e "rest-api/utils/errors"
	"rest-api/utils/validator"

	"github.com/go-chi/chi"
)

func New(db *sql.DB, v *validator.Validator, errResp *e.ErrorResponse) *chi.Mux {
	r := chi.NewRouter()

	r.Route("/v1", func(r chi.Router) {
		UserRouter(r, db, v, errResp)
	})
	return r
}
