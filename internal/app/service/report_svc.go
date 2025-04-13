package service

import (
	"context"
	"fmt"
	"log/slog"
	"mahir-trade-be/internal/app/infra"
	"mahir-trade-be/internal/app/models"
	"mahir-trade-be/internal/app/repo/google"
	"mahir-trade-be/internal/app/repo/postgres"
	"mahir-trade-be/pkg/middleware"
	"math"
	"mime/multipart"
	"net/http"
	"strings"

	"go.uber.org/dig"
)

type (
	Report struct {
		ID int `json:"id"`
	}
	UploadContent struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}

	ReportSvc interface {
		CreateReport(ctx context.Context, req models.Report) (resp models.DefaultResponse, err error)
		GetReports(ctx context.Context, req models.PaginationRequest) (resp models.DefaultPaginationResponseData, err error)
		GetReportByID(ctx context.Context, id int64) (resp models.DefaultResponse, err error)
		UpdateReport(ctx context.Context, req models.Report) (resp models.DefaultResponse, err error)
		DeleteReport(ctx context.Context, id int64, deletedBy string) (resp models.DefaultResponse, err error)
		UploadContent(ctx context.Context, files map[string][]*multipart.FileHeader) (resp models.DefaultResponse, err error)
	}

	ReportSvcImpl struct {
		dig.In

		ReportRepo postgres.ReportRepo
		BucketRepo google.BucketRepo
		GoogleCfg  *infra.GoogleCfg
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
	if len(reports) > 0 {
		var reportsResp []models.Report
		for _, report := range reports {
			parsedThumbnailURL := r.BucketRepo.URLParser(report.ReportThumbnailURL)
			if report.ReportThumbnailURL != "" {
				report.ReportThumbnailURL, err = r.BucketRepo.PresignedURL(ctx, parsedThumbnailURL.BucketName, report.ReportThumbnailURL)
				if err != nil {
					slog.ErrorContext(ctx, fmt.Sprintf("[service][GetReports] while PresignedURL err : %v", err.Error()))
				}
			}
			reportsResp = append(reportsResp, report)
		}
		dataResp.Data = reportsResp
	} else {
		dataResp.Data = []models.Report{}
	}

	{
		dataResp.Code = http.StatusOK
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
	if report.ReportThumbnailURL == "" {
		parsedThumbnailURL := r.BucketRepo.URLParser(report.ReportThumbnailURL)
		if report.ReportThumbnailURL == "" {
			slog.InfoContext(ctx, fmt.Sprintf("[service][GetReportByID] while URLParser err : %v", parsedThumbnailURL))
			return resp, err
		}
		report.ReportThumbnailURL, err = r.BucketRepo.PresignedURL(ctx, parsedThumbnailURL.BucketName, report.ReportThumbnailURL)
		if err != nil {
			slog.ErrorContext(ctx, fmt.Sprintf("[service][GetReportByID] while PresignedURL err : %v", err.Error()))
		}
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

func (r *ReportSvcImpl) UploadContent(ctx context.Context, files map[string][]*multipart.FileHeader) (resp models.DefaultResponse, err error) {
	var uploadContents []UploadContent

	_, ok := ctx.Value(middleware.UserData).(middleware.UserCtxReq)
	if !ok {
		resp.Code = http.StatusUnauthorized
		resp.Message = "unauthorized"
		resp.Error = "user not found"
		slog.ErrorContext(ctx, "[service][UploadContent] user not found")

		return resp, err
	}

	for key, fileHeaders := range files {
		for _, fileHeader := range fileHeaders {
			fileName := fmt.Sprintf("contents/%s", strings.ToLower(key))
			blobFile, egErr := fileHeader.Open()
			if egErr != nil {
				slog.ErrorContext(ctx, "failed to open file", "egErr", fileHeader.Filename)
				err = fmt.Errorf("something went wrong")
				return
			}
			fileUpload := google.FileUpload{
				Filename:        fileName,
				Size:            fileHeader.Size,
				BucketName:      r.GoogleCfg.EducationBucketName,
				FileContentType: fileHeader.Header.Get("Content-Type"),
			}
			if strings.Contains(key, "thumbnail") {
				fileUpload.BucketName = r.GoogleCfg.FileBucketName
				fileUpload.Filename = fmt.Sprintf("thumbnails/%s", strings.ToLower(key))
			}
			uploadFileUrl, errUpload := r.BucketRepo.UploadStreamFile(ctx, fileUpload, blobFile)
			if errUpload != nil {
				slog.ErrorContext(ctx, "failed to upload file", "errUpload", fileHeader.Filename)
				err = fmt.Errorf("something went wrong")
				return
			}

			uploadContents = append(uploadContents, UploadContent{
				Key:   key,
				Value: uploadFileUrl,
			})

			errCloseFile := blobFile.Close()
			if errCloseFile != nil {
				slog.ErrorContext(ctx, "failed to close file", "errCloseFile", fileHeader.Filename)
				err = fmt.Errorf("something went wrong")
				return
			}
		}
	}

	resp.Code = http.StatusOK
	resp.Message = "Success"
	resp.Data = uploadContents

	return
}
