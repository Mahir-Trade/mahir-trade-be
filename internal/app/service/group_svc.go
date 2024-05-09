package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"mahir-trade-be/internal/app/models"
	"mahir-trade-be/internal/app/repo/postgres"
	"math"
	"net/http"

	"go.uber.org/dig"
)

type (
	GroupSvc interface {
		CreateGroup(ctx context.Context, req models.Group) (resp models.DefaultResponse, err error)
		GetGroupByID(ctx context.Context, groupId int64) (resp models.DefaultResponse, err error)
		GetGroups(ctx context.Context, req models.GetGroupsRequest) (resp models.DefaultPaginationResponseData, err error)
		UpdateGroup(ctx context.Context, req models.Group) (resp models.DefaultResponse, err error)
		DeleteGroup(ctx context.Context, groupId int64) (resp models.DefaultResponse, err error)
	}

	GroupSvcImpl struct {
		dig.In

		GroupRepo postgres.GroupRepo
	}
)

func NewGroupSvc(impl GroupSvcImpl) GroupSvc {
	return &impl
}

func (g *GroupSvcImpl) CreateGroup(ctx context.Context, req models.Group) (resp models.DefaultResponse, err error) {
	{
		resp.Code = http.StatusCreated
		resp.Message = "Success"
	}

	groupId, err := g.GroupRepo.CreateGroup(ctx, req)
	if err != nil {
		resp.Code = http.StatusBadRequest
		resp.Message = "bad request"
		resp.Error = err.Error()
		slog.ErrorContext(ctx, fmt.Sprintf("[service][CreateGroup] while CreateGroup err : %v", err.Error()))

		return resp, err
	}

	type Group struct {
		ID int `json:"id"`
	}

	resp.Data = Group{ID: groupId}

	return
}

func (g *GroupSvcImpl) GetGroupByID(ctx context.Context, groupId int64) (resp models.DefaultResponse, err error) {
	{
		resp.Code = http.StatusOK
		resp.Message = "Success"
	}

	group, err := g.GroupRepo.GetGroupByID(ctx, groupId)
	if err != nil {
		resp.Code = http.StatusBadRequest
		resp.Message = "bad request"
		resp.Error = err.Error()
		slog.ErrorContext(ctx, fmt.Sprintf("[service][GetGroupByID] while GetGroupByID err : %v", err))

		return resp, err
	}

	resp.Data = group

	return
}

func (g *GroupSvcImpl) GetGroups(ctx context.Context, req models.GetGroupsRequest) (resp models.DefaultPaginationResponseData, err error) {
	var dataResp models.DefaultResponse
	{
		dataResp.Code = http.StatusOK
		dataResp.Message = "Success"
	}

	groups, totalData, err := g.GroupRepo.GetGroups(ctx, req)
	if err != nil {
		dataResp.Code = http.StatusBadRequest
		dataResp.Message = "bad request"
		dataResp.Error = err.Error()
		resp.Results = dataResp
		slog.ErrorContext(ctx, fmt.Sprintf("[service][GetGroups] while GetGroups err : %v", err))

		return resp, err
	}

	// mapping response
	{
		dataResp.Data = groups
		resp.Page = uint(req.Page)
		resp.Limit = uint(req.Limit)

		totalPage := math.Ceil(float64(totalData) / float64(req.Limit))
		resp.TotalPages = uint(totalPage)

		resp.TotalItems = uint(totalData)
		resp.HasNext = req.Page < int64(resp.TotalPages)
		resp.HasPrevious = req.Page > 1
		resp.Results = dataResp
	}

	return
}

func (g *GroupSvcImpl) UpdateGroup(ctx context.Context, req models.Group) (resp models.DefaultResponse, err error) {
	{
		resp.Code = http.StatusOK
		resp.Message = "Success"
	}

	group, err := g.GroupRepo.GetGroupByID(ctx, req.ID)
	if err != nil {
		resp.Code = http.StatusBadRequest
		resp.Message = "bad request"
		resp.Error = err.Error()
		slog.ErrorContext(ctx, fmt.Sprintf("[service][UpdateGroup] while GetGroupByID err : %v", err))

		return resp, err
	}

	if group.ID == 0 {
		errMsg := "group not found"

		resp.Code = http.StatusBadRequest
		resp.Message = errMsg
		resp.Error = errors.New(errMsg).Error()
		slog.ErrorContext(ctx, fmt.Sprintf("[service][UpdateGroup] err : %v", errMsg))

		return resp, errors.New(errMsg)
	}

	if req.GroupName == group.GroupName {
		errMsg := "group name is same as before"

		resp.Code = http.StatusBadRequest
		resp.Message = errMsg
		resp.Error = errors.New(errMsg).Error()
		slog.ErrorContext(ctx, fmt.Sprintf("[service][UpdateGroup] err : %v", errMsg))

		return resp, errors.New(errMsg)
	}

	err = g.GroupRepo.UpdateGroup(ctx, req)
	if err != nil {
		resp.Code = http.StatusBadRequest
		resp.Message = "bad request"
		resp.Error = err.Error()
		slog.ErrorContext(ctx, fmt.Sprintf("[service][UpdateGroup] while UpdateGroup err : %v", err))

		return resp, err
	}

	return
}

func (g *GroupSvcImpl) DeleteGroup(ctx context.Context, groupId int64) (resp models.DefaultResponse, err error) {
	{
		resp.Code = http.StatusOK
		resp.Message = "Success"
	}

	// if groupId == 0 {
	// 	errMsg := errors.New("group id is required")
	// 	slog.ErrorContext(ctx, fmt.Sprintf("[service][DeleteGroup] err : %v", errMsg.Error()))
	// 	return errMsg
	// }

	group, err := g.GroupRepo.GetGroupByID(ctx, groupId)
	if err != nil {
		resp.Code = http.StatusBadRequest
		resp.Message = "bad request"
		resp.Error = err
		slog.ErrorContext(ctx, fmt.Sprintf("[service][DeleteGroup] while GetGroupByID err : %v", err))

		return resp, err
	}

	if group.ID == 0 {
		errMsg := "group not found"

		resp.Code = http.StatusBadRequest
		resp.Message = errMsg
		resp.Error = err
		slog.ErrorContext(ctx, fmt.Sprintf("[service][DeleteGroup] err : %v", errMsg))

		return resp, errors.New(errMsg)
	}

	err = g.GroupRepo.SoftDeleteGroup(ctx, groupId, "SYSTEM")
	if err != nil {
		resp.Code = http.StatusBadRequest
		resp.Message = "bad request"
		resp.Error = err
		slog.ErrorContext(ctx, fmt.Sprintf("[service][DeleteGroup] while DeleteGroup err : %v", err))

		return resp, err
	}

	return
}
