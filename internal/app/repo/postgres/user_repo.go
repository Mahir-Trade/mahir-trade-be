package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"mahir-trade-be/internal/app/models"
	"mahir-trade-be/internal/app/repo/postgres/queries"

	"go.uber.org/dig"
)

type (
	UserRepo interface {
		CreateUser(ctx context.Context, req models.User) (id int, err error)
	}

	UserRepoImpl struct {
		dig.In

		*sql.DB
	}
)

func NewUserRepo(impl UserRepoImpl) UserRepo {
	return &impl
}

func (u *UserRepoImpl) CreateUser(ctx context.Context, req models.User) (id int, err error) {
	_, err = u.QueryContext(ctx, queries.QueryCreateUser, req.Email, req.Fullname, req.PhoneNumber, req.Username, req.Password)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("error while CreateUser err: %v", err.Error()))
		return id, err
	}

	return id, nil
}