package controller

import (
	"log/slog"
	"mahir-trade-be/internal/app/models"
	"mahir-trade-be/internal/app/service"
	"mahir-trade-be/internal/app/service/utils"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"go.uber.org/dig"
)

type (
	AdminCtrl interface {
		AdminLogin(ec echo.Context) error
		AdminRegistration(ec echo.Context) error
		UpdateTypeUser(ec echo.Context) error
	}

	AdminCtrlImpl struct {
		dig.In

		AdminSvc service.AdminSvc
	}
)

func NewAdminCtrl(impl AdminCtrlImpl) AdminCtrl {
	return &impl
}

func (ox *AdminCtrlImpl) AdminLogin(ec echo.Context) error {
	ctx := ec.Request().Context()

	defer func() {
		if r := recover(); r != nil {
			slog.Error("CreateGroup - something went wrong", r)
		}
	}()

	var admin service.AdminLoginRequest

	if err := ec.Bind(&admin); err != nil {
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	validate := utils.Validate

	err := validate.Struct(admin)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   errors,
		})
	}

	resp, err := ox.AdminSvc.AdminLogin(ctx, admin)
	if err != nil {
		slog.ErrorContext(ctx, "[controller][AdminLogin] error while AdminLogin err: %v", err)
		return ec.JSON(http.StatusInternalServerError, resp)
	}

	return ec.JSON(resp.Code, resp)
}

func (ox *AdminCtrlImpl) AdminRegistration(ec echo.Context) error {
	ctx := ec.Request().Context()

	defer func() {
		if r := recover(); r != nil {
			slog.Error("CreateGroup - something went wrong", r)
		}
	}()

	var admin models.Admin

	if err := ec.Bind(&admin); err != nil {
		slog.ErrorContext(ctx, "[controller][AdminRegistration] error while ec.Bind(&admin) err: %v", err)
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	validate := utils.Validate

	err := validate.Struct(admin)
	if err != nil {
		slog.ErrorContext(ctx, "[controller][AdminRegistration] error while validate.Struct(admin) err: %v", err)
		errors := err.(validator.ValidationErrors)
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   errors,
		})
	}

	resp, err := ox.AdminSvc.AdminRegistration(ctx, admin)
	if err != nil {
		slog.ErrorContext(ctx, "[controller][AdminRegistration] error while AdminRegistration err: %v", err)
		return ec.JSON(http.StatusInternalServerError, resp)
	}

	return ec.JSON(resp.Code, resp)
}

func (ox *AdminCtrlImpl) UpdateTypeUser(ec echo.Context) error {
	ctx := ec.Request().Context()

	defer func() {
		if r := recover(); r != nil {
			slog.Error("CreateGroup - something went wrong", r)
		}
	}()

	userID := ec.Param("user_id")
	if userID == "" {
		slog.ErrorContext(ctx, "[controller][UpdateTypeUser] error while ec.Param(\"user_id\") err:")
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "user id is required",
		})
	}

	id, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		slog.ErrorContext(ctx, "[controller][UpdateTypeUser] error while strconv.ParseInt(userID, 10, 64) err: %v", err)
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "user id is required",
		})
	}

	var req service.UpdateTypeUserRequest

	if err := ec.Bind(&req); err != nil {
		slog.ErrorContext(ctx, "[controller][UpdateTypeUser] error while ec.Bind(&req) err: %v", err)
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	validate := utils.Validate

	err = validate.Struct(req)
	if err != nil {
		slog.ErrorContext(ctx, "[controller][UpdateTypeUser] error while validate.Struct(req) err: %v", err)
		errors := err.(validator.ValidationErrors)
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   errors,
		})
	}

	resp, err := ox.AdminSvc.UpdateTypeUser(ctx, req.IsActive, id)
	if err != nil {
		slog.ErrorContext(ctx, "[controller][UpdateTypeUser] error while UpdateTypeUser err: %v", err)
		return ec.JSON(http.StatusInternalServerError, resp)
	}

	return ec.JSON(resp.Code, resp)
}
