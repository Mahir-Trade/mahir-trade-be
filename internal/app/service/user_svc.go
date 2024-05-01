package service

import (
	"context"
	"fmt"
	"log/slog"
	"mahir-trade-be/internal/app/models"
	"mahir-trade-be/internal/app/repo/postgres"
	"mahir-trade-be/internal/app/service/utils"
	"strings"

	"go.uber.org/dig"
)

type (
	UserSvc interface {
		UserRegistration(ctx context.Context, req models.User) (err error)
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
