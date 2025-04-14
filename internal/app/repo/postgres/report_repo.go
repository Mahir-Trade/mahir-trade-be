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
		GetReports(ctx context.Context, req models.PaginationRequest) (reports []models.Report, totalCount int64, err error)
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
	rows, err := r.QueryContext(ctx, queries.QueryCreateReport, req.ReportName, req.ReportThumbnailURL, req.Contents, req.CreatedBy)
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

func (r *ReportRepoImpl) GetReports(ctx context.Context, req models.PaginationRequest) (reports []models.Report, totalCount int64, err error) {
	if req.Limit == 0 {
		req.Limit = 10
	}

	offset := 0
	if req.Page > 1 {
		offset = int((req.Page - 1) * req.Limit)
	}

	query := queries.QueryGetReports
	queryParams := []interface{}{}

	if req.Search != "" {
		query += " AND report_name ILIKE '%' || $1 || '%'"
		queryParams = append(queryParams, req.Search)
	}

	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(queryParams)+1, len(queryParams)+2)
	queryParams = append(queryParams, req.Limit, offset)

	rows, err := r.QueryContext(ctx, query, queryParams...)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("error while GetReports err: %v", err.Error()))
		return reports, totalCount, err
	}

	defer rows.Close()

	for rows.Next() {
		var report models.Report
		err = rows.Scan(&totalCount, &report.ID, &report.ReportName, &report.ReportThumbnailURL, &report.Contents, &report.CreatedBy, &report.UpdatedBy, &report.CreatedAt, &report.UpdatedAt)
		if err != nil {
			slog.ErrorContext(ctx, fmt.Sprintf("error while GetReports err: %v", err.Error()))
			return reports, totalCount, err
		}

		reports = append(reports, report)
	}
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("error while marshalling reports err: %v", err.Error()))
		return reports, totalCount, err
	}
	return reports, totalCount, nil
}

func (r *ReportRepoImpl) GetReportByID(ctx context.Context, id int64) (report models.Report, err error) {
	row := r.QueryRowContext(ctx, queries.QueryGetReportByID, id)

	err = row.Scan(&report.ID, &report.ReportName, &report.ReportThumbnailURL, &report.Contents, &report.CreatedBy, &report.UpdatedBy, &report.CreatedAt, &report.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return report, fmt.Errorf("report not found")
		}

		slog.ErrorContext(ctx, fmt.Sprintf("error while GetReportByID err: %v", err.Error()))
		return report, err
	}

	return report, nil
}

func (r *ReportRepoImpl) UpdateReport(ctx context.Context, req models.Report) (err error) {
	_, err = r.ExecContext(ctx, queries.QueryUpdateReport, req.ReportName, req.ReportThumbnailURL, req.Contents, req.UpdatedBy, req.ID)
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
