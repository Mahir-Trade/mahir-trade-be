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
	DiscordAccountRepo interface {
		CreateDiscordAccount(ctx context.Context, req models.DiscordAccount) (id int, err error)
		GetDiscordAccountByUserID(ctx context.Context, userID int64) (discordAccount models.DiscordAccount, err error)
		DeleteDiscordAccountByUserID(ctx context.Context, id int64) (err error)
	}

	DiscordAccountRepoImpl struct {
		dig.In

		*sql.DB
	}
)

func NewDiscordAccountRepo(impl DiscordAccountRepoImpl) DiscordAccountRepo {
	return &impl
}

func (d *DiscordAccountRepoImpl) CreateDiscordAccount(ctx context.Context, req models.DiscordAccount) (id int, err error) {
	rows, err := d.QueryContext(ctx, queries.QueryCreateDiscordAccount, req.UserID, req.DiscordAccountID, req.Username, req.Email)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("error while CreateDiscordAccount err: %v", err.Error()))
		return id, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			slog.ErrorContext(ctx, fmt.Sprintf("error while CreateDiscordAccount err: %v", err.Error()))
			return id, err
		}
	}

	return id, nil
}

func (d *DiscordAccountRepoImpl) GetDiscordAccountByUserID(ctx context.Context, userID int64) (discordAccount models.DiscordAccount, err error) {
	row := d.QueryRowContext(ctx, queries.QueryGetDiscordAccountByUserID, userID)
	err = row.Scan(&discordAccount.ID, &discordAccount.UserID, &discordAccount.DiscordAccountID, &discordAccount.Username, &discordAccount.Email, &discordAccount.CreatedAt, &discordAccount.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return discordAccount, fmt.Errorf("discord account with user id %d not found", userID)
		}

		slog.ErrorContext(ctx, fmt.Sprintf("error while GetGroup err: %v", err.Error()))
		return discordAccount, err
	}

	return discordAccount, nil
}

func (d *DiscordAccountRepoImpl) DeleteDiscordAccountByUserID(ctx context.Context, id int64) (err error) {
	_, err = d.ExecContext(ctx, queries.QuerySoftDeleteDiscordAccount, id)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("error while DeleteDiscordAccountByUserID err: %v", err.Error()))
		return err
	}

	return nil
}
