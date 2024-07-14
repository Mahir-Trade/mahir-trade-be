package service

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"mahir-trade-be/internal/app/models"
	"mahir-trade-be/internal/app/repo/postgres"
	"mahir-trade-be/pkg/middleware"
	"math"
	"net/http"

	"go.uber.org/dig"
)

type (
	ModuleRequest struct {
		ModuleName   string `json:"module_name" validate:"required"`
		GroupID      int64  `json:"group_id,omitempty"`
		ThumbnailUrl string `json:"thumbnail_url,omitempty"`
		Tag          string `json:"tag,omitempty"`
		CreatedBy    string `json:"created_by,omitempty"`
		UpdatedBy    string `json:"updated_by,omitempty"`
	}

	ModuleResponse struct {
		ID           int64  `json:"id"`
		UUID         string `json:"uuid"`
		GroupID      int64  `json:"group_id"`
		GroupName    string `json:"group_name"`
		ModuleName   string `json:"module_name"`
		ThumbnailUrl string `json:"thumbnail_url"`
		Tag          string `json:"tag"`
		CreatedBy    string `json:"created_by"`
		UpdatedBy    string `json:"updated_by"`
		CreatedAt    string `json:"created_at,omitempty"`
		UpdatedAt    string `json:"updated_at,omitempty"`
	}

	ModuleSvc interface {
		CreateModule(ctx context.Context, req ModuleRequest) (resp models.DefaultResponse, err error)
		GetModuleByID(ctx context.Context, moduleID int64) (resp models.DefaultResponse, err error)
		GetModules(ctx context.Context, req models.PaginationRequest) (resp models.DefaultPaginationResponseData, err error)
		GetModulesByGroupID(ctx context.Context, groupID int64) (resp models.DefaultResponse, err error)
		UpdateModule(ctx context.Context, moduleID int64, req ModuleRequest) (resp models.DefaultResponse, err error)
		DeleteModule(ctx context.Context, moduleID int64) (resp models.DefaultResponse, err error)
	}

	ModuleSvcImpl struct {
		dig.In

		ModuleRepo    postgres.ModuleRepo
		GroupRepo     postgres.GroupRepo
		SubModuleRepo postgres.SubModuleRepo
	}
)

func NewModuleSvc(impl ModuleSvcImpl) ModuleSvc {
	return &impl
}

func (m *ModuleSvcImpl) CreateModule(ctx context.Context, req ModuleRequest) (resp models.DefaultResponse, err error) {
	{
		resp.Code = http.StatusCreated
		resp.Message = "Success"
		resp.Data = struct{}{}
	}

	adminData, ok := ctx.Value(middleware.UserData).(middleware.UserCtxReq)
	if !ok {
		resp.Code = http.StatusBadRequest
		resp.Message = "bad request"
		resp.Error = errors.New("something went wrong, we will fix it soon").Error()
		return
	}

	model := models.Module{
		ModuleName: req.ModuleName,
	}

	{
		if req.GroupID != 0 {
			_, err = m.GroupRepo.GetGroupByID(ctx, req.GroupID)
			if err != nil {
				resp.Code = http.StatusBadRequest
				resp.Message = "bad request"
				resp.Error = err.Error()
				return
			}

			model.GroupID = sql.NullInt64{
				Int64: req.GroupID,
				Valid: true,
			}
		}

		if req.ThumbnailUrl != "" {
			model.ThumbnailUrl = sql.NullString{
				String: req.ThumbnailUrl,
				Valid:  true,
			}
		}

		if req.Tag != "" {
			model.Tag = sql.NullString{
				String: req.Tag,
				Valid:  true,
			}
		}

		model.CreatedBy = adminData.Username
	}

	moduleId, err := m.ModuleRepo.CreateModule(ctx, model)
	if err != nil {
		resp.Code = http.StatusBadRequest
		resp.Message = "bad request"
		resp.Error = err.Error()

		return resp, err
	}

	type Module struct {
		ID int `json:"id"`
	}

	resp.Data = Module{ID: moduleId}

	return
}

func (m *ModuleSvcImpl) GetModuleByID(ctx context.Context, moduleID int64) (resp models.DefaultResponse, err error) {

	{
		resp.Code = http.StatusOK
		resp.Message = "Success"
		resp.Data = struct{}{}
	}

	module, err := m.ModuleRepo.GetModuleByID(ctx, moduleID)
	if err != nil {
		slog.ErrorContext(ctx, "[service][GetModuleByID] error while GetModuleByID err", "%v", err.Error())
		resp.Code = http.StatusBadRequest
		resp.Message = "bad request"
		resp.Error = err.Error()

		return resp, err
	}

	respData := ModuleResponse{
		ID:         module.ID,
		UUID:       module.UUID,
		ModuleName: module.ModuleName,
		CreatedBy:  module.CreatedBy,
		UpdatedBy:  module.UpdatedBy,
		CreatedAt:  module.CreatedAt,
		UpdatedAt:  module.UpdatedAt,
	}
	if module.GroupID.Valid {
		respData.GroupID = module.GroupID.Int64
		group, errGroup := m.GroupRepo.GetGroupByID(ctx, module.GroupID.Int64)
		if errGroup != nil {
			slog.ErrorContext(ctx, "[service][GetModuleByID] error while GetGroupByID err", "%v", errGroup.Error())
			resp.Code = http.StatusBadRequest
			resp.Message = "bad request"
			resp.Error = errGroup.Error()
			return resp, errGroup
		}
		respData.GroupName = group.GroupName
	}

	if module.ThumbnailUrl.Valid {
		respData.ThumbnailUrl = module.ThumbnailUrl.String
	}

	if module.Tag.Valid {
		respData.Tag = module.Tag.String
	}

	resp.Data = respData

	return
}

func (m *ModuleSvcImpl) GetModules(ctx context.Context, req models.PaginationRequest) (resp models.DefaultPaginationResponseData, err error) {
	var dataResp models.DefaultResponse
	{
		dataResp.Code = http.StatusOK
		dataResp.Message = "Success"
		dataResp.Data = struct{}{}
	}

	modules, totalData, err := m.ModuleRepo.GetModules(ctx, req)
	if err != nil {
		slog.ErrorContext(ctx, "[service][GetModules] error while GetModules err", "%v", err.Error())
		dataResp.Code = http.StatusBadRequest
		dataResp.Message = "bad request"
		dataResp.Error = errors.New("something went wrong, we will fix it soon").Error()
		return
	}

	var respData []ModuleResponse
	for _, module := range modules {
		data := ModuleResponse{
			ID:         module.ID,
			UUID:       module.UUID,
			ModuleName: module.ModuleName,
			CreatedBy:  module.CreatedBy,
			UpdatedBy:  module.UpdatedBy,
			CreatedAt:  module.CreatedAt,
			UpdatedAt:  module.UpdatedAt,
		}
		if module.GroupID.Valid {
			group, errGroup := m.GroupRepo.GetGroupByID(ctx, module.GroupID.Int64)
			if errGroup != nil {
				slog.ErrorContext(ctx, "[service][GetModules] error while GetGroupByID err", "%v", errGroup.Error())
			}
			data.GroupID = module.GroupID.Int64
			data.GroupName = group.GroupName
		}
		if module.ThumbnailUrl.Valid {
			data.ThumbnailUrl = module.ThumbnailUrl.String
		}
		if module.Tag.Valid {
			data.Tag = module.Tag.String
		}
		respData = append(respData, data)
	}

	{
		dataResp.Data = respData
		resp.Page = uint(req.Page)
		resp.Limit = uint(req.Limit)

		totalPage := math.Ceil(float64(totalData) / float64(req.Limit))
		resp.TotalPages = uint(totalPage)

		resp.TotalItems = uint(len(respData))
		resp.HasNext = resp.Page < resp.TotalPages
		resp.HasPrevious = resp.Page > 1
		resp.Results = dataResp
	}

	return
}

func (m *ModuleSvcImpl) UpdateModule(ctx context.Context, moduleID int64, req ModuleRequest) (resp models.DefaultResponse, err error) {
	{
		resp.Code = http.StatusOK
		resp.Message = "Success"
		resp.Data = struct{}{}
	}

	module, err := m.ModuleRepo.GetModuleByID(ctx, moduleID)
	if err != nil {
		resp.Code = http.StatusBadRequest
		resp.Message = "bad request"
		resp.Error = err.Error()
		slog.ErrorContext(ctx, "[service][UpdateModule] while GetModuleByID err", "%v", err)
		return
	}

	if module.ID == 0 {
		errMsg := "module not found"
		resp.Code = http.StatusBadRequest
		resp.Message = errMsg
		resp.Error = errors.New(errMsg).Error()
		slog.ErrorContext(ctx, "[service][UpdateModule] err", "%v", errMsg)

		return resp, errors.New(errMsg)
	}

	if req.ModuleName == module.ModuleName {
		errMsg := "module name is same as before"
		resp.Code = http.StatusBadRequest
		resp.Message = errMsg
		resp.Error = errors.New(errMsg).Error()
		slog.ErrorContext(ctx, "[service][UpdateModule] err", "%v", errMsg)

		return resp, errors.New(errMsg)
	}
	model := models.Module{
		ID:         moduleID,
		ModuleName: req.ModuleName,
		UpdatedBy:  req.UpdatedBy,
	}

	if req.Tag != "" {
		model.Tag = sql.NullString{
			String: req.Tag,
			Valid:  true,
		}
	}

	if req.ThumbnailUrl != "" {
		model.ThumbnailUrl = sql.NullString{
			String: req.ThumbnailUrl,
			Valid:  true,
		}
	}

	err = m.ModuleRepo.UpdateModule(ctx, model)
	if err != nil {
		resp.Code = http.StatusBadGateway
		resp.Message = "error while update module, please try again"
		resp.Error = err.Error()
		slog.ErrorContext(ctx, "[service][UpdateModule] while UpdateModule err", "%v", err)
		return
	}

	resp.Data = struct{}{}

	return
}

func (m *ModuleSvcImpl) GetModulesByGroupID(ctx context.Context, groupID int64) (resp models.DefaultResponse, err error) {
	{
		resp.Code = http.StatusOK
		resp.Message = "Success"
		resp.Data = struct{}{}
	}

	modules, err := m.ModuleRepo.GetModulesByGroupID(ctx, groupID)
	if err != nil {
		slog.ErrorContext(ctx, "[service][GetModulesByGroupID] error while GetModulesByGroupID err", "%v", err.Error())
		resp.Code = http.StatusBadRequest
		resp.Message = "bad request"
		resp.Error = errors.New("something went wrong, we will fix it soon").Error()
		return
	}

	var respData []ModuleResponse
	for _, module := range modules {
		data := ModuleResponse{
			ID:         module.ID,
			UUID:       module.UUID,
			ModuleName: module.ModuleName,
			CreatedBy:  module.CreatedBy,
			UpdatedBy:  module.UpdatedBy,
			CreatedAt:  module.CreatedAt,
			UpdatedAt:  module.UpdatedAt,
		}
		if module.GroupID.Valid {
			group, errGroup := m.GroupRepo.GetGroupByID(ctx, module.GroupID.Int64)
			if errGroup != nil {
				slog.ErrorContext(ctx, "[service][GetModulesByGroupID] error while GetGroupByID err", "%v", errGroup.Error())
				resp.Code = http.StatusInternalServerError
				resp.Message = "internal server error, we will fix it soon"
				resp.Error = errors.New("something went wrong, we will fix it soon").Error()
				return
			}
			data.GroupID = module.GroupID.Int64
			data.GroupName = group.GroupName
		}

		if module.ThumbnailUrl.Valid {
			data.ThumbnailUrl = module.ThumbnailUrl.String
		}

		if module.Tag.Valid {
			data.Tag = module.Tag.String
		}

		respData = append(respData, data)
	}

	{
		resp.Data = respData
	}
	return
}

func (m *ModuleSvcImpl) DeleteModule(ctx context.Context, moduleID int64) (resp models.DefaultResponse, err error) {
	{
		resp.Code = http.StatusOK
		resp.Message = "Success"
	}

	adminData, ok := ctx.Value(middleware.UserData).(middleware.UserCtxReq)
	if !ok {
		resp.Code = http.StatusBadRequest
		resp.Message = "bad request"
		resp.Error = errors.New("something went wrong, we will fix it soon").Error()
		return
	}

	module, err := m.ModuleRepo.GetModuleByID(ctx, moduleID)
	if err != nil {
		resp.Code = http.StatusBadRequest
		resp.Message = "bad request"
		resp.Error = err
		slog.ErrorContext(ctx, "[service][DeleteModule] error while GetModuleByID err", "%v", err)
		return
	}

	if module.ID == 0 {
		errMsg := "module not found"
		resp.Code = http.StatusBadRequest
		resp.Message = errMsg
		resp.Error = err
		slog.ErrorContext(ctx, "[service][DeleteModule] err", "%v", errMsg)

		return resp, errors.New(errMsg)
	}

	err = m.ModuleRepo.SoftDeleteModule(ctx, moduleID, adminData.Username)
	if err != nil {
		slog.ErrorContext(ctx, "[service][DeleteModule] error while SoftDeleteModule err", "%v", err)
		resp.Code = http.StatusInternalServerError
		resp.Message = "internal server error, we will fix it soon"
		resp.Error = errors.New("something went wrong, we will fix it soon").Error()
		return
	}

	err = m.SubModuleRepo.RemoveModuleIDFromSubModules(ctx, moduleID, adminData.Username)
	if err != nil {
		slog.ErrorContext(ctx, "[service][DeleteModule] error while RemoveModuleIDFromSubModules err", "%v", err)
		resp.Code = http.StatusInternalServerError
		resp.Message = "internal server error, we will fix it soon"
		resp.Error = errors.New("something went wrong, we will fix it soon").Error()
		return
	}

	return
}
