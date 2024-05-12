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
		UpdateModule(ec echo.Context) error
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
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	validate := utils.Validate

	err := validate.Struct(module)
	if err != nil {
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
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
		})
	}

	moduleID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
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
	return nil
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
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   "module id is required",
		})
	}

	moduleID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   "module id is required",
		})
	}

	if moduleID == 0 {
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   "module id is required",
		})
	}

	var module service.ModuleRequest

	if err := ec.Bind(&module); err != nil {
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	validate := utils.Validate
	err = validate.Struct(module)
	if err != nil {
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
