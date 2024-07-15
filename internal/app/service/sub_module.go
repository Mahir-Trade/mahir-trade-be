package service

import (
	"context"
	"database/sql"
	"io"
	"log/slog"
	"mahir-trade-be/internal/app/infra"
	"mahir-trade-be/internal/app/models"
	"mahir-trade-be/internal/app/repo/google"
	"mahir-trade-be/internal/app/repo/postgres"
	"mahir-trade-be/pkg/middleware"
	"math"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"go.uber.org/dig"
)

type (
	SubModuleRequest struct {
		ModuleID      int64  `json:"module_id"`
		SubModuleName string `json:"sub_module_name" validate:"required"`
		Title         string `json:"title" validate:"required"`
		VideoURL      string `json:"video_url" validate:"required"`
	}

	SubModuleResponse struct {
		ID            int64  `json:"id"`
		UUID          string `json:"uuid"`
		ModuleID      int64  `json:"module_id"`
		ModuleName    string `json:"module_name"`
		SubModuleName string `json:"sub_module_name"`
		Title         string `json:"title"`
		VideoURL      string `json:"video_url"`
		CreatedBy     string `json:"created_by"`
		UpdatedBy     string `json:"updated_by"`
		CreatedAt     string `json:"created_at,omitempty"`
		UpdatedAt     string `json:"updated_at,omitempty"`
	}

	MarkSubModuleAsWatchedRequest struct {
		SubModuleID int64 `json:"sub_module_id"`
	}

	SubModuleSvc interface {
		CreateSubModule(ctx context.Context, req SubModuleRequest) (resp models.DefaultResponse, err error)
		GetSubModuleByID(ctx context.Context, id int64) (resp models.DefaultResponse, err error)
		GetSubModules(ctx context.Context, req models.PaginationRequest) (resp models.DefaultPaginationResponseData, err error)
		GetSubModulesByModuleID(ctx context.Context, moduleID int64) (resp models.DefaultResponse, err error)
		UpdateSubModule(ctx context.Context, id int64, req SubModuleRequest) (resp models.DefaultResponse, err error)

		SoftDeleteSubModule(ctx context.Context, subModuleId int64) (resp models.DefaultResponse, err error)
		MarkSubModuleAsWatched(ctx context.Context, req MarkSubModuleAsWatchedRequest) (resp models.DefaultResponse, err error)
		UploadFile(ctx context.Context, req google.FileUpload, src multipart.File) (resp models.DefaultResponse, err error)
	}

	SubModuleSvcImpl struct {
		dig.In

		ModuleRepo        postgres.ModuleRepo
		SubModuleRepo     postgres.SubModuleRepo
		BucketRepo        google.BucketRepo
		UserSubModuleRepo postgres.UserSubModuleRepo

		GoogleCfg *infra.GoogleCfg
	}
)

func NewSubModuleSvc(impl SubModuleSvcImpl) SubModuleSvc {
	return &impl
}

func (s *SubModuleSvcImpl) CreateSubModule(ctx context.Context, req SubModuleRequest) (resp models.DefaultResponse, err error) {
	{
		resp.Code = http.StatusOK
		resp.Message = "Success"
		resp.Data = struct{}{}
	}
	adminData, ok := ctx.Value(middleware.UserData).(middleware.UserCtxReq)
	if !ok {
		slog.ErrorContext(ctx, "[service][CreateSubModule] error while get user data from context")
		resp.Code = http.StatusUnauthorized
		resp.Message = "Unauthorized"
		resp.Error = "forbidden access, role not allowed"
		return
	}

	subModuleEntry := models.SubModule{
		SubModuleName: req.SubModuleName,
		Title:         req.Title,
		VideoURL:      req.VideoURL,
		CreatedBy:     adminData.Username,
	}

	if req.ModuleID != 0 {
		_, err = s.ModuleRepo.GetModuleByID(ctx, req.ModuleID)
		if err != nil {
			slog.ErrorContext(ctx, "[service][CreateSubModule] error while get module by id err: %v", err)
			resp.Code = http.StatusNotFound
			resp.Message = "module not found"
			resp.Error = err.Error()
			return
		}

		subModuleEntry.ModuleID = sql.NullInt64{
			Int64: req.ModuleID,
			Valid: true,
		}
	}

	id, err := s.SubModuleRepo.CreateSubModule(ctx, subModuleEntry)
	if err != nil {
		slog.ErrorContext(ctx, "[service][CreateSubModule] error while create sub module err: %v", err)
		resp.Code = http.StatusInternalServerError
		resp.Message = "Internal Server Error, Please try again later"
		resp.Error = "Internal Server Error"
		return
	}
	resp.Data = struct {
		ID int `json:"id"`
	}{
		ID: id,
	}
	return
}

func (s *SubModuleSvcImpl) GetSubModuleByID(ctx context.Context, id int64) (resp models.DefaultResponse, err error) {
	{
		resp.Code = http.StatusOK
		resp.Message = "Success"
		resp.Data = struct{}{}
	}
	subModule, err := s.SubModuleRepo.GetSubModuleByID(ctx, id)
	if err != nil {
		slog.ErrorContext(ctx, "[service][GetSubModuleByID] error while get sub module by id err: %v", err)
		resp.Code = http.StatusInternalServerError
		resp.Message = "Internal Server Error, Please try again later"
		resp.Error = "Internal Server Error"
		return
	}

	url, err := s.BucketRepo.PresignedURL(ctx, s.GoogleCfg.VideoBucketName, subModule.VideoURL)
	if err != nil {
		slog.ErrorContext(ctx, "[service][GetSubModuleByID] error while get presigned url err: %v", err)
		resp.Code = http.StatusInternalServerError
		resp.Message = "Internal Server Error, Please try again later"
		resp.Error = "Internal Server Error"
		return
	}

	subModule.VideoURL = url

	resp.Data = subModule
	return
}

func (s *SubModuleSvcImpl) GetSubModules(ctx context.Context, req models.PaginationRequest) (resp models.DefaultPaginationResponseData, err error) {
	var (
		dataResp       models.DefaultResponse
		submodulesData []SubModuleResponse
	)
	{
		dataResp.Code = http.StatusOK
		dataResp.Message = "Success"
		dataResp.Data = struct{}{}

		req.Page = req.Page - 1
	}

	subModules, totalData, err := s.SubModuleRepo.GetSubModules(ctx, req)
	if err != nil {
		slog.ErrorContext(ctx, "[service][GetSubModules] error while get sub modules err: %v", err)
		dataResp.Code = http.StatusNotFound
		dataResp.Message = "Data Not Found"
		dataResp.Error = "Data Not Found"
		return
	}

	for _, subModule := range subModules {
		submoduleData := SubModuleResponse{
			ID:            subModule.ID,
			UUID:          subModule.UUID,
			SubModuleName: subModule.SubModuleName,
			Title:         subModule.Title,
			VideoURL:      subModule.VideoURL,
			CreatedBy:     subModule.CreatedBy,
			UpdatedBy:     subModule.UpdatedBy,
			CreatedAt:     subModule.CreatedAt,
			UpdatedAt:     subModule.UpdatedAt,
		}

		if subModule.ModuleID.Valid {
			module, errModule := s.ModuleRepo.GetModuleByID(ctx, subModule.ModuleID.Int64)
			if errModule != nil {
				slog.ErrorContext(ctx, "[service][GetSubModules] error while get module by id err: %v", err)
			}
			submoduleData.ModuleID = module.ID
			submoduleData.ModuleName = module.ModuleName
		}
		submodulesData = append(submodulesData, submoduleData)
	}

	{
		dataResp.Data = submodulesData
		resp.Page = uint(req.Page) + 1
		resp.Limit = uint(req.Limit)

		totalPage := math.Ceil(float64(totalData) / float64(req.Limit))
		resp.TotalPages = uint(totalPage)

		resp.TotalItems = uint(len(submodulesData))
		resp.HasNext = resp.Page < resp.TotalPages
		resp.HasPrevious = resp.Page > 1
		resp.Results = dataResp
	}

	return

}

func (s *SubModuleSvcImpl) UpdateSubModule(ctx context.Context, id int64, req SubModuleRequest) (resp models.DefaultResponse, err error) {
	{
		resp.Code = http.StatusOK
		resp.Message = "Success"
		resp.Data = struct{}{}
	}

	adminData, ok := ctx.Value(middleware.UserData).(middleware.UserCtxReq)
	if !ok {
		slog.ErrorContext(ctx, "[service][UpdateSubModule] error while get user data from context")
		resp.Code = http.StatusUnauthorized
		resp.Message = "Unauthorized"
		resp.Error = "forbidden access, role not allowed"
		return
	}

	entryUpdateSubModule := models.SubModule{
		ID:            id,
		SubModuleName: req.SubModuleName,
		Title:         req.Title,
		VideoURL:      req.VideoURL,
		UpdatedBy:     adminData.Username,
	}

	if req.ModuleID == 0 {
		slog.ErrorContext(ctx, "[service][UpdateSubModule] error while update sub module, module id is required")
		resp.Code = http.StatusBadRequest
		resp.Message = "Bad Request"
		resp.Error = "module ID is required"
		return
	}
	entryUpdateSubModule.ModuleID = sql.NullInt64{
		Int64: req.ModuleID,
		Valid: true,
	}
	err = s.SubModuleRepo.UpdateSubModule(ctx, entryUpdateSubModule)
	if err != nil {
		slog.ErrorContext(ctx, "[service][UpdateSubModule] error while update sub module err: %v", err)
		resp.Code = http.StatusInternalServerError
		resp.Message = "error while updating sub Module"
		resp.Error = "error while Updating sub module"
		return
	}

	return
}

func (s *SubModuleSvcImpl) SoftDeleteSubModule(ctx context.Context, subModuleId int64) (resp models.DefaultResponse, err error) {
	{
		resp.Code = http.StatusOK
		resp.Message = "Success"
		resp.Data = struct{}{}
	}

	adminData, ok := ctx.Value(middleware.UserData).(middleware.UserCtxReq)
	if !ok {
		slog.ErrorContext(ctx, "[service][SoftDeleteSubModule] error while get user data from context")
		resp.Code = http.StatusUnauthorized
		resp.Message = "Unauthorized"
		resp.Error = "forbidden access, role not allowed"
		return
	}

	err = s.SubModuleRepo.SoftDeleteSubModule(ctx, subModuleId, adminData.Username)
	if err != nil {
		slog.ErrorContext(ctx, "[service][SoftDeleteSubModule] error while soft delete sub module err: %v", err)
		resp.Code = http.StatusInternalServerError
		resp.Message = "Internal Server Error, Please try again later"
		resp.Error = "Internal Server Error"
		return
	}

	return
}

func (s *SubModuleSvcImpl) GetSubModulesByModuleID(ctx context.Context, moduleID int64) (resp models.DefaultResponse, err error) {
	var (
		submodulesData []SubModuleResponse
	)
	{
		resp.Code = http.StatusOK
		resp.Message = "Success"
		resp.Data = struct{}{}
	}

	subModules, err := s.SubModuleRepo.GetSubModulesByModuleID(ctx, moduleID)
	if err != nil {
		slog.ErrorContext(ctx, "[service][GetSubModulesByModuleID] error while get sub modules by module id err: %v", err)
		resp.Code = http.StatusNotFound
		resp.Message = "Data Not Found"
		resp.Error = "Data Not Found"
		return
	}

	for _, subModule := range subModules {
		submoduleData := SubModuleResponse{
			ID:            subModule.ID,
			UUID:          subModule.UUID,
			SubModuleName: subModule.SubModuleName,
			Title:         subModule.Title,
			VideoURL:      subModule.VideoURL,
			CreatedBy:     subModule.CreatedBy,
			UpdatedBy:     subModule.UpdatedBy,
			CreatedAt:     subModule.CreatedAt,
			UpdatedAt:     subModule.UpdatedAt,
		}

		if subModule.ModuleID.Valid {
			module, errModule := s.ModuleRepo.GetModuleByID(ctx, subModule.ModuleID.Int64)
			if errModule != nil {
				slog.ErrorContext(ctx, "[service][GetSubModulesByModuleID] error while get module by id err: %v", err)
			}
			submoduleData.ModuleID = module.ID
			submoduleData.ModuleName = module.ModuleName
			submodulesData = append(submodulesData, submoduleData)
		}
	}

	{
		resp.Data = submodulesData
	}

	return
}

func (s *SubModuleSvcImpl) MarkSubModuleAsWatched(ctx context.Context, req MarkSubModuleAsWatchedRequest) (resp models.DefaultResponse, err error) {
	{
		resp.Code = http.StatusOK
		resp.Message = "Success"
	}

	userData, ok := ctx.Value(middleware.UserData).(middleware.UserCtxReq)
	if !ok {
		slog.ErrorContext(ctx, "[service][MarkSubModuleAsWatched] error while get user data from context")
		resp.Code = http.StatusUnauthorized
		resp.Message = "Unauthorized"
		resp.Error = "forbidden access, role not allowed"
		return
	}

	_, err = s.SubModuleRepo.GetSubModuleByID(ctx, req.SubModuleID)
	if err != nil {
		slog.ErrorContext(ctx, "[service][MarkSubModuleAsWatched] error while get sub module by id err: %v", err)
		resp.Code = http.StatusNotFound
		resp.Message = "Sub Module Not Found"
		resp.Error = "Sub Module Not Found"
		return
	}

	userSubModule, err := s.UserSubModuleRepo.GetUserSubModuleBySubModuleIDAndUserID(ctx, userData.UserID, req.SubModuleID)
	if err != nil {
		slog.ErrorContext(ctx, "[service][MarkSubModuleAsWatched] error while get sub module by id err: %v", err)
		resp.Code = http.StatusBadRequest
		resp.Message = "bad request"
		resp.Error = err.Error()
		return
	}

	if userSubModule.ID != 0 {
		return
	}

	userSubModuleReq := models.UserSubModule{
		UserID:      userData.UserID,
		SubModuleID: req.SubModuleID,
		CreatedBy:   userData.Email,
	}

	_, err = s.UserSubModuleRepo.CreateUserSubModule(ctx, userSubModuleReq)
	if err != nil {
		slog.ErrorContext(ctx, "[service][MarkSubModuleAsWatched] error while create user sub module err: %v", err)
		resp.Code = http.StatusInternalServerError
		resp.Message = "Internal Server Error, Please try again later"
		resp.Error = "Internal Server Error"
		return
	}

	return
}

func (s *SubModuleSvcImpl) UploadFile(ctx context.Context, req google.FileUpload, src multipart.File) (resp models.DefaultResponse, err error) {
	{
		resp.Code = http.StatusOK
		resp.Message = "Success"
	}

	tempDir := "./temp"
	if _, err = os.Stat(tempDir); os.IsNotExist(err) {
		if err = os.MkdirAll(tempDir, 0755); err != nil {
			slog.ErrorContext(ctx, "[service][UploadFile] error while MkdirAll err: %v", err)
			resp.Code = http.StatusInternalServerError
			resp.Message = "Unable to create temporary file"
			resp.Error = err.Error()

			return
		}
	}

	tempFile, err := os.CreateTemp("./temp", "upload-*.tmp")
	if err != nil {
		slog.ErrorContext(ctx, "[service][UploadFile] error while CreateTemp err: %v", err)
		resp.Code = http.StatusInternalServerError
		resp.Message = "Unable to create temporary file"
		resp.Error = err.Error()

		return
	}
	defer func() {
		tempFile.Close()
		if err := os.Remove(tempFile.Name()); err != nil {
			slog.ErrorContext(ctx, "[service][UploadFile] Failed to remove temporary file %s: %v", tempFile.Name(), err)
		}
	}()

	if _, err = io.Copy(tempFile, src); err != nil {
		slog.ErrorContext(ctx, "[service][UploadFile] error while Copy err: %v", err)
		resp.Code = http.StatusInternalServerError
		resp.Message = "Unable to save temporary file"
		resp.Error = err.Error()

		return
	}

	ext := filepath.Ext(req.Filename)
	contentType := mime.TypeByExtension(ext)
	req.FileContentType = contentType
	req.LocalFilePath = tempFile.Name()

	if strings.Contains(contentType, "video") {
		req.BucketName = s.GoogleCfg.VideoBucketName
	} else {
		req.BucketName = s.GoogleCfg.ImageBucketName
	}

	url, err := s.BucketRepo.UploadFile(ctx, req, src)
	if err != nil {
		slog.ErrorContext(ctx, "[service][UploadFile] error while upload file err: %v", err)
		resp.Code = http.StatusInternalServerError
		resp.Message = "Internal Server Error, Please try again later"
		resp.Error = "Internal Server Error"

		return
	}

	resp.Data = struct {
		URL string `json:"url"`
	}{
		URL: url,
	}

	return
}
