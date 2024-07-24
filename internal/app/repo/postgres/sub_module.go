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
	SubModuleRepo interface {
		CreateSubModule(ctx context.Context, req models.SubModule) (id int, err error)
		GetSubModuleByID(ctx context.Context, id int64) (subModule models.SubModule, err error)
		GetSubModulesByModuleID(ctx context.Context, moduleID int64) (subModules []models.SubModule, err error)
		GetSubModules(ctx context.Context, req models.PaginationRequest) (subModules []models.SubModule, totalCount int64, err error)
		UpdateSubModule(ctx context.Context, req models.SubModule) (err error)
		SoftDeleteSubModule(ctx context.Context, subModuleId int64, operator string) (err error)
		RemoveModuleIDFromSubModules(ctx context.Context, moduleID int64, operator string) (err error)
	}

	SubModuleRepoImpl struct {
		dig.In

		*sql.DB
	}
)

func NewSubModuleRepo(impl SubModuleRepoImpl) SubModuleRepo {
	return &impl
}

func (s *SubModuleRepoImpl) CreateSubModule(ctx context.Context, req models.SubModule) (id int, err error) {
	var (
		row *sql.Row
	)

	if req.ModuleID.Valid {
		row = s.QueryRowContext(ctx, queries.QueryCreateSubModule, req.ModuleID, req.SubModuleName, req.Title, req.VideoURL, req.CreatedBy)
	} else {
		row = s.QueryRowContext(ctx, queries.QueryCreateSubModuleWithoutModuleID, req.SubModuleName, req.Title, req.VideoURL, req.CreatedBy)
	}

	err = row.Scan(&id)
	if err != nil {
		slog.ErrorContext(ctx, "[subModuleRepoImpl][CreateSubModule] error while row.Scan", "%v", err.Error())
		err = fmt.Errorf("something went wrong, we will fix it soon")
		return
	}

	return

}

func (s *SubModuleRepoImpl) GetSubModuleByID(ctx context.Context, id int64) (subModule models.SubModule, err error) {
	row := s.QueryRowContext(ctx, queries.QueryGetSubModuleByID, id)
	err = row.Scan(&subModule.ID, &subModule.UUID, &subModule.ModuleID, &subModule.SubModuleName, &subModule.Title, &subModule.VideoURL, &subModule.CreatedBy, &subModule.UpdatedBy, &subModule.CreatedAt, &subModule.UpdatedAt)
	if err != nil {
		return subModule, err
	}

	return subModule, nil
}

func (s *SubModuleRepoImpl) GetSubModules(ctx context.Context, req models.PaginationRequest) (subModules []models.SubModule, totalCount int64, err error) {
	rows, err := s.QueryContext(ctx, queries.QueryGetSubModules, req.Limit, req.Page)
	if err != nil {
		return subModules, totalCount, err
	}

	defer rows.Close()

	for rows.Next() {
		var subModule models.SubModule
		err = rows.Scan(&totalCount, &subModule.ID, &subModule.UUID, &subModule.ModuleID, &subModule.SubModuleName, &subModule.Title, &subModule.VideoURL, &subModule.CreatedBy, &subModule.UpdatedBy, &subModule.CreatedAt, &subModule.UpdatedAt)
		if err != nil {
			return subModules, totalCount, err
		}

		subModules = append(subModules, subModule)
	}

	return subModules, totalCount, nil
}

func (s *SubModuleRepoImpl) GetSubModulesByModuleID(ctx context.Context, moduleID int64) (subModules []models.SubModule, err error) {
	rows, err := s.QueryContext(ctx, queries.QueryGetSubModuleByModuleID, moduleID)
	if err != nil {
		slog.ErrorContext(ctx, "[subModuleRepoImpl][GetSubModulesByModuleID] error while QueryContext", "%v", err.Error())
		err = fmt.Errorf("internal server error")
		return
	}

	defer rows.Close()

	for rows.Next() {
		var subModule models.SubModule
		err = rows.Scan(&subModule.ID, &subModule.UUID, &subModule.ModuleID, &subModule.SubModuleName, &subModule.Title, &subModule.VideoURL, &subModule.CreatedBy, &subModule.UpdatedBy, &subModule.CreatedAt, &subModule.UpdatedAt)
		if err != nil {
			slog.ErrorContext(ctx, "[subModuleRepoImpl][GetSubModulesByModuleID] error while rows.Scan", "%v", err.Error())
			err = fmt.Errorf("data not match")
			return
		}

		subModules = append(subModules, subModule)
	}

	if len(subModules) == 0 {
		slog.ErrorContext(ctx, "[subModuleRepoImpl][GetSubModulesByModuleID] sub module with module id %d not found", moduleID)
		err = fmt.Errorf("sub module with module id %d not found", moduleID)
		return
	}

	return subModules, nil
}

func (s *SubModuleRepoImpl) UpdateSubModule(ctx context.Context, req models.SubModule) (err error) {
	result, err := s.ExecContext(ctx, queries.QueryUpdateSubModule, req.SubModuleName, req.Title, req.VideoURL, req.UpdatedBy, req.ModuleID, req.ID)
	if err != nil {
		slog.ErrorContext(ctx, "[subModuleRepoImpl][UpdateSubModule] error while ExecContext", "%v", err.Error())
		err = fmt.Errorf("internal server error")
		return
	}

	if affected, _ := result.RowsAffected(); affected == 0 {
		slog.ErrorContext(ctx, "[subModuleRepoImpl][UpdateSubModule] sub module with id %d not found", req.ID)
		err = fmt.Errorf("sub module with id %d not found", req.ID)
		return
	}

	return nil
}

func (s *SubModuleRepoImpl) SoftDeleteSubModule(ctx context.Context, subModuleId int64, operator string) (err error) {
	_, err = s.ExecContext(ctx, queries.QuerySoftDeleteSubModule, operator, operator, subModuleId)
	if err != nil {
		return err
	}

	return nil
}

func (s *SubModuleRepoImpl) RemoveModuleIDFromSubModules(ctx context.Context, moduleID int64, operator string) (err error) {
	_, err = s.ExecContext(ctx, queries.QueryRemoveModuleIDFromSubModules, operator, moduleID)
	if err != nil {
		slog.ErrorContext(ctx, "[subModuleRepoImpl][RemoveModuleIDFromSubModules] error while ExecContext", "%v", err.Error())
		err = fmt.Errorf("internal server error")
		return err
	}

	return
}
