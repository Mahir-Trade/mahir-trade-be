package postgres

import (
	"context"
	"database/sql"
	"mahir-trade-be/internal/app/repo/postgres/queries"

	"go.uber.org/dig"
)

type (
	EmailTemplateRepo interface {
		GetByKey(ctx context.Context, key string) (body string, err error)
	}

	EmailTemplateRepoImpl struct {
		dig.In

		*sql.DB
	}
)

func NewEmailTemplateRepo(impl EmailTemplateRepoImpl) EmailTemplateRepo {
	return &impl
}

func (e *EmailTemplateRepoImpl) GetByKey(ctx context.Context, key string) (body string, err error) {
	err = e.QueryRowContext(ctx, queries.QueryGetByKey, key).Scan(&body)
	if err != nil {
		return
	}

	return
}
