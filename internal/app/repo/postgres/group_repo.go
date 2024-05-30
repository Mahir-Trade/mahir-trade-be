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
	GroupRepo interface {
		CreateGroup(ctx context.Context, req models.Group) (id int, err error)
		GetGroupByID(ctx context.Context, id int64) (group models.Group, err error)
		GetGroups(ctx context.Context, req models.GetGroupsRequest) (groups []models.Group, totalCount int64, err error)
		UpdateGroup(ctx context.Context, req models.Group) (err error)
		SoftDeleteGroup(ctx context.Context, groupId int64, operator string) (err error)
	}

	GroupRepoImpl struct {
		dig.In

		*sql.DB
	}
)

func NewGroupRepo(impl GroupRepoImpl) GroupRepo {
	return &impl
}

func (g *GroupRepoImpl) CreateGroup(ctx context.Context, req models.Group) (id int, err error) {
	rows, err := g.QueryContext(ctx, queries.QueryCreateGroup, req.GroupName)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("error while CreateGroup err: %v", err.Error()))
		return id, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			slog.ErrorContext(ctx, fmt.Sprintf("error while CreateGroup err: %v", err.Error()))
			return id, err
		}
	}

	return id, nil
}

func (g *GroupRepoImpl) GetGroupByID(ctx context.Context, id int64) (group models.Group, err error) {
	row := g.QueryRowContext(ctx, queries.QueryGetGroupByID, id)
	err = row.Scan(&group.ID, &group.UUID, &group.GroupName, &group.CreatedAt, &group.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return group, fmt.Errorf("group with id %d not found", id)
		}

		slog.ErrorContext(ctx, fmt.Sprintf("error while GetGroup err: %v", err.Error()))
		return group, err
	}

	return group, nil
}

func (g *GroupRepoImpl) GetGroups(ctx context.Context, req models.GetGroupsRequest) (groups []models.Group, totalCount int64, err error) {
	args := []interface{}{}
	query := queries.QueryGetGroups

	if !req.ShowAll {
		if req.Limit == 0 {
			req.Limit = 10
		}

		offset := 0
		if req.Page > 1 {
			offset = int((req.Page - 1) * req.Limit)
		}

		args = append(args, req.Limit, offset)
		query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(args)-1, len(args))
	}

	rows, err := g.QueryContext(ctx, query, args...)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("error while GetGroups err: %v", err.Error()))
		return groups, totalCount, err
	}

	defer rows.Close()

	for rows.Next() {
		var group models.Group
		err = rows.Scan(&group.ID, &group.UUID, &group.GroupName, &group.CreatedAt, &group.UpdatedAt)
		if err != nil {
			slog.ErrorContext(ctx, fmt.Sprintf("error while GetGroups err: %v", err.Error()))
			return groups, totalCount, err
		}

		groups = append(groups, group)
	}

	// get total data
	row := g.QueryRowContext(ctx, queries.QueryGetTotalGroups)
	err = row.Scan(&totalCount)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("error while GetGroups err: %v", err.Error()))
		return groups, totalCount, err
	}

	return
}

func (g *GroupRepoImpl) UpdateGroup(ctx context.Context, req models.Group) (err error) {
	_, err = g.ExecContext(ctx, queries.QueryUpdateGroup, req.GroupName, req.ID)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("error while UpdateGroup err: %v", err.Error()))
		return err
	}

	return nil
}

func (g *GroupRepoImpl) SoftDeleteGroup(ctx context.Context, groupId int64, operator string) (err error) {
	_, err = g.ExecContext(ctx, queries.QuertSoftDeleteGroup, operator, operator, groupId)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("error while SoftDeleteGroup err: %v", err.Error()))
		return err
	}

	return nil
}
