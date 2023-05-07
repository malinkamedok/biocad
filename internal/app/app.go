package app

import (
	"biocad/internal/configs"
	v1 "biocad/internal/controller/http/v1"
	"biocad/internal/usecase"
	"biocad/internal/usecase/repo"
	"biocad/pkg/httpserver"
	"biocad/pkg/postgres"
	"github.com/go-chi/chi/v5"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func Run(cfg *configs.Config) {

	pg, err := postgres.New(cfg)
	if err != nil {
		log.Fatal("Cannot connect to Postgres")
	}

	u := usecase.NewUserUseCase(repo.NewUserRepo(pg))

	handler := chi.NewRouter()

	v1.NewRouter(handler, u)

	serv := httpserver.New(handler, httpserver.Port(cfg.AppPort))
	interruption := make(chan os.Signal, 1)
	signal.Notify(interruption, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interruption:
		log.Printf("signal: " + s.String())
	case err = <-serv.Notify():
		log.Printf("Notify from http server")
	}

	err = serv.Shutdown()
	if err != nil {
		log.Printf("Http server shutdown")
	}
}
