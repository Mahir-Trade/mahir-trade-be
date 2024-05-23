package service

import (
	"context"
	"fmt"
	"log/slog"
	"mahir-trade-be/internal/app/models"
	"mahir-trade-be/internal/app/repo/postgres"
	"math"
	"net/http"

	"go.uber.org/dig"
)

type (
	ReportSvc interface {
		CreateReport(ctx context.Context, req models.Report) (resp models.DefaultResponse, err error)
		GetReports(ctx context.Context, req models.PaginationRequest) (resp models.DefaultPaginationResponseData, err error)
		GetReportByID(ctx context.Context, id int64) (resp models.DefaultResponse, err error)
		UpdateReport(ctx context.Context, req models.Report) (resp models.DefaultResponse, err error)
		DeleteReport(ctx context.Context, id int64, deletedBy string) (resp models.DefaultResponse, err error)
	}

	ReportSvcImpl struct {
		dig.In

		ReportRepo postgres.ReportRepo
	}
)

func NewReportSvc(impl ReportSvcImpl) ReportSvc {
	return &impl
}

func (r *ReportSvcImpl) CreateReport(ctx context.Context, req models.Report) (resp models.DefaultResponse, err error) {
	{
		resp.Code = http.StatusCreated
		resp.Message = "Success"
	}

	reportId, err := r.ReportRepo.CreateReport(ctx, req)
	if err != nil {
		resp.Code = http.StatusBadRequest
		resp.Message = "bad request"
		resp.Error = err.Error()
		slog.ErrorContext(ctx, fmt.Sprintf("[service][CreateReport] while CreateReport err : %v", err.Error()))

		return resp, err
	}

	type Report struct {
		ID int `json:"id"`
	}

	resp.Data = Report{ID: int(reportId)}

	return
}

func (r *ReportSvcImpl) GetReports(ctx context.Context, req models.PaginationRequest) (resp models.DefaultPaginationResponseData, err error) {
	var dataResp models.DefaultResponse
	{
		dataResp.Code = http.StatusOK
		dataResp.Message = "Success"
	}

	reports, totalCount, err := r.ReportRepo.GetReports(ctx, req)
	if err != nil {
		dataResp.Code = http.StatusBadRequest
		dataResp.Message = "bad request"
		dataResp.Error = err.Error()
		slog.ErrorContext(ctx, fmt.Sprintf("[service][GetReports] while GetReports err : %v", err.Error()))

		return resp, err
	}

	// convert to response
	{
		dataResp.Data = reports
		resp.Page = uint(req.Page)
		resp.Limit = uint(req.Limit)

		totalPage := math.Ceil(float64(totalCount) / float64(req.Limit))
		resp.TotalPages = uint(totalPage)

		resp.TotalItems = uint(totalCount)
		resp.HasNext = req.Page < int64(resp.TotalPages)
		resp.HasPrevious = req.Page > 1
		resp.Results = dataResp
	}

	return
}

func (r *ReportSvcImpl) GetReportByID(ctx context.Context, id int64) (resp models.DefaultResponse, err error) {
	{
		resp.Code = http.StatusOK
		resp.Message = "Success"
	}

	report, err := r.ReportRepo.GetReportByID(ctx, id)
	if err != nil {
		resp.Code = http.StatusBadRequest
		resp.Message = "bad request"
		resp.Error = err.Error()
		slog.ErrorContext(ctx, fmt.Sprintf("[service][GetReportByID] while GetReportByID err : %v", err.Error()))

		return resp, err
	}

	resp.Data = report

	return
}

func (r *ReportSvcImpl) UpdateReport(ctx context.Context, req models.Report) (resp models.DefaultResponse, err error) {
	{
		resp.Code = http.StatusOK
		resp.Message = "Success"
	}

	report, err := r.ReportRepo.GetReportByID(ctx, int64(req.ID))
	if err != nil {
		resp.Code = http.StatusBadRequest
		resp.Message = "bad request"
		resp.Error = err.Error()
		slog.ErrorContext(ctx, fmt.Sprintf("[service][GetReportByID] while GetReportByID err : %v", err.Error()))

		return resp, err
	}

	if report.ID == 0 {
		resp.Code = http.StatusBadRequest
		resp.Message = "bad request"
		resp.Error = "report not found"
		slog.ErrorContext(ctx, fmt.Sprintf("[service][UpdateReport] while GetReportByID err : report not found"))

		return resp, err
	}

	err = r.ReportRepo.UpdateReport(ctx, req)
	if err != nil {
		resp.Code = http.StatusBadRequest
		resp.Message = "bad request"
		resp.Error = err.Error()
		slog.ErrorContext(ctx, fmt.Sprintf("[service][UpdateReport] while UpdateReport err : %v", err.Error()))

		return resp, err
	}

	return
}

func (r *ReportSvcImpl) DeleteReport(ctx context.Context, id int64, deletedBy string) (resp models.DefaultResponse, err error) {
	{
		resp.Code = http.StatusOK
		resp.Message = "Success"
	}

	report, err := r.ReportRepo.GetReportByID(ctx, id)
	if err != nil {
		resp.Code = http.StatusBadRequest
		resp.Message = "bad request"
		resp.Error = err.Error()
		slog.ErrorContext(ctx, fmt.Sprintf("[service][DeleteReport] while GetReportByID err : %v", err.Error()))

		return resp, err
	}

	if report.ID == 0 {
		resp.Code = http.StatusBadRequest
		resp.Message = "bad request"
		resp.Error = "report not found"
		slog.ErrorContext(ctx, fmt.Sprintf("[service][DeleteReport] report with id %d not found", id))

		return resp, err
	}

	err = r.ReportRepo.SoftDeleteReport(ctx, id, deletedBy)
	if err != nil {
		resp.Code = http.StatusBadRequest
		resp.Message = "bad request"
		resp.Error = err.Error()
		slog.ErrorContext(ctx, fmt.Sprintf("[service][DeleteReport] while SoftDeleteReport err : %v", err.Error()))

		return resp, err
	}

	return
}
