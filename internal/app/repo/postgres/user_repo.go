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
		FindUserByEmailOrUsername(ctx context.Context, req string) (user models.User, err error)
		GetUserByID(ctx context.Context, id int64) (user models.User, err error)
		GetUserByUUID(ctx context.Context, uuid string) (user models.User, err error)
		UpdateTypeUser(ctx context.Context, isActive bool, operator string, id int64) (typeUser bool, err error)
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

func (u *UserRepoImpl) FindUserByEmailOrUsername(ctx context.Context, req string) (user models.User, err error) {
	row := u.QueryRowContext(ctx, queries.QueryFindUserByEmail, req)
	err = row.Scan(&user.UserID, &user.UUID, &user.Fullname, &user.PhoneNumber, &user.Username, &user.Password)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("error while FindUserByEmailOrUsername err: %v", err.Error()))
		return user, err
	}

	return user, nil
}

func (u *UserRepoImpl) GetUserByID(ctx context.Context, id int64) (user models.User, err error) {
	row := u.QueryRowContext(ctx, queries.QueryGetUserByID, id)
	err = row.Scan(&user.UserID, &user.UUID, &user.Fullname, &user.PhoneNumber, &user.Username, &user.Password)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("error while GetUserByID err: %v", err.Error()))
		return user, err
	}

	return user, nil
}

func (u *UserRepoImpl) GetUserByUUID(ctx context.Context, uuid string) (user models.User, err error) {
	row := u.QueryRowContext(ctx, queries.QueryGetUserByUUID, uuid)
	err = row.Scan(&user.UserID, &user.UUID, &user.Fullname, &user.PhoneNumber, &user.Username, &user.Password)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("error while GetUserByUUID err: %v", err.Error()))
		return user, err
	}

	return user, nil
}

func (u *UserRepoImpl) UpdateTypeUser(ctx context.Context, isActive bool, operator string, id int64) (typeUser bool, err error) {
	_, err = u.QueryContext(ctx, queries.QueryUpdateTypeUser, isActive, operator, id)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("error while UpdateTypeUser err: %v", err.Error()))
		return typeUser, err
	}

	return true, nil
}
