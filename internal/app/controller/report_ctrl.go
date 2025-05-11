package controller

import (
	"log/slog"
	"mahir-trade-be/internal/app/models"
	"mahir-trade-be/internal/app/service"
	"mahir-trade-be/internal/app/service/utils"
	"mahir-trade-be/pkg/middleware"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"go.uber.org/dig"
)

type (
	ReportCtrl interface {
		CreateReport(ec echo.Context) error
		GetReports(ec echo.Context) error
		GetReportByID(ec echo.Context) error
		UpdateReport(ec echo.Context) error
		DeleteReport(ec echo.Context) error
		UploadContent(ec echo.Context) error
	}

	ReportCtrlImpl struct {
		dig.In

		ReportSvc service.ReportSvc
	}
)

func NewReportCtrl(impl ReportCtrlImpl) ReportCtrl {
	return &impl
}

func (r *ReportCtrlImpl) CreateReport(ec echo.Context) error {
	ctx := ec.Request().Context()

	defer func() {
		if r := recover(); r != nil {
			slog.Error("CreateGroup - something went wrong", r)
		}
	}()

	userData, ok := ctx.Value(middleware.UserData).(middleware.UserCtxReq)
	if !ok {
		return ec.JSON(http.StatusUnauthorized, models.DefaultResponse{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized",
			Error:   "Unauthorized",
		})
	}

	var report models.Report
	if err := ec.Bind(&report); err != nil {
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	validate := utils.Validate
	err := validate.Struct(report)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   errors.Error(),
		})
	}

	if err := report.ValidateContents(); err != nil {
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	report.CreatedBy = userData.Email
	resp, err := r.ReportSvc.CreateReport(ctx, report)
	if err != nil {
		return ec.JSON(resp.Code, resp)
	}

	return ec.JSON(http.StatusCreated, resp)
}

func (r *ReportCtrlImpl) GetReports(ec echo.Context) error {
	ctx := ec.Request().Context()

	defer func() {
		if r := recover(); r != nil {
			slog.Error("GetReports - something went wrong", r)
		}
	}()

	limit, err := strconv.ParseInt(ec.QueryParam("limit"), 10, 64)
	if err != nil {
		limit = 10
	}

	page, err := strconv.ParseInt(ec.QueryParam("page"), 10, 64)
	if err != nil {
		page = 1
	}

	search := ec.QueryParam("search")
	sortBy := ec.QueryParam("sortBy")

	req := models.PaginationRequest{
		Limit:  limit,
		Page:   page,
		Search: search,
		SortBy: sortBy,
	}

	resp, err := r.ReportSvc.GetReports(ctx, req)
	if err != nil {
		slog.Error("GetReports - something went wrong", err)
		return ec.JSON(http.StatusBadRequest, resp)
	}

	return ec.JSON(http.StatusOK, resp)
}

func (r *ReportCtrlImpl) GetReportByID(ec echo.Context) error {
	ctx := ec.Request().Context()

	defer func() {
		if r := recover(); r != nil {
			slog.Error("GetPackageByID - something went wrong", r)
		}
	}()

	id, err := strconv.ParseInt(ec.Param("id"), 10, 64)
	if err != nil {
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	resp, err := r.ReportSvc.GetReportByID(ctx, id)
	if err != nil {
		return ec.JSON(http.StatusBadRequest, resp)
	}

	return ec.JSON(http.StatusOK, resp)
}

func (r *ReportCtrlImpl) UpdateReport(ec echo.Context) error {
	ctx := ec.Request().Context()

	defer func() {
		if r := recover(); r != nil {
			slog.Error("UpdateReport - something went wrong", r)
		}
	}()

	userData, ok := ctx.Value(middleware.UserData).(middleware.UserCtxReq)
	if !ok {
		return ec.JSON(http.StatusUnauthorized, models.DefaultResponse{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized",
			Error:   "Unauthorized",
		})
	}

	id, err := strconv.ParseInt(ec.Param("id"), 10, 64)
	if err != nil {
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	var report models.Report
	if err := ec.Bind(&report); err != nil {
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	report.ID = id
	validate := utils.Validate
	err = validate.Struct(report)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   errors.Error(),
		})
	}

	report.UpdatedBy = userData.Email
	resp, err := r.ReportSvc.UpdateReport(ctx, report)
	if err != nil {
		return ec.JSON(resp.Code, resp)
	}

	return ec.JSON(http.StatusOK, resp)
}

func (r *ReportCtrlImpl) DeleteReport(ec echo.Context) error {
	ctx := ec.Request().Context()

	defer func() {
		if r := recover(); r != nil {
			slog.Error("DeleteReport - something went wrong", r)
		}
	}()

	userData, ok := ctx.Value(middleware.UserData).(middleware.UserCtxReq)
	if !ok {
		return ec.JSON(http.StatusUnauthorized, models.DefaultResponse{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized",
			Error:   "Unauthorized",
		})
	}

	id, err := strconv.ParseInt(ec.Param("id"), 10, 64)
	if err != nil {
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	resp, err := r.ReportSvc.DeleteReport(ctx, id, userData.Email)
	if err != nil {
		return ec.JSON(resp.Code, resp)
	}
	return ec.JSON(http.StatusOK, resp)
}

func (r *ReportCtrlImpl) UploadContent(ec echo.Context) error {
	ctx := ec.Request().Context()

	defer func() {
		if r := recover(); r != nil {
			slog.Error("UploadContent - something went wrong", r)
		}
	}()

	files, _, err := utils.ParseMultipartForm(ec,
		[]utils.FileType{
			{MimeType: "image/jpeg", Size: 2},
			{MimeType: "image/jpg", Size: 2},
			{MimeType: "image/png", Size: 2},
			{MimeType: "image/webp", Size: 2},
		})
	if err != nil {
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	resp, err := r.ReportSvc.UploadContent(ctx, files)
	if err != nil {
		return ec.JSON(resp.Code, resp)
	}

	return ec.JSON(http.StatusOK, resp)
}
