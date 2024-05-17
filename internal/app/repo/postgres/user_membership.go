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
