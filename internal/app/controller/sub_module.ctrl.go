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
	SubModuleCtrl interface {
		CreateSubModule(ec echo.Context) error
		GetSubModules(ec echo.Context) error
		GetSubModuleByID(ec echo.Context) error
		GetSubModulesByModuleID(ec echo.Context) error
		UpdateSubModule(ec echo.Context) error
		SoftDeleteSubModule(ec echo.Context) error
		MarkSubModuleAsWatched(ec echo.Context) error
	}

	SubModuleCtrlImpl struct {
		dig.In

		SubModuleSvc service.SubModuleSvc
	}
)

func NewSubModuleCtrl(impl SubModuleCtrlImpl) SubModuleCtrl {
	return &impl
}

func (ox *SubModuleCtrlImpl) CreateSubModule(ec echo.Context) error {
	ctx := ec.Request().Context()

	defer func() {
		if r := recover(); r != nil {
			slog.Error("CreateGroup - something went wrong", r)
		}
	}()

	var req service.SubModuleRequest

	if err := ec.Bind(&req); err != nil {
		slog.Error("CreateSubModule - error binding request body", err)
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	validate := utils.Validate

	err := validate.Struct(req)
	if err != nil {
		slog.Error("CreateSubModule - validation error", err)
		errors := err.(validator.ValidationErrors)
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   errors.Error(),
		})
	}

	res, err := ox.SubModuleSvc.CreateSubModule(ctx, req)
	if err != nil {
		slog.Error("CreateSubModule - something went wrong", err)
		return ec.JSON(res.Code, res)
	}

	return ec.JSON(res.Code, res)

}

func (ox *SubModuleCtrlImpl) GetSubModules(ec echo.Context) error {
	ctx := ec.Request().Context()

	defer func() {
		if r := recover(); r != nil {
			slog.Error("GetSubModules - something went wrong", r)
		}
	}()

	var req models.PaginationRequest

	if err := ec.Bind(&req); err != nil {
		slog.Error("GetSubModules - error binding request body", err)
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	validate := utils.Validate
	err := validate.Struct(req)
	if err != nil {
		slog.Error("GetSubModules - validation error", err)
		errors := err.(validator.ValidationErrors)
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   errors.Error(),
		})
	}

	res, err := ox.SubModuleSvc.GetSubModules(ctx, req)
	if err != nil {
		slog.Error("GetSubModules - something went wrong", err)
		return ec.JSON(http.StatusInternalServerError, res)
	}

	return ec.JSON(http.StatusOK, res)
}

func (ox *SubModuleCtrlImpl) GetSubModuleByID(ec echo.Context) error {
	ctx := ec.Request().Context()

	defer func() {
		if r := recover(); r != nil {
			slog.Error("GetSubModuleByID - something went wrong", r)
		}
	}()
	id := ec.Param("sub_module_id")
	if id == "" {
		slog.Error("GetSubModuleByID - invalid request body")
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
		})
	}

	subModuleID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		slog.Error("GetSubModuleByID - parsing error", err)
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   "invalid module id",
		})
	}

	if subModuleID == 0 {
		slog.Error("GetSubModuleByID - module id is required")
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   "module id is required",
		})
	}

	res, err := ox.SubModuleSvc.GetSubModuleByID(ctx, subModuleID)
	if err != nil {
		slog.Error("GetSubModuleByID - something went wrong", err)
		return ec.JSON(http.StatusInternalServerError, res)
	}

	return ec.JSON(http.StatusOK, res)
}

func (ox *SubModuleCtrlImpl) GetSubModulesByModuleID(ec echo.Context) error {
	ctx := ec.Request().Context()

	defer func() {
		if r := recover(); r != nil {
			slog.Error("GetSubModuleByModuleID - something went wrong", r)
		}
	}()

	id := ec.Param("module_id")
	if id == "" {
		slog.Error("GetSubModuleByModuleID - invalid request body")
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   "module id is required",
		})
	}

	moduleID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		slog.Error("GetSubModuleByModuleID - parsing error", err)
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   "module id is required",
		})
	}

	if moduleID == 0 {
		slog.Error("GetSubModuleByModuleID - module id is required")
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   "module id is required",
		})
	}

	res, err := ox.SubModuleSvc.GetSubModulesByModuleID(ctx, moduleID)
	if err != nil {
		slog.Error("GetSubModuleByModuleID - something went wrong", err)
		return ec.JSON(res.Code, res)
	}

	return ec.JSON(http.StatusOK, res)
}

func (ox *SubModuleCtrlImpl) UpdateSubModule(ec echo.Context) error {
	ctx := ec.Request().Context()

	defer func() {
		if r := recover(); r != nil {
			slog.Error("UpdateSubModule - something went wrong", r)
		}
	}()

	id := ec.Param("sub_module_id")
	if id == "" {
		slog.Error("UpdateSubModule - invalid request body")
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   "module id is required",
		})
	}

	subModuleID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		slog.Error("UpdateSubModule - parsing error", err)
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   "module id is required",
		})
	}

	var req service.SubModuleRequest

	if err := ec.Bind(&req); err != nil {
		slog.Error("UpdateSubModule - error binding request body", err)
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	validate := utils.Validate
	err = validate.Struct(req)
	if err != nil {
		slog.Error("UpdateSubModule - validation error", err)
		errors := err.(validator.ValidationErrors)
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   errors.Error(),
		})
	}

	res, err := ox.SubModuleSvc.UpdateSubModule(ctx, subModuleID, req)
	if err != nil {
		slog.Error("UpdateSubModule - something went wrong", err)
		return ec.JSON(res.Code, res)
	}

	return ec.JSON(res.Code, res)
}

func (ox *SubModuleCtrlImpl) SoftDeleteSubModule(ec echo.Context) error {
	ctx := ec.Request().Context()

	defer func() {
		if r := recover(); r != nil {
			slog.Error("SoftDeleteSubModule - something went wrong", r)
		}
	}()

	id := ec.Param("sub_module_id")
	if id == "" {
		slog.Error("SoftDeleteSubModule - invalid request body")
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   "module id is required",
		})
	}

	subModuleID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		slog.Error("SoftDeleteSubModule - parsing error", err)
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   "module id is required",
		})
	}

	if subModuleID == 0 {
		slog.Error("SoftDeleteSubModule - module id is required")
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   "module id is required",
		})
	}

	res, err := ox.SubModuleSvc.SoftDeleteSubModule(ctx, subModuleID)
	if err != nil {
		slog.Error("SoftDeleteSubModule - something went wrong", err)
		return ec.JSON(http.StatusInternalServerError, res)
	}

	return ec.JSON(http.StatusOK, res)
}

func (ox *SubModuleCtrlImpl) MarkSubModuleAsWatched(ec echo.Context) error {
	ctx := ec.Request().Context()

	defer func() {
		if r := recover(); r != nil {
			slog.Error("MarkSubModuleAsWatched - something went wrong", r)
		}
	}()

	var req service.MarkSubModuleAsWatchedRequest

	if err := ec.Bind(&req); err != nil {
		slog.Error("MarkSubModuleAsWatched - error binding request body", err)
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	validate := utils.Validate
	err := validate.Struct(req)
	if err != nil {
		slog.Error("MarkSubModuleAsWatched - validation error", err)
		errors := err.(validator.ValidationErrors)
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   errors.Error(),
		})
	}

	res, err := ox.SubModuleSvc.MarkSubModuleAsWatched(ctx, req)
	if err != nil {
		slog.Error("MarkSubModuleAsWatched - something went wrong", err)
		return ec.JSON(res.Code, res)
	}

	return ec.JSON(res.Code, res)
}
