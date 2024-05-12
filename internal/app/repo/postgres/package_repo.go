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
	PackageRepo interface {
		CreatePackage(ctx context.Context, req models.Package) (id int, err error)
		GetPackages(ctx context.Context, req models.GetPackagesRequest) (packages []models.Package, totalCount int64, err error)
		GetPackageByID(ctx context.Context, id int64) (res models.Package, err error)
		UpdatePackage(ctx context.Context, req models.Package) (err error)
		SoftDeletePackage(ctx context.Context, id int64, deletedBy string) (err error)
	}

	PackageRepoImpl struct {
		dig.In

		*sql.DB
	}
)

func NewPackageRepo(impl PackageRepoImpl) PackageRepo {
	return &impl
}

func (p *PackageRepoImpl) CreatePackage(ctx context.Context, req models.Package) (id int, err error) {
	rows, err := p.QueryContext(ctx, queries.QueryCreatePackage, req.Price, req.DurationInMonth, req.Description)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("error while CreatePackage err: %v", err.Error()))
		return id, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			slog.ErrorContext(ctx, fmt.Sprintf("error while CreatePackage err: %v", err.Error()))
			return id, err
		}
	}

	return id, nil
}

func (p *PackageRepoImpl) GetPackages(ctx context.Context, req models.GetPackagesRequest) (packages []models.Package, totalCount int64, err error) {
	if req.Limit == 0 {
		req.Limit = 10
	}

	offset := 0
	if req.Page > 1 {
		offset = int((req.Page - 1) * req.Limit)
	}

	rows, err := p.QueryContext(ctx, queries.QueryGetPackages, req.Limit, offset)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("error while GetPackages err: %v", err.Error()))
		return packages, totalCount, err
	}

	defer rows.Close()

	for rows.Next() {
		var pack models.Package
		err = rows.Scan(&pack.ID, &pack.Price, &pack.DurationInMonth, &pack.Description, &pack.CreatedBy, &pack.UpdatedBy, &pack.CreatedAt, &pack.UpdatedAt)
		if err != nil {
			slog.ErrorContext(ctx, fmt.Sprintf("error while GetPackages err: %v", err.Error()))
			return packages, totalCount, err
		}

		packages = append(packages, pack)
	}

	totalCountRow := p.QueryRowContext(ctx, queries.QueryGetTotalPackages)
	err = totalCountRow.Scan(&totalCount)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("error while GetTotalPackages err: %v", err.Error()))
		return packages, totalCount, err
	}

	return packages, totalCount, nil
}

func (p *PackageRepoImpl) GetPackageByID(ctx context.Context, id int64) (res models.Package, err error) {
	row := p.QueryRowContext(ctx, queries.QueryGetPackageByID, id)
	err = row.Scan(&res.ID, &res.Price, &res.DurationInMonth, &res.Description, &res.CreatedBy, &res.UpdatedBy, &res.CreatedAt, &res.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return res, fmt.Errorf("package with id %d not found", id)
		}

		slog.ErrorContext(ctx, fmt.Sprintf("error while GetPackageByID err: %v", err.Error()))
		return res, err
	}

	return res, nil
}

func (p *PackageRepoImpl) UpdatePackage(ctx context.Context, req models.Package) (err error) {
	_, err = p.ExecContext(ctx, queries.QueryUpdatePackage, req.Price, req.DurationInMonth, req.Description, req.ID)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("error while UpdatePackage err: %v", err.Error()))
		return err
	}

	return nil
}

func (p *PackageRepoImpl) SoftDeletePackage(ctx context.Context, id int64, deletedBy string) (err error) {
	_, err = p.ExecContext(ctx, queries.QuertSoftDeletePackage, deletedBy, deletedBy, id)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("error while SoftDeletePackage err: %v", err.Error()))
		return err
	}

	return nil
}
