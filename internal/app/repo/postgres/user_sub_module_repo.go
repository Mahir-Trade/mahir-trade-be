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
	UserSubModuleRepo interface {
		CreateUserSubModule(ctx context.Context, req models.UserSubModule) (id int64, err error)
		GetUserSubModuleBySubModuleIDAndUserID(ctx context.Context, userID, subModuleID int64) (resp models.UserSubModule, err error)
	}

	UserSubModuleRepoImpl struct {
		dig.In

		*sql.DB
	}
)

func NewUserSubModuleRepo(impl UserSubModuleRepoImpl) UserSubModuleRepo {
	return &impl
}

func (u *UserSubModuleRepoImpl) CreateUserSubModule(ctx context.Context, req models.UserSubModule) (id int64, err error) {
	err = u.QueryRowContext(ctx, queries.QueryCreateUserSubModule, req.UserID, req.SubModuleID, req.CreatedBy).Scan(&id)
	if err != nil {
		slog.ErrorContext(ctx, "[userSubModuleRepoImpl][CreateUserSubModule] error while row.Scan", "%v", err.Error())
		err = fmt.Errorf("something went wrong, we will fix it soon")
		return
	}

	return
}

func (u *UserSubModuleRepoImpl) GetUserSubModuleBySubModuleIDAndUserID(ctx context.Context, userID, subModuleID int64) (resp models.UserSubModule, err error) {
	row := u.QueryRowContext(ctx, queries.QueryGetUserSubModuleBySubModuleIDAndUserID, subModuleID, userID)
	err = row.Scan(
		&resp.ID,
		&resp.UUID,
		&resp.UserID,
		&resp.SubModuleID,
		&resp.CreatedBy,
		&resp.UpdatedBy,
		&resp.CreatedAt,
		&resp.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return resp, nil
		}

		slog.ErrorContext(ctx, "[userSubModuleRepoImpl][GetUserSubModuleBySubModuleIDAndUserID] error while row.Scan", "%v", err.Error())
		return resp, err
	}

	return resp, nil
}
