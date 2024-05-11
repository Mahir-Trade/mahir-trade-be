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
	ReportRepo interface {
		CreateReport(ctx context.Context, req models.Report) (id int64, err error)
		GetReports(ctx context.Context, req models.GetPackagesRequest) (reports []models.Report, totalCount int64, err error)
		GetReportByID(ctx context.Context, id int64) (report models.Report, err error)
		UpdateReport(ctx context.Context, req models.Report) (err error)
		SoftDeleteReport(ctx context.Context, id int64, deletedBy string) (err error)
	}

	ReportRepoImpl struct {
		dig.In

		*sql.DB
	}
)

func NewReportRepo(impl ReportRepoImpl) ReportRepo {
	return &impl
}

func (r *ReportRepoImpl) CreateReport(ctx context.Context, req models.Report) (id int64, err error) {
	rows, err := r.QueryContext(ctx, queries.QueryCreateReport, req.ReportThumbnailURL, req.ReportFileURL)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("error while CreateReport err: %v", err.Error()))
		return id, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			slog.ErrorContext(ctx, fmt.Sprintf("error while CreateReport err: %v", err.Error()))
			return id, err
		}
	}

	return id, nil
}

func (r *ReportRepoImpl) GetReports(ctx context.Context, req models.GetPackagesRequest) (reports []models.Report, totalCount int64, err error) {
	if req.Limit == 0 {
		req.Limit = 10
	}

	offset := 0
	if req.Page > 1 {
		offset = int((req.Page - 1) * req.Limit)
	}

	rows, err := r.QueryContext(ctx, queries.QueryGetReports, req.Limit, offset)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("error while GetReports err: %v", err.Error()))
		return reports, totalCount, err
	}

	defer rows.Close()

	for rows.Next() {
		var report models.Report
		err = rows.Scan(&report.ID, &report.ReportThumbnailURL, &report.ReportFileURL, &report.CreatedBy, &report.UpdatedBy, &report.CreatedAt, &report.UpdatedAt)
		if err != nil {
			slog.ErrorContext(ctx, fmt.Sprintf("error while GetReports err: %v", err.Error()))
			return reports, totalCount, err
		}

		reports = append(reports, report)
	}

	totalCountRow := r.QueryRowContext(ctx, queries.QueryGetTotalReports)
	err = totalCountRow.Scan(&totalCount)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("error while GetTotalReports err: %v", err.Error()))
		return reports, totalCount, err
	}

	return reports, totalCount, nil
}

func (r *ReportRepoImpl) GetReportByID(ctx context.Context, id int64) (report models.Report, err error) {
	row := r.QueryRowContext(ctx, queries.QueryGetReportByID, id)

	err = row.Scan(&report.ID, &report.ReportThumbnailURL, &report.ReportFileURL, &report.CreatedBy, &report.UpdatedBy, &report.CreatedAt, &report.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return report, fmt.Errorf("package with id %d not found", id)
		}

		slog.ErrorContext(ctx, fmt.Sprintf("error while GetReportByID err: %v", err.Error()))
		return report, err
	}

	return report, nil
}

func (r *ReportRepoImpl) UpdateReport(ctx context.Context, req models.Report) (err error) {
	_, err = r.ExecContext(ctx, queries.QueryUpdateReport, req.ReportThumbnailURL, req.ReportFileURL, req.ID)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("error while UpdateReport err: %v", err.Error()))
		return err
	}

	return nil
}

func (r *ReportRepoImpl) SoftDeleteReport(ctx context.Context, id int64, deletedBy string) (err error) {
	_, err = r.ExecContext(ctx, queries.QuertSoftDeleteReport, deletedBy, deletedBy, id)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("error while SoftDeleteReport err: %v", err.Error()))
		return err
	}

	return nil
}
