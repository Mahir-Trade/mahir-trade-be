package postgres

import (
	"context"
	"database/sql"

	"go.uber.org/dig"
)

type (
	GeneralLog struct {
		ID        int    `json:"id"`
		UserID    int    `json:"user_id"`
		RawBody   string `json:"raw_body"`
		CreatedBy string `json:"created_by"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
		UpdatedBy string `json:"updated_by"`
	}

	GeneralLogRepo interface {
		CreateGeneralLog(ctx context.Context, req GeneralLog) error
	}

	GeneralLogRepoImpl struct {
		dig.In

		*sql.DB
	}
)

func NewGeneralLogRepo(impl GeneralLogRepoImpl) GeneralLogRepo {
	return &impl
}

func (g *GeneralLogRepoImpl) CreateGeneralLog(ctx context.Context, req GeneralLog) error {
	_, err := g.ExecContext(ctx, "INSERT INTO general_logs (user_id, raw_body, created_by) VALUES ($1, $2, $3)", req.UserID, req.RawBody, req.CreatedBy)
	if err != nil {
		return err
	}

	return nil
}
