package postgres

import (
	"context"
	"database/sql"
	"log/slog"
	"mahir-trade-be/internal/app/models"
	"mahir-trade-be/internal/app/repo/postgres/queries"

	"go.uber.org/dig"
)

type (
	TransactionRepo interface {
		CreateTransaction(ctx context.Context, req models.Transaction) (id int, err error)
	}

	TransactionRepoImpl struct {
		dig.In

		*sql.DB
	}
)

func NewTransactionRepo(impl TransactionRepoImpl) TransactionRepo {
	return &impl
}

func (t *TransactionRepoImpl) CreateTransaction(ctx context.Context, req models.Transaction) (id int, err error) {
	rows, err := t.QueryContext(ctx, queries.QueryCreateTransaction, req.OrderID, req.Amount, req.SettlementDate, req.WebhookID, req.CreatedBy)
	if err != nil {
		slog.ErrorContext(ctx, "[transactionRepoImpl][CreateTransaction] error while QueryContext", "%v", err.Error())
		return id, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			slog.ErrorContext(ctx, "[transactionRepoImpl][CreateTransaction] error while Scan", "%v", err.Error())
			return id, err
		}
	}

	return id, nil
}
