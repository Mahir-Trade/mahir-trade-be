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
	OrderRepo interface {
		CreateOrder(ctx context.Context, req models.Order) (id int, err error)
		GetOrderByPaymentCode(ctx context.Context, paymentCode string) (resp models.Order, err error)
		UpdateOrderStatus(ctx context.Context, req models.Order) (err error)
	}

	OrderRepoImpl struct {
		dig.In

		*sql.DB
	}
)

func NewOrderRepo(impl OrderRepoImpl) OrderRepo {
	return &impl
}

func (o *OrderRepoImpl) CreateOrder(ctx context.Context, req models.Order) (id int, err error) {
	if req.Status == "" {
		req.Status = "pending"
	}

	rows, err := o.QueryContext(ctx, queries.QueryCreateOrder, req.UserID, req.PackageID, req.Status, req.PaymentCode, req.PaymentURL, req.CreatedBy)
	if err != nil {
		slog.ErrorContext(ctx, "[orderRepoImpl][CreateOrder] error while QueryContext", "%v", err.Error())
		return id, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			slog.ErrorContext(ctx, "[orderRepoImpl][CreateOrder] error while Scan", "%v", err.Error())
			return id, err
		}
	}

	return id, nil
}

func (o *OrderRepoImpl) GetOrderByPaymentCode(ctx context.Context, paymentCode string) (resp models.Order, err error) {
	row := o.QueryRowContext(ctx, queries.QueryGetOrderByPaymentCode, paymentCode)
	err = row.Scan(&resp.ID, &resp.UserID, &resp.PackageID, &resp.Status, &resp.PaymentCode, &resp.PaymentURL, &resp.CreatedBy, &resp.UpdatedBy, &resp.CreatedAt, &resp.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return resp, fmt.Errorf("order with payment code %s not found", paymentCode)
		}

		slog.ErrorContext(ctx, "[orderRepoImpl][GetOrderByPaymentCode] error while Scan", "%v", err.Error())
		return resp, err
	}

	return resp, nil
}

func (o *OrderRepoImpl) UpdateOrderStatus(ctx context.Context, req models.Order) (err error) {
	_, err = o.ExecContext(ctx, queries.QueryUpdateOrderStatus, req.Status, req.UpdatedBy, req.ID)
	if err != nil {
		slog.ErrorContext(ctx, "[orderRepoImpl][UpdateOrderStatus] error while ExecContext", "%v", err.Error())
		return err
	}

	return nil
}
