package service

import (
	"context"
	"errors"
	"log/slog"
	"mahir-trade-be/internal/app/models"
	"mahir-trade-be/internal/app/repo/postgres"
	"mahir-trade-be/internal/app/service/utils"
	"mahir-trade-be/pkg/middleware"
	"net/http"
	"strings"

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

	AdminSvc interface {
		AdminLogin(ctx context.Context, req AdminLoginRequest) (resp models.DefaultResponse, err error)
		AdminRegistration(ctx context.Context, req models.Admin) (resp models.DefaultResponse, err error)
		UpdateTypeUser(ctx context.Context, isActive bool, id int64) (resp models.DefaultResponse, err error)
	}

	AdminSvcImpl struct {
		dig.In

		AdminRepo postgres.AdminRepo
		UserRepo  postgres.UserRepo
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
		Token string `json:"token"`
		Exp   int64  `json:"exp"`
	}{
		Token: token,
		Exp:   exp,
	}

	return

}

func (a *AdminSvcImpl) AdminRegistration(ctx context.Context, req models.Admin) (resp models.DefaultResponse, err error) {
	{
		resp = models.DefaultResponse{
			Code:    http.StatusOK,
			Message: "Registration Success",
			Data:    struct{}{},
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

	_, err = a.AdminRepo.CreateAdmin(ctx, req)
	if err != nil {
		slog.ErrorContext(ctx, "[service][AdminRegistration] error while CreateAdmin err: %v", err)
		resp.Code = http.StatusInternalServerError
		resp.Message = "Bad Request"
		resp.Error = err.Error()
		return
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

	if user.IsActivate == isActive {
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
