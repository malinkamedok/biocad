package v1

import (
	"biocad/internal/usecase"
	"github.com/go-chi/chi/v5"
)

func NewRouter(handler *chi.Mux, u usecase.UserContract) {
	handler.Route("/data", func(r chi.Router) {
		NewUserRoutes(r, u)
	})
}
