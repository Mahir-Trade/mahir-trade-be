package postgres

import (
	"context"
	"database/sql"
	"mahir-trade-be/internal/app/models"
	"mahir-trade-be/internal/app/repo/postgres/queries"

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
