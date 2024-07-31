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
	UserRepo interface {
		CreateUser(ctx context.Context, req models.User) (id int, err error)
		FindUserByEmailOrUsername(ctx context.Context, req string) (user models.User, err error)
		FindUserByEmailAndUsername(ctx context.Context, email, username string) (user models.User, err error)
		GetUserByID(ctx context.Context, id int64) (user models.User, err error)
		GetUserByUUID(ctx context.Context, uuid string) (user models.User, err error)
		UpdateTypeUser(ctx context.Context, isActive bool, operator string, id int64) (typeUser bool, err error)
		GetAllUser(ctx context.Context, req models.PaginationRequest) (users []models.GetUsersBOResponse, totalCount int64, err error)
	}

	UserRepoImpl struct {
		dig.In

		*sql.DB
	}
)

func NewUserRepo(impl UserRepoImpl) UserRepo {
	return &impl
}

func (u *UserRepoImpl) CreateUser(ctx context.Context, req models.User) (id int, err error) {
	rows, err := u.QueryContext(ctx, queries.QueryCreateUser, req.Email, req.PhoneNumber, req.Username, req.Password)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return id, fmt.Errorf("email or username already exist")
		}

		slog.ErrorContext(ctx, fmt.Sprintf("error while CreateUser err: %v", err.Error()))
		return id, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			return id, err
		}
	}

	return id, nil
}

func (u *UserRepoImpl) FindUserByEmailOrUsername(ctx context.Context, req string) (user models.User, err error) {
	row := u.QueryRowContext(ctx, queries.QueryFindUserByEmailOrUsename, req)
	err = row.Scan(&user.UserID, &user.UUID, &user.PhoneNumber, &user.Username, &user.Email, &user.Password)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("error while FindUserByEmailOrUsername err: %v", err.Error()))
		return user, err
	}

	return user, nil
}

func (u *UserRepoImpl) GetUserByID(ctx context.Context, id int64) (user models.User, err error) {
	row := u.QueryRowContext(ctx, queries.QueryGetUserByID, id)
	err = row.Scan(&user.UserID, &user.UUID, &user.PhoneNumber, &user.Username, &user.Email, &user.Password, &user.IsActive)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("error while GetUserByID err: %v", err.Error()))

		if err == sql.ErrNoRows {
			return user, fmt.Errorf("user not found")
		}

		return user, err
	}

	return user, nil
}

func (u *UserRepoImpl) GetUserByUUID(ctx context.Context, uuid string) (user models.User, err error) {
	row := u.QueryRowContext(ctx, queries.QueryGetUserByUUID, uuid)
	err = row.Scan(&user.UserID, &user.UUID, &user.PhoneNumber, &user.Username, &user.Email, &user.Password)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("error while GetUserByUUID err: %v", err.Error()))
		return user, err
	}

	return user, nil
}

func (u *UserRepoImpl) UpdateTypeUser(ctx context.Context, isActive bool, operator string, id int64) (typeUser bool, err error) {
	_, err = u.QueryContext(ctx, queries.QueryUpdateTypeUser, isActive, operator, id)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("error while UpdateTypeUser err: %v", err.Error()))
		return typeUser, err
	}

	return true, nil
}
func (u *UserRepoImpl) FindUserByEmailAndUsername(ctx context.Context, email, username string) (user models.User, err error) {
	row := u.QueryRowContext(ctx, queries.QueryFindUserByEmailAndUsername, email, username)
	err = row.Scan(&user.UserID, &user.UUID, &user.PhoneNumber, &user.Username, &user.Email, &user.Password)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("error while FindUserByEmailAndUsername err: %v", err.Error()))
		return user, err
	}

	return user, nil
}

func (u *UserRepoImpl) GetAllUser(ctx context.Context, req models.PaginationRequest) (users []models.GetUsersBOResponse, totalCount int64, err error) {
	query := queries.QueryGetUsers
	queryParams := []interface{}{}

	if req.Search != "" {
		// currently only can find by email
		query += " AND u.email ILIKE '%' || $1 || '%'"
		queryParams = append(queryParams, strings.ToLower(req.Search))
	}

	query += fmt.Sprintf(" ORDER BY u.created_at DESC LIMIT $%d OFFSET $%d", len(queryParams)+1, len(queryParams)+2)
	queryParams = append(queryParams, req.Limit, int((req.Page-1)*req.Limit))

	rows, err := u.QueryContext(ctx, query, queryParams...)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("[repo][GetAllUsers] while QueryContext, err: %v", err.Error()))
		return
	}

	defer rows.Close()

	for rows.Next() {
		var user models.GetUsersBOResponse
		err = rows.Scan(
			&totalCount,
			&user.UserID,
			&user.UUID,
			&user.PhoneNumber,
			&user.Email,
			&user.Username,
			&user.IsActive,
			&user.AccountType,
			&user.MembershipExpiredData,
			&user.CreatedAt,
			&user.CreatedBy,
		)
		if err != nil {
			slog.ErrorContext(ctx, fmt.Sprintf("[repo][GetAlltUser] while Scan, err: %v", err.Error()))
			return
		}

		users = append(users, user)
	}

	return
}
