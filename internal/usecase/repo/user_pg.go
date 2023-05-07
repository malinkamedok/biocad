package repo

import (
	"biocad/internal/usecase"
	"biocad/pkg/postgres"
)

type UserRepo struct {
	*postgres.Postgres
}

var _ usecase.UserRp = (*UserRepo)(nil)

func NewUserRepo(pg *postgres.Postgres) *UserRepo {
	return &UserRepo{pg}
}
