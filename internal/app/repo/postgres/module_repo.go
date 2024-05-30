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
	ModuleRepo interface {
		CreateModule(ctx context.Context, req models.Module) (id int, err error)
		GetModuleByID(ctx context.Context, id int64) (module models.Module, err error)
		UpdateModule(ctx context.Context, req models.Module) (err error)
		GetModules(ctx context.Context, req models.PaginationRequest) (modules []models.Module, totalCount int64, err error)
		GetModulesByGroupID(ctx context.Context, groupID int64) (modules []models.Module, err error)
		SoftDeleteModule(ctx context.Context, moduleId int64, operator string) (err error)
	}

	ModuleRepoImpl struct {
		dig.In

		*sql.DB
	}
)

func NewModuleRepo(impl ModuleRepoImpl) ModuleRepo {
	return &impl
}
func (m *ModuleRepoImpl) CreateModule(ctx context.Context, req models.Module) (id int, err error) {
	var row *sql.Row

	if req.GroupID.Valid {
		if req.Tag.Valid {
			row = m.QueryRowContext(ctx, queries.QueryCreateModuleWithGroupIDAndTag, req.GroupID, req.ModuleName, req.ThumbnailUrl, req.Tag, req.CreatedBy)
		} else {
			row = m.QueryRowContext(ctx, queries.QueryCreateModuleWithGroupID, req.GroupID, req.ModuleName, req.ThumbnailUrl, req.CreatedBy)
		}
	} else {
		if req.Tag.Valid {
			row = m.QueryRowContext(ctx, queries.QueryCreateModuleWithoutGroupID, req.ModuleName, req.ThumbnailUrl, req.Tag, req.CreatedBy)
		} else {
			row = m.QueryRowContext(ctx, queries.QueryCreateModuleWithoutGroupIDAndTag, req.ModuleName, req.ThumbnailUrl, req.CreatedBy)
		}
	}

	err = row.Scan(&id)
	if err != nil {
		slog.ErrorContext(ctx, "[moduleRepoImpl][CreateModule] error while row.Scan", "%v", err.Error())
		err = fmt.Errorf("something went wrong, we will fix it soon")
		return
	}

	return
}

func (m *ModuleRepoImpl) GetModuleByID(ctx context.Context, id int64) (module models.Module, err error) {
	row := m.QueryRowContext(ctx, queries.QueryGetModuleByID, id)
	err = row.Scan(&module.ID, &module.UUID, &module.GroupID, &module.ModuleName, &module.ThumbnailUrl, &module.Tag, &module.CreatedBy, &module.CreatedAt, &module.UpdatedAt, &module.UpdatedBy)
	if err != nil {
		slog.ErrorContext(ctx, "[moduleRepoImpl][GetModuleByID] error while row.Scan", "%v", err.Error())
		return module, err
	}

	return module, nil
}

func (m *ModuleRepoImpl) UpdateModule(ctx context.Context, req models.Module) (err error) {
	var (
		result sql.Result
	)
	if req.ThumbnailUrl.Valid && req.Tag.Valid {
		result, err = m.ExecContext(ctx, queries.QueryUpdateModuleWithThumbnailAndTag, req.ModuleName, req.ThumbnailUrl, req.Tag, req.UpdatedBy, req.ID)
	} else if req.ThumbnailUrl.Valid {
		result, err = m.ExecContext(ctx, queries.QueryUpdateModuleWithThumbnail, req.ModuleName, req.ThumbnailUrl, req.UpdatedBy, req.ID)
	} else if req.Tag.Valid {
		result, err = m.ExecContext(ctx, queries.QueryUpdateModuleWithTag, req.ModuleName, req.Tag, req.UpdatedBy, req.ID)
	} else {
		result, err = m.ExecContext(ctx, queries.QueryUpdateModule, req.ModuleName, req.UpdatedBy, req.ID)
	}
	if err != nil {
		slog.ErrorContext(ctx, "[moduleRepoImpl][UpdateModule] error while ExecContext", "%v", err.Error())
		err = fmt.Errorf("something went wrong, we will fix it soon")
		return err
	}

	if affected, _ := result.RowsAffected(); affected == 0 {
		slog.ErrorContext(ctx, "[moduleRepoImpl][UpdateModule] error while RowsAffected", "%v", err.Error())
		err = fmt.Errorf("something went wrong, we will fix it soon")
		return err
	}

	return nil
}

func (m *ModuleRepoImpl) GetModules(ctx context.Context, req models.PaginationRequest) (modules []models.Module, totalCount int64, err error) {
	var (
		rows  *sql.Rows
		query = queries.QueryGetModules
	)

	args := []interface{}{}

	if !req.ShowAll {
		if req.Search != "" {
			query += " AND m.module_name ILIKE '%' || $1 || '%'"
			args = append(args, req.Search)
		}

		args = append(args, req.Limit, req.Page)
		query += fmt.Sprintf(" ORDER BY m.created_at DESC LIMIT $%d OFFSET $%d", len(args)-1, len(args))
	}

	rows, err = m.QueryContext(ctx, query, args...)
	if err != nil {
		slog.ErrorContext(ctx, "[moduleRepoImpl][searchNotExist] error while QueryContext", "%v", err.Error())
		err = fmt.Errorf("something went wrong, we will fix it soon")
		return
	}

	defer rows.Close()

	for rows.Next() {
		var module models.Module
		err = rows.Scan(&totalCount, &module.ID, &module.UUID, &module.GroupID, &module.ModuleName, &module.ThumbnailUrl, &module.Tag, &module.CreatedBy, &module.CreatedAt, &module.UpdatedAt, &module.UpdatedBy)
		if err != nil {
			slog.ErrorContext(ctx, "[moduleRepoImpl][scan] error while rows.Scan", "%v", err.Error())
			return
		}

		modules = append(modules, module)
	}

	return
}

func (m *ModuleRepoImpl) GetModulesByGroupID(ctx context.Context, groupID int64) (modules []models.Module, err error) {
	rows, err := m.QueryContext(ctx, queries.QueryGetModulesByGroupID, groupID)
	if err != nil {
		slog.ErrorContext(ctx, "[moduleRepoImpl][GetModulesByGroupID] error while QueryContext", "%v", err.Error())
		return
	}

	defer rows.Close()

	for rows.Next() {
		var module models.Module
		err = rows.Scan(&module.ID, &module.UUID, &module.GroupID, &module.ModuleName, &module.ThumbnailUrl, &module.CreatedBy, &module.CreatedAt, &module.UpdatedAt, &module.UpdatedBy)
		if err != nil {
			slog.ErrorContext(ctx, "[moduleRepoImpl][GetModulesByGroupID] error while rows.Scan", "%v", err.Error())
			return
		}

		modules = append(modules, module)
	}

	return
}

func (m *ModuleRepoImpl) SoftDeleteModule(ctx context.Context, moduleId int64, operator string) (err error) {
	result, err := m.ExecContext(ctx, queries.QuerySoftDeleteModule, operator, operator, moduleId)
	if err != nil {
		slog.ErrorContext(ctx, "[moduleRepoImpl][SoftDeleteModule] error while ExecContext", "%v", err.Error())
		err = fmt.Errorf("something went wrong, we will fix it soon")
		return err
	}

	if affected, _ := result.RowsAffected(); affected == 0 {
		slog.ErrorContext(ctx, "[moduleRepoImpl][SoftDeleteModule] error while RowsAffected", "%v", err.Error())
		err = fmt.Errorf("something went wrong, we will fix it soon")
		return err
	}

	return nil
}
