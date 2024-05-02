package service

import (
	"context"
	"fmt"
	"log/slog"
	"mahir-trade-be/internal/app/models"
	"mahir-trade-be/internal/app/repo/postgres"
	"mahir-trade-be/internal/app/service/utils"
	"net/http"
	"strings"

	"go.uber.org/dig"
)

type (
	LoginReq struct {
		Identity string `json:"identity" validate:"required"`
		Password string `json:"password" validate:"required"`
	}
	JWTData struct {
		Email   string `json:"email"`
		UserID  string `json:"user_id"`
		Usename string `json:"username"`
	}

	UserSvc interface {
		UserRegistration(ctx context.Context, req models.User) (err error)
		UserLogin(ctx context.Context, req LoginReq) (resp models.DefaultResponse, err error)
	}

	UserSvcImpl struct {
		dig.In

		UserRepo postgres.UserRepo
	}
)

func NewUserSvc(impl UserSvcImpl) UserSvc {
	return &impl
}

func (u *UserSvcImpl) UserRegistration(ctx context.Context, req models.User) (err error) {

	{
		req.Email = strings.ToLower(strings.TrimSpace(req.Email))
		req.PhoneNumber = strings.TrimSpace(req.PhoneNumber)
		req.Username = strings.ToLower(strings.TrimSpace(req.Username))
		req.Password, err = utils.HashPassword(req.Password)
		if err != nil {
			slog.ErrorContext(ctx, fmt.Sprintf("[service][HashPassword] err : %v", err))
			return err
		}
	}

	_, err = u.UserRepo.CreateUser(ctx, req)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("[service][UserRegistration] err : %v", err.Error()))
		return err
	}

	return nil
}

func (u *UserSvcImpl) UserLogin(ctx context.Context, req LoginReq) (resp models.DefaultResponse, err error) {

	{
		resp.Code = http.StatusOK
		resp.Message = "success"
		req.Identity = strings.ToLower(strings.TrimSpace(req.Identity))
	}

	user, err := u.UserRepo.FindUserByEmailOrUsername(ctx, req.Identity)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("[service][UserLogin][FindUserByEmail] err : %v", err))
		err = fmt.Errorf("invalid email or password")
		resp.Code = http.StatusUnauthorized
		resp.Message = "invalid email or password"
		resp.Error = err
		return
	}

	if err = utils.VerifyPassword(req.Password, user.Password); err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("[service][UserLogin][VerifyPassword] err : %v", err))
		err = fmt.Errorf("invalid email or password")
		resp.Code = http.StatusUnauthorized
		resp.Message = "invalid email or password"
		resp.Error = err
		return
	}

	token, exp, err := utils.Sign(JWTData{
		Email:   user.Email,
		UserID:  user.UUID,
		Usename: user.Username,
	})

	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("[service][UserLogin][Sign] err : %v", err))
		err = fmt.Errorf("internal server error, we will fix it soon")
		resp.Code = http.StatusInternalServerError
		resp.Message = "internal server error, we will fix it soon"
		return
	}

	resp.Data = struct {
		Token  string `json:"token"`
		Expire int64  `json:"expire"`
	}{
		Token:  token,
		Expire: exp,
	}

	return
}
