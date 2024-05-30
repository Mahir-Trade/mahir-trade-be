package controller

import (
	"errors"
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
	GroupCtrl interface {
		CreateGroup(ec echo.Context) error
		GetGroupByID(ec echo.Context) error
		GetGroups(ec echo.Context) error
		UpdateGroup(ec echo.Context) error
		DeleteGroup(ec echo.Context) error
	}

	GroupCtrlImpl struct {
		dig.In

		GroupSvc service.GroupSvc
	}
)

func NewGroupCtrl(impl GroupCtrlImpl) GroupCtrl {
	return &impl
}

func (ox *GroupCtrlImpl) CreateGroup(ec echo.Context) error {
	ctx := ec.Request().Context()

	defer func() {
		if r := recover(); r != nil {
			slog.Error("CreateGroup - something went wrong", r)
		}
	}()

	var group models.Group

	if err := ec.Bind(&group); err != nil {
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	validate := utils.Validate

	err := validate.Struct(group)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   errors.Error(),
		})
	}

	resp, err := ox.GroupSvc.CreateGroup(ctx, group)
	if err != nil {
		slog.Error("CreateGroup - something went wrong", err)
		return ec.JSON(http.StatusBadRequest, resp)
	}

	return ec.JSON(resp.Code, resp)
}

func (ox *GroupCtrlImpl) GetGroupByID(ec echo.Context) error {
	ctx := ec.Request().Context()

	defer func() {
		if r := recover(); r != nil {
			slog.Error("GetGroupByID - something went wrong", r)
		}
	}()

	id := ec.Param("id")
	if id == "" {
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "group id is required",
		})
	}

	groupID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		slog.Error("GetGroupByID - something went wrong", err)
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "group id is required",
		})
	}

	if groupID == 0 {
		slog.Error("GetGroupByID - something went wrong", errors.New("group id is required"))
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "group id is required",
		})
	}

	resp, err := ox.GroupSvc.GetGroupByID(ctx, groupID)
	if err != nil {
		slog.Error("GetGroupByID - something went wrong", err)
		return ec.JSON(http.StatusBadRequest, resp)
	}

	return ec.JSON(resp.Code, resp)
}

func (ox *GroupCtrlImpl) GetGroups(ec echo.Context) error {
	ctx := ec.Request().Context()

	defer func() {
		if r := recover(); r != nil {
			slog.Error("GetGroups - something went wrong", r)
		}
	}()

	var req models.GetGroupsRequest
	if err := ec.Bind(&req); err != nil {
		slog.Error("GetGroups - something went wrong, bind error", err)
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request",
			Error:   "invalid pagination request",
		})
	}

	resp, err := ox.GroupSvc.GetGroups(ctx, req)
	if err != nil {
		slog.Error("GetGroups - something went wrong", err)
		return ec.JSON(http.StatusBadRequest, resp)
	}

	return ec.JSON(http.StatusOK, resp)
}

func (ox *GroupCtrlImpl) UpdateGroup(ec echo.Context) error {
	var (
		errMsg string
		ctx    = ec.Request().Context()
	)

	defer func() {
		if r := recover(); r != nil {
			slog.Error("UpdateGroup - something went wrong", r)
		}
	}()

	id := ec.Param("id")
	if id == "" {
		errMsg = "group id is required"
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: errMsg,
			Error:   errors.New(errMsg),
		})
	}

	groupID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		slog.Error("GetGroupByID - something went wrong", err)
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "group id is required",
			Error:   err.Error(),
		})
	}

	if groupID == 0 {
		errMsg = "group id is required"
		slog.Error("GetGroupByID - something went wrong", errors.New("group id is required"))
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: errMsg,
			Error:   errors.New(errMsg),
		})
	}

	var group models.Group

	if err := ec.Bind(&group); err != nil {
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	validate := utils.Validate

	err = validate.Struct(group)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   errors.Error(),
		})
	}

	group.ID = groupID

	resp, err := ox.GroupSvc.UpdateGroup(ctx, group)
	if err != nil {
		slog.Error("UpdateGroup - something went wrong", err)
		return ec.JSON(http.StatusBadRequest, resp)
	}

	return ec.JSON(resp.Code, resp)
}

func (ox *GroupCtrlImpl) DeleteGroup(ec echo.Context) error {
	var (
		errMsg string
		ctx    = ec.Request().Context()
	)

	defer func() {
		if r := recover(); r != nil {
			slog.Error("DeleteGroup - something went wrong", r)
		}
	}()

	id := ec.Param("id")
	if id == "" {
		errMsg = "group id is required"
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: errMsg,
			Error:   errors.New(errMsg),
		})
	}

	groupID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		slog.Error("DeleteGroup - something went wrong", err)
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "group id is required",
			Error:   err.Error(),
		})
	}

	if groupID == 0 {
		errMsg = "group id is required"
		slog.Error("DeleteGroup - something went wrong", errors.New("group id is required"))
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: errMsg,
			Error:   errors.New(errMsg),
		})
	}

	resp, err := ox.GroupSvc.DeleteGroup(ctx, groupID)
	if err != nil {
		slog.Error("DeleteGroup - something went wrong", err)
		return ec.JSON(http.StatusBadRequest, resp)
	}

	return ec.JSON(resp.Code, resp)
}
