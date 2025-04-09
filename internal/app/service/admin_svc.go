package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"mahir-trade-be/internal/app/models"
	"mahir-trade-be/internal/app/repo/postgres"
	"mahir-trade-be/internal/app/service/utils"
	"mahir-trade-be/pkg/middleware"
	"math"
	"net/http"
	"strings"
	"time"

	"go.uber.org/dig"
)

type (
	AdminLoginRequest struct {
		Identity string `json:"identity" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	UpdateTypeUserRequest struct {
		IsActive bool `json:"is_active" validate:"required"`
	}

	ToggleUserMembershipRequest struct {
		UserIds []int64 `json:"user_ids" validate:"required,dive"`
	}

	AdminSvc interface {
		AdminLogin(ctx context.Context, req AdminLoginRequest) (resp models.DefaultResponse, err error)
		AdminRegistration(ctx context.Context, req models.Admin) (resp models.DefaultResponse, err error)
		UpdateTypeUser(ctx context.Context, isActive bool, id int64) (resp models.DefaultResponse, err error)
		GetDetailAdminInfo(ctx context.Context, username string) (resp models.DefaultResponse, err error)
		GetAllUsers(ctx context.Context, req models.PaginationRequest) (resp models.DefaultPaginationResponseData, err error)
		ToggleInactiveUserMembership(ctx context.Context, req ToggleUserMembershipRequest) (resp models.DefaultResponse, err error)
	}

	AdminSvcImpl struct {
		dig.In
		AdminRepo          postgres.AdminRepo
		UserRepo           postgres.UserRepo
		UserMembershipRepo postgres.UserMembershipRepo
	}
)

func NewAdminSvc(impl AdminSvcImpl) AdminSvc {
	return &impl
}

func (a *AdminSvcImpl) AdminLogin(ctx context.Context, req AdminLoginRequest) (resp models.DefaultResponse, err error) {
	{
		resp = models.DefaultResponse{
			Code:    http.StatusOK,
			Message: "Login Success",
			Data:    struct{}{},
		}
	}

	admin, err := a.AdminRepo.FindByUsername(ctx, req.Identity)
	if err != nil {
		slog.ErrorContext(ctx, "[service][AdminLogin] error while FindByUsername err: %v", err)
		resp.Code = http.StatusUnauthorized
		resp.Message = "Invalid Username or Password"
		return
	}

	if err = utils.VerifyPassword(req.Password, admin.Password); err != nil {
		slog.ErrorContext(ctx, "[service][AdminLogin] error while VerifyPassword err", "%v", err)
		resp.Code = http.StatusUnauthorized
		resp.Message = "Invalid Username or Password"
		return
	}

	dataJwt := middleware.UserCtxReq{
		UserID:   admin.AdminID,
		Email:    admin.Email,
		Username: admin.Username,
	}

	token, exp, err := utils.Sign(dataJwt)
	if err != nil {
		slog.ErrorContext(ctx, "[service][AdminLogin] error while Sign err: %v", err)
		resp.Code = http.StatusInternalServerError
		resp.Message = "Internal Server Error, we will fix it soon"
		resp.Error = "Internal Server Error, we will fix it soon"
		return
	}

	resp.Data = struct {
		Token string    `json:"token"`
		Exp   time.Time `json:"exp"`
	}{
		Token: token,
		Exp:   exp,
	}

	return

}

func (a *AdminSvcImpl) AdminRegistration(ctx context.Context, req models.Admin) (resp models.DefaultResponse, err error) {
	{
		resp = models.DefaultResponse{
			Code:    http.StatusCreated,
			Message: "Registration Success",
		}
		req.Email = strings.ToLower(strings.TrimSpace(req.Email))
		req.Username = strings.ToLower(strings.TrimSpace(req.Username))
		req.Password, err = utils.HashPassword(req.Password)
		if err != nil {
			slog.ErrorContext(ctx, "[service][AdminRegistration] error while HashPassword err: %v", err)
			resp.Code = http.StatusBadRequest
			resp.Message = "Bad Request"
			resp.Error = err.Error()
			return
		}
	}

	adminId, err := a.AdminRepo.CreateAdmin(ctx, req)
	if err != nil {
		slog.ErrorContext(ctx, "[service][AdminRegistration] error while CreateAdmin err: %v", err)
		resp.Code = http.StatusInternalServerError
		resp.Message = "Bad Request"
		resp.Error = err.Error()
		return
	}

	dataJwt := middleware.UserCtxReq{
		UserID:   int64(adminId),
		Email:    req.Email,
		Username: req.Username,
	}

	token, exp, err := utils.Sign(dataJwt)
	if err != nil {
		slog.ErrorContext(ctx, "[service][AdminRegistration] error while Sign err: %v", err)
		resp.Code = http.StatusInternalServerError
		resp.Message = utils.ErrorInternalServer
		resp.Error = err.Error()
		return
	}

	resp.Data = struct {
		Token string    `json:"token"`
		Exp   time.Time `json:"exp"`
	}{
		Token: token,
		Exp:   exp,
	}

	return
}

func (a *AdminSvcImpl) UpdateTypeUser(ctx context.Context, isActive bool, id int64) (resp models.DefaultResponse, err error) {
	{
		resp = models.DefaultResponse{
			Code:    http.StatusOK,
			Message: "Success",
			Data:    struct{}{},
		}
	}

	adminData, ok := ctx.Value(middleware.UserData).(middleware.UserCtxReq)
	if !ok {
		err = errors.New("failed to get admin data from context")
		slog.ErrorContext(ctx, "[service][UpdateTypeUser] error while getting admin data from context: %v", err)
		resp.Code = http.StatusInternalServerError
		resp.Message = "Internal Server Error"
		resp.Error = err.Error()
		return
	}

	user, err := a.UserRepo.GetUserByID(ctx, id)
	if err != nil {
		slog.ErrorContext(ctx, "[service][UpdateTypeUser] error while GetUserByID err: %v", err)
		resp.Code = http.StatusBadRequest
		resp.Message = "Bad Request"
		resp.Error = err.Error()
	}

	if user.UserID == 0 {
		resp.Code = http.StatusNotFound
		resp.Message = "User Not Found"
		return
	}

	if user.IsActive == isActive {
		resp.Code = http.StatusBadRequest
		resp.Message = "User already in this state"
		return
	}

	state, err := a.UserRepo.UpdateTypeUser(ctx, isActive, adminData.Username, int64(user.UserID))
	if err != nil {
		slog.ErrorContext(ctx, "[service][UpdateTypeUser] error while UpdateTypeUser err: %v", err)
		resp.Code = http.StatusInternalServerError
		resp.Message = "Internal Server Error"
		resp.Error = err.Error()
		return
	}

	resp.Data = state

	return
}

func (a *AdminSvcImpl) GetDetailAdminInfo(ctx context.Context, username string) (resp models.DefaultResponse, err error) {
	{
		resp = models.DefaultResponse{
			Code:    http.StatusOK,
			Message: "Success",
			Data:    struct{}{},
		}
	}

	admin, err := a.AdminRepo.FindByUsername(ctx, username)
	if err != nil {
		slog.ErrorContext(ctx, "[service][GetDetailAdminInfo] error while GetAdminByID err: %v", err)
		resp.Code = http.StatusInternalServerError
		resp.Message = "Internal Server Error"
		resp.Error = err.Error()
		return
	}

	admin.Password = ""
	resp.Data = admin

	return
}

func (a *AdminSvcImpl) GetAllUsers(ctx context.Context, req models.PaginationRequest) (resp models.DefaultPaginationResponseData, err error) {
	var dataResp models.DefaultResponse
	{
		dataResp.Code = http.StatusOK
		dataResp.Message = "success"
	}

	users, totalCount, err := a.UserRepo.GetAllUser(ctx, req)
	if err != nil {
		dataResp.Code = http.StatusBadRequest
		dataResp.Message = utils.ErrorBadRequest
		dataResp.Error = err.Error()
		slog.ErrorContext(ctx, fmt.Sprintf("[service][GetAllUsers] while GetAllUser err: %v", err.Error()))

		return resp, err
	}

	{
		dataResp.Data = users
		resp.Page = uint(req.Page)
		resp.Limit = uint(req.Limit)

		resp.TotalPages = uint(math.Ceil(float64(totalCount) / float64(req.Limit)))
		resp.TotalItems = uint(totalCount)
		resp.HasNext = req.Page < int64(resp.TotalPages)
		resp.HasPrevious = req.Page > 1
		resp.Results = dataResp
	}

	return
}

func (a *AdminSvcImpl) ToggleInactiveUserMembership(ctx context.Context, req ToggleUserMembershipRequest) (resp models.DefaultResponse, err error) {
	{
		resp = models.DefaultResponse{
			Code:    http.StatusOK,
			Message: "Success",
			Data:    struct{}{},
		}
	}

	if len(req.UserIds) == 0 {
		resp.Code = http.StatusBadRequest
		resp.Message = "User IDs are required"
		return resp, errors.New("user IDs are required")
	}

	adminData, ok := ctx.Value(middleware.UserData).(middleware.UserCtxReq)
	if !ok {
		err = errors.New("failed to get admin data from context")
		slog.ErrorContext(ctx, "[service][ToggleInactiveUserMembership] error while getting admin data from context: %v", err)
		resp.Code = http.StatusForbidden
		resp.Message = "Internal Server Error"
		resp.Error = err.Error()
		return resp, err
	}

	err = a.UserMembershipRepo.UpdateUserMembershipExpired(ctx, req.UserIds, adminData.Username)
	if err != nil {
		resp.Code = http.StatusBadRequest
		resp.Message = "Failed to toggle user membership"
	}

	return resp, nil
}
