package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"mahir-trade-be/internal/app/models"
	"mahir-trade-be/internal/app/repo/postgres/queries"
	"strings"

	"go.uber.org/dig"
)

type (
	AdminRepo interface {
		FindByUsername(ctx context.Context, username string) (admin models.Admin, err error)
		CreateAdmin(ctx context.Context, req models.Admin) (id int, err error)
		SoftDeleteAdmin(ctx context.Context, id int64, operator string) (err error)
	}

	AdminRepoImpl struct {
		dig.In

		*sql.DB
	}
)

func NewAdminRepo(impl AdminRepoImpl) AdminRepo {
	return &impl
}

func (a *AdminRepoImpl) FindByUsername(ctx context.Context, username string) (admin models.Admin, err error) {
	err = a.QueryRowContext(ctx, queries.QueryFindByEmail, username).Scan(&admin.AdminID, &admin.UUID, &admin.Email, &admin.Username, &admin.Password)
	return admin, err
}

func (a *AdminRepoImpl) CreateAdmin(ctx context.Context, req models.Admin) (id int, err error) {
	var rows *sql.Rows

	rows, err = a.QueryContext(ctx, queries.QueryCreateAdmin, req.Email, req.Username, req.Password)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("error while CreateAdmin err: %v", err.Error()))
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return id, fmt.Errorf("email or username already exist")
		}

		return id, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			return id, err
		}
	}

	return id, err
}

func (a *AdminRepoImpl) SoftDeleteAdmin(ctx context.Context, id int64, operator string) (err error) {
	_, err = a.ExecContext(ctx, queries.QuerySoftDeleteAdmin, operator, operator, id)
	return err
}
