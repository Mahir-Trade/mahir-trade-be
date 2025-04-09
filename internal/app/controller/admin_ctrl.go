package controller

import (
	"fmt"
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
	AdminCtrl interface {
		AdminLogin(ec echo.Context) error
		AdminRegistration(ec echo.Context) error
		UpdateTypeUser(ec echo.Context) error
		GetDetailUserForBO(ec echo.Context) error
		GetDetailAdminInfo(ec echo.Context) error
		GetAllUsers(ec echo.Context) error
		StartMembershipProgram(ec echo.Context) error
	}

	AdminCtrlImpl struct {
		dig.In

		AdminSvc service.AdminSvc
		UserSvc  service.UserSvc
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
		return ec.JSON(resp.Code, resp)
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
		return ec.JSON(resp.Code, resp)
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

func (ox *AdminCtrlImpl) GetDetailUserForBO(ec echo.Context) error {
	ctx := ec.Request().Context()

	defer func() {
		if r := recover(); r != nil {
			slog.Error("GetDetailUserForBO - something went wrong", r)
		}
	}()

	userID := ec.Param("user_id")
	userIDInt, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		slog.Error("GetDetailUserForBO - something went wrong", err)
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "bad request",
			Error:   err.Error(),
		})
	}

	res, err := ox.UserSvc.GetDetailUserForBO(ctx, userIDInt)
	if err != nil {
		slog.Error("GetDetailUserForBO - something went wrong", err)
		return ec.JSON(http.StatusBadRequest, res)
	}

	return ec.JSON(http.StatusOK, res)
}

func (ox *AdminCtrlImpl) GetDetailAdminInfo(ec echo.Context) error {
	ctx := ec.Request().Context()

	defer func() {
		if r := recover(); r != nil {
			slog.Error("GetDetailAdminInfo - something went wrong", r)
		}
	}()

	userData, ok := ctx.Value(middleware.UserData).(middleware.UserCtxReq)
	if !ok {
		slog.Error("GetDetailAdminInfo - something went wrong")
		return ec.JSON(http.StatusInternalServerError, models.DefaultResponse{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		})
	}

	res, err := ox.AdminSvc.GetDetailAdminInfo(ctx, userData.Username)
	if err != nil {
		slog.Error("GetDetailAdminInfo - something went wrong", err)
		return ec.JSON(http.StatusInternalServerError, res)
	}

	return ec.JSON(http.StatusOK, res)
}

func (ox *AdminCtrlImpl) GetAllUsers(ec echo.Context) error {
	ctx := ec.Request().Context()

	defer func() {
		if r := recover(); r != nil {
			slog.Error("GetAllUsers - something went wrong", r)
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
	MembershipStatus := ec.QueryParam("membershipStatus")
	req := models.PaginationRequest{
		Limit:            limit,
		Page:             page,
		Search:           search,
		MembershipStatus: MembershipStatus,
	}

	res, err := ox.AdminSvc.GetAllUsers(ctx, req)
	if err != nil {
		slog.Error("GetAllUsers - something went wrong", err)
		return ec.JSON(http.StatusBadRequest, res)
	}

	return ec.JSON(http.StatusOK, res)
}

func (ox *AdminCtrlImpl) StartMembershipProgram(ec echo.Context) error {
	ctx := ec.Request().Context()

	defer func() {
		if r := recover(); r != nil {
			slog.Error("StartMembershipProgram - something went wrong", r)
		}
	}()

	var req models.StartMembershipProgramRequest

	if err := ec.Bind(&req); err != nil {
		slog.Error("StartMembershipProgram - something went wrong", err)
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	validate := utils.Validate

	err := validate.Struct(req)
	if err != nil {
		slog.Error("StartMembershipProgram - something went wrong", err)
		errs := err.(validator.ValidationErrors)
		var errorMessages []string
		for _, e := range errs {
			errorMessages = append(errorMessages, fmt.Sprintf("Field %s is %s", e.Field(), e.Tag()))
		}
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   errorMessages,
		})
	}

	resp, err := ox.AdminSvc.StartMembershipProgram(ctx, req)
	if err != nil {
		slog.Error("StartMembershipProgram - something went wrong", err)
		return ec.JSON(resp.Code, resp)
	}

	return ec.JSON(resp.Code, resp)
}
