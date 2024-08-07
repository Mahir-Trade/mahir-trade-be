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
	ModuleCtrl interface {
		CreateModule(ec echo.Context) error
		GetModuleByID(ec echo.Context) error
		GetModules(ec echo.Context) error
		GetModulesByGroupID(ec echo.Context) error
		UpdateModule(ec echo.Context) error
		DeleteModule(ec echo.Context) error
		GetPercetangeMarkWatchedModulesUser(ec echo.Context) error
	}

	ModuleCtrlImpl struct {
		dig.In

		ModuleSvc service.ModuleSvc
	}
)

func NewModuleCtrl(impl ModuleCtrlImpl) ModuleCtrl {
	return &impl
}

func (ox *ModuleCtrlImpl) CreateModule(ec echo.Context) error {
	ctx := ec.Request().Context()

	defer func() {
		if r := recover(); r != nil {
			slog.Error("CreateGroup - something went wrong", r)
		}
	}()

	var module service.ModuleRequest

	if err := ec.Bind(&module); err != nil {
		slog.Error("CreateModule - error binding request body", err)
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	validate := utils.Validate

	err := validate.Struct(module)
	if err != nil {
		slog.Error("CreateModule - validation error", err)
		errors := err.(validator.ValidationErrors)
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   errors.Error(),
		})
	}

	resp, err := ox.ModuleSvc.CreateModule(ctx, module)
	if err != nil {
		slog.Error("CreateModule - something went wrong", err)
		return ec.JSON(http.StatusBadRequest, resp)
	}
	return ec.JSON(http.StatusCreated, resp)
}

func (ox *ModuleCtrlImpl) GetModuleByID(ec echo.Context) error {
	ctx := ec.Request().Context()

	id := ec.Param("module_id")
	if id == "" {
		slog.Error("GetModuleByID - invalid request body")
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
		})
	}

	moduleID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		slog.Error("GetModuleByID - parsing error", err)
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   "invalid module id",
		})
	}

	if moduleID == 0 {
		slog.Error("GetModuleByID - module id is required")
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   "module id is required",
		})
	}

	resp, err := ox.ModuleSvc.GetModuleByID(ctx, moduleID)
	if err != nil {
		slog.Error("GetModuleByID - something went wrong", err)
		return ec.JSON(http.StatusBadRequest, resp)
	}
	return ec.JSON(http.StatusOK, resp)
}

func (ox *ModuleCtrlImpl) GetModules(ec echo.Context) error {
	ctx := ec.Request().Context()

	defer func() {
		if r := recover(); r != nil {
			slog.Error("CreateGroup - something went wrong", r)
		}
	}()

	var req models.PaginationRequest

	if err := ec.Bind(&req); err != nil {
		slog.Error("GetModules - something went wrong, bind error", err)
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   "invalid pagination request",
		})
	}

	validate := utils.Validate
	err := validate.Struct(req)
	if err != nil {
		slog.Error("GetModules - something went wrong, validation error", err)
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   "invalid pagination request",
		})
	}

	resp, err := ox.ModuleSvc.GetModules(ctx, req)
	if err != nil {
		slog.Error("GetModules - something went wrong, service error", err)
		return ec.JSON(http.StatusBadRequest, resp)
	}

	return ec.JSON(http.StatusOK, resp)

}

func (ox *ModuleCtrlImpl) GetModulesByGroupID(ec echo.Context) error {
	ctx := ec.Request().Context()

	defer func() {
		if r := recover(); r != nil {
			slog.Error("CreateGroup - something went wrong", r)
		}
	}()

	id := ec.Param("group_id")
	if id == "" {
		slog.Error("GetModulesByGroupID - invalid request body")
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   "group id is required",
		})
	}

	groupID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		slog.Error("GetModulesByGroupID - parsing error", err)
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   "group id is required",
		})
	}

	if groupID == 0 {
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   "group id is required",
		})
	}

	resp, err := ox.ModuleSvc.GetModulesByGroupID(ctx, groupID)
	if err != nil {
		slog.Error("GetModulesByGroupID - something went wrong", err)
		return ec.JSON(http.StatusBadRequest, resp)
	}

	return ec.JSON(http.StatusOK, resp)
}

func (ox *ModuleCtrlImpl) UpdateModule(ec echo.Context) error {
	ctx := ec.Request().Context()

	defer func() {
		if r := recover(); r != nil {
			slog.Error("CreateGroup - something went wrong", r)
		}
	}()

	id := ec.Param("module_id")
	if id == "" {
		slog.Error("UpdateModule - invalid request body")
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   "module id is required",
		})
	}

	moduleID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		slog.Error("UpdateModule - parsing error", err)
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   "module id is required",
		})
	}

	if moduleID == 0 {
		slog.Error("UpdateModule - module id is required")
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   "module id is required",
		})
	}

	var module service.ModuleRequest

	if err := ec.Bind(&module); err != nil {
		slog.Error("UpdateModule - error binding request body", err)
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	validate := utils.Validate
	err = validate.Struct(module)
	if err != nil {
		slog.Error("UpdateModule - validation error", err)
		errors := err.(validator.ValidationErrors)
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   errors.Error(),
		})
	}

	resp, err := ox.ModuleSvc.UpdateModule(ctx, moduleID, module)
	if err != nil {
		slog.Error("UpdateModule - something went wrong", err)
		return ec.JSON(http.StatusBadRequest, resp)
	}

	return ec.JSON(http.StatusOK, resp)

}

func (ox *ModuleCtrlImpl) DeleteModule(ec echo.Context) error {
	ctx := ec.Request().Context()

	defer func() {
		if r := recover(); r != nil {
			slog.Error("CreateGroup - something went wrong", r)
		}
	}()

	id := ec.Param("module_id")
	if id == "" {
		slog.Error("DeleteModule - invalid request body")
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   "module id is required",
		})
	}

	moduleID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		slog.Error("DeleteModule - parsing error", err)
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   "module id is required",
		})
	}

	if moduleID == 0 {
		slog.Error("DeleteModule - module id is required")
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   "module id is required",
		})
	}

	resp, err := ox.ModuleSvc.DeleteModule(ctx, moduleID)
	if err != nil {
		slog.Error("DeleteModule - something went wrong", err)
		return ec.JSON(http.StatusBadRequest, resp)
	}

	return ec.JSON(http.StatusOK, resp)
}

func (ox *ModuleCtrlImpl) GetPercetangeMarkWatchedModulesUser(ec echo.Context) error {
	ctx := ec.Request().Context()

	defer func() {
		if r := recover(); r != nil {
			slog.Error("GetPercetangeMarkWatchedModulesUser - something went wrong", r)
		}
	}()
	id := ec.Param("module_id")
	if id == "" {
		slog.Error("DeleteModule - invalid request body")
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   "module id is required",
		})
	}
	moduleID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		slog.Error("DeleteModule - parsing error", err)
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   "module id is required",
		})
	}

	if moduleID == 0 {
		slog.Error("DeleteModule - module id is required")
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   "module id is required",
		})
	}

	resp, err := ox.ModuleSvc.GetPercentageMarkWatchedModulesUser(ctx, moduleID)
	if err != nil {
		slog.Error("GetPercetangeMarkWatchedModulesUser - something went wrong", err)
		return ec.JSON(http.StatusBadRequest, resp)
	}

	return ec.JSON(http.StatusOK, resp)
}
