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
	UserMembershipRepo interface {
		CreateUserMembership(ctx context.Context, req models.UserMembership) (id int, err error)
		UpdateUserMembershipByUserID(ctx context.Context, req models.UserMembership) (err error)
		UpdateBulkUserMembership(ctx context.Context) (err error)
		GetUserMembershipByUserID(ctx context.Context, userID int64) (resp models.UserMembership, err error)
	}

	UserMembershipRepoImpl struct {
		dig.In

		*sql.DB
	}
)

func NewUserMembershipRepo(impl UserMembershipRepoImpl) UserMembershipRepo {
	return &impl
}

func (u *UserMembershipRepoImpl) CreateUserMembership(ctx context.Context, req models.UserMembership) (id int, err error) {
	rows, err := u.QueryContext(ctx, queries.QueryCreateUserMembership, req.UserID, req.PackageID, req.ExpiredAt, req.IsMembershipActive, req.CreatedBy)
	if err != nil {
		slog.ErrorContext(ctx, "[userMembershipRepoImpl][CreateUserMembership] error while QueryContext", "%v", err.Error())
		return id, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			slog.ErrorContext(ctx, "[userMembershipRepoImpl][CreateUserMembership] error while Scan", "%v", err.Error())
			return id, err
		}
	}

	return id, nil
}

func (u UserMembershipRepoImpl) UpdateUserMembershipByUserID(ctx context.Context, req models.UserMembership) (err error) {
	_, err = u.ExecContext(ctx, queries.QueryUpdateUserMembershipExpired, req.ExpiredAt, req.IsMembershipActive, req.UpdatedBy, req.UserID)
	if err != nil {
		slog.ErrorContext(ctx, "[userMembershipRepoImpl][UpdateUserMembershipByUserID] error while ExecContext", "%v", err.Error())
		return err
	}

	return nil
}

func (u *UserMembershipRepoImpl) GetUserMembershipByUserID(ctx context.Context, userID int64) (resp models.UserMembership, err error) {
	row := u.QueryRowContext(ctx, queries.QueryGetUserMembershipByUserID, userID)
	err = row.Scan(&resp.ID, &resp.UserID, &resp.PackageID, &resp.ExpiredAt, &resp.IsMembershipActive, &resp.CreatedBy, &resp.UpdatedBy, &resp.CreatedAt, &resp.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return resp, nil
		}
		slog.ErrorContext(ctx, "[userMembershipRepoImpl][GetUserMembershipByUserID] error while Scan", "%v", err.Error())
		return resp, err
	}

	return resp, nil
}

func (u *UserMembershipRepoImpl) GetUserMemberships(ctx context.Context) (resp []models.UserMembership, err error) {
	rows, err := u.QueryContext(ctx, queries.QueryGetUserMemberships)
	if err != nil {
		slog.ErrorContext(ctx, "[userMembershipRepoImpl][GetUserMemberships] error while QueryContext", "%v", err.Error())
		return resp, err
	}

	defer rows.Close()

	for rows.Next() {
		var userMembership models.UserMembership
		err = rows.Scan(&userMembership.ID, &userMembership.UserID, &userMembership.PackageID, &userMembership.ExpiredAt, &userMembership.IsMembershipActive, &userMembership.CreatedBy, &userMembership.UpdatedBy, &userMembership.CreatedAt, &userMembership.UpdatedAt)
		if err != nil {
			slog.ErrorContext(ctx, "[userMembershipRepoImpl][GetUserMemberships] error while Scan", "%v", err.Error())
			return resp, err
		}

		resp = append(resp, userMembership)
	}

	return resp, nil
}

func (u *UserMembershipRepoImpl) UpdateBulkUserMembership(ctx context.Context) (err error) {
	tx, err := u.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		slog.ErrorContext(ctx, "[userMembershipRepoImpl][UpdateBulkUserMembership] error while BeginTx", "%v", err.Error())
		return err
	}

	_, err = tx.ExecContext(ctx, queries.QueryBulkUpdateUserMembershipExpired)
	if err != nil {
		tx.Rollback()
		slog.ErrorContext(ctx, "[userMembershipRepoImpl][UpdateBulkUserMembership] error while ExecContext", "%v", err.Error())
		return err
	}

	err = tx.Commit()
	if err != nil {
		slog.ErrorContext(ctx, "[userMembershipRepoImpl][UpdateBulkUserMembership] error while Commit", "%v", err.Error())
		return err
	}

	return nil
}
