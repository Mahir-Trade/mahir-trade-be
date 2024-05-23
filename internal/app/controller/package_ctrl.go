package controller

import (
	"errors"
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
	PackageCtrl interface {
		CreatePackage(ec echo.Context) error
		GetPackages(ec echo.Context) error
		GetPackageByID(ec echo.Context) error
		UpdatePackage(ec echo.Context) error
		DeletePackage(ec echo.Context) error
	}

	PackageCtrlImpl struct {
		dig.In

		PackageSvc service.PackageSvc
	}
)

func NewPackageCtrl(impl PackageCtrlImpl) PackageCtrl {
	return &impl
}

func (ox *PackageCtrlImpl) CreatePackage(ec echo.Context) error {
	ctx := ec.Request().Context()

	defer func() {
		if r := recover(); r != nil {
			slog.Error("CreateGroup - something went wrong", r)
		}
	}()

	userData, ok := ctx.Value(middleware.UserData).(middleware.UserCtxReq)
	if !ok {
		slog.Error("CreateGroup - something went wrong", errors.New("user data not found"))
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "user data not found",
		})
	}

	var pkg models.Package
	if err := ec.Bind(&pkg); err != nil {
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	validate := utils.Validate
	err := validate.Struct(pkg)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   errors.Error(),
		})
	}

	pkg.CreatedBy = userData.Email

	resp, err := ox.PackageSvc.CreatePackage(ctx, pkg)
	if err != nil {
		return ec.JSON(http.StatusBadRequest, resp)
	}

	return ec.JSON(http.StatusCreated, resp)
}

func (ox *PackageCtrlImpl) GetPackages(ec echo.Context) error {
	ctx := ec.Request().Context()

	defer func() {
		if r := recover(); r != nil {
			slog.Error("GetPackages - something went wrong", r)
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

	req := models.PaginationRequest{
		Limit: limit,
		Page:  page,
	}

	resp, err := ox.PackageSvc.GetPackages(ctx, req)
	if err != nil {
		slog.Error("GetPackages - something went wrong", err)
		return ec.JSON(http.StatusBadRequest, resp)
	}

	return ec.JSON(http.StatusOK, resp)
}

func (ox *PackageCtrlImpl) GetPackageByID(ec echo.Context) error {
	ctx := ec.Request().Context()

	defer func() {
		if r := recover(); r != nil {
			slog.Error("GetPackageByID - something went wrong", r)
		}
	}()

	id := ec.Param("id")
	if id == "" {
		slog.Error("GetPackageByID - package id is required")
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "package id is required",
		})
	}

	pkgID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		slog.Error("GetPackageByID - invalid package id", err)
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "invalid package id",
		})
	}

	if pkgID == 0 {
		slog.Error("GetPackageByID - something went wrong", errors.New("package id is required"))
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "package id is required",
		})
	}

	resp, err := ox.PackageSvc.GetPackageByID(ctx, pkgID)
	if err != nil {
		slog.Error("GetPackageByID - something went wrong", err)
		return ec.JSON(http.StatusBadRequest, resp)
	}

	return ec.JSON(http.StatusOK, resp)
}

func (ox *PackageCtrlImpl) UpdatePackage(ec echo.Context) error {
	var (
		errMsg string
		ctx    = ec.Request().Context()
	)

	defer func() {
		if r := recover(); r != nil {
			slog.Error("UpdatePackage - something went wrong", r)
		}
	}()

	userData, ok := ctx.Value(middleware.UserData).(middleware.UserCtxReq)
	if !ok {
		slog.Error("UpdatePackage - something went wrong", errors.New("user data not found"))
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "user data not found",
		})
	}

	id := ec.Param("id")
	if id == "" {
		errMsg = "package id is required"
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: errMsg,
			Error:   errors.New(errMsg),
		})
	}

	packageID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		slog.Error("UpdatePackage - something went wrong", err)
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "package id is required",
			Error:   err.Error(),
		})
	}

	if packageID == 0 {
		errMsg = "package id is required"
		slog.Error("UpdatePackage - something went wrong", errors.New("package id is required"))
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: errMsg,
			Error:   errors.New(errMsg),
		})
	}

	var pkg models.Package
	if err := ec.Bind(&pkg); err != nil {
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	pkg.ID = packageID
	validate := utils.Validate
	err = validate.Struct(pkg)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   errors.Error(),
		})
	}

	pkg.UpdatedBy = userData.Email

	resp, err := ox.PackageSvc.UpdatePackage(ctx, pkg)
	if err != nil {
		slog.Error("UpdatePackage - something went wrong", err)
		return ec.JSON(http.StatusBadRequest, resp)
	}

	return ec.JSON(http.StatusOK, resp)
}

func (ox *PackageCtrlImpl) DeletePackage(ec echo.Context) error {
	var (
		errMsg string
		ctx    = ec.Request().Context()
	)

	defer func() {
		if r := recover(); r != nil {
			slog.Error("DeletePackage - something went wrong", r)
		}
	}()

	userData, ok := ctx.Value(middleware.UserData).(middleware.UserCtxReq)
	if !ok {
		slog.Error("DeletePackage - something went wrong", errors.New("user data not found"))
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "user data not found",
		})
	}

	id := ec.Param("id")
	if id == "" {
		errMsg = "package id is required"
		slog.Error("DeletePackage - something went wrong", errors.New("package id is required"))
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: errMsg,
			Error:   errors.New(errMsg),
		})
	}

	packageID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		slog.Error("DeletePackage - something went wrong", err)
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "package id is required",
			Error:   err.Error(),
		})
	}

	if packageID == 0 {
		errMsg = "package id is required"
		slog.Error("DeletePackage - something went wrong", errors.New("package id is required"))
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: errMsg,
			Error:   errors.New(errMsg),
		})
	}

	resp, err := ox.PackageSvc.DeletePackage(ctx, packageID, userData.Email)
	if err != nil {
		slog.Error("DeletePackage - something went wrong", err)
		return ec.JSON(http.StatusBadRequest, resp)
	}

	return ec.JSON(http.StatusOK, resp)
}
