package controller

import (
	"log/slog"
	"mahir-trade-be/internal/app/models"
	"mahir-trade-be/internal/app/repo/discord"
	"mahir-trade-be/internal/app/service"
	"mahir-trade-be/internal/app/service/utils"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"go.uber.org/dig"
)

type (
	ErrorMessage struct {
		Indonesian string `json:"indonesian"`
		English    string `json:"english"`
		Error      any    `json:"error,omitempty"`
	}

	Resp struct {
		Code       int    `json:"code"`
		Message    string `json:"message,omitempty"`
		Data       any    `json:"data,omitempty"`
		ErrMessage string `json:"error_message,omitempty"`
	}

	AuthCtrl interface {
		UserRegistration(ec echo.Context) error
		UserLogin(ec echo.Context) error
		AssignRoleDiscordToUser(ec echo.Context) error
		RemoveRoleDiscordUser(ec echo.Context) error
	}

	AuthCtrlImpl struct {
		dig.In

		UserSvc service.UserSvc
	}
)

func NewAuthCtrl(impl AuthCtrlImpl) AuthCtrl {
	return &impl
}

func (ox *AuthCtrlImpl) UserRegistration(ec echo.Context) error {
	ctx := ec.Request().Context()

	defer func() {
		if r := recover(); r != nil {
			slog.Error("UserRegistration - something went wrong", r)
		}
	}()

	var user models.User

	if err := ec.Bind(&user); err != nil {
		return ec.JSON(http.StatusBadRequest, ErrorMessage{
			Indonesian: "Invalid request body",
			English:    "Invalid request body",
			Error:      err.Error(),
		})
	}

	validate := utils.Validate

	err := validate.Struct(user)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		return ec.JSON(http.StatusBadRequest, ErrorMessage{
			Indonesian: "Invalid request body",
			English:    "invalid request body",
			Error:      errors.Error(),
		})
	}

	err = ox.UserSvc.UserRegistration(ctx, user)
	if err != nil {
		return ec.JSON(http.StatusBadRequest, ErrorMessage{
			Indonesian: "bad request",
			English:    "bad request",
			Error:      err.Error(),
		})
	}

	return ec.JSON(http.StatusOK, Resp{Code: http.StatusCreated, Message: "successfully registered"})
}

func (ox *AuthCtrlImpl) UserLogin(ec echo.Context) error {

	ctx := ec.Request().Context()

	defer func() {
		if r := recover(); r != nil {
			slog.Error("UserLogin - something went wrong", r)
		}
	}()

	var user service.LoginReq

	if err := ec.Bind(&user); err != nil {
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	validate := utils.Validate

	err := validate.Struct(user)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   errors.Error(),
		})
	}

	res, err := ox.UserSvc.UserLogin(ctx, user)
	if err != nil {
		slog.Error("UserLogin - something went wrong", err)
		return ec.JSON(http.StatusBadRequest, res)
	}

	return ec.JSON(http.StatusOK, res)

}

func (ox *AuthCtrlImpl) AssignRoleDiscordToUser(ec echo.Context) error {
	ctx := ec.Request().Context()

	defer func() {
		if r := recover(); r != nil {
			slog.Error("AssignRoleDiscordToUser - something went wrong", r)
		}
	}()

	userUUID := ec.Get("user_id").(string)

	var discordReq discord.GetDiscordUserRequest
	if err := ec.Bind(&discordReq); err != nil {
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	res, err := ox.UserSvc.AssignRoleDiscordToUser(ctx, userUUID, discordReq.Username)
	if err != nil {
		slog.Error("AssignRoleDiscordToUser - something went wrong", err)
		return ec.JSON(http.StatusBadRequest, res)
	}

	return ec.JSON(http.StatusOK, res)
}

func (ox *AuthCtrlImpl) RemoveRoleDiscordUser(ec echo.Context) error {
	ctx := ec.Request().Context()

	defer func() {
		if r := recover(); r != nil {
			slog.Error("RemoveRoleDiscordUser - something went wrong", r)
		}
	}()

	userUUID := ec.Get("user_id").(string)

	var discordReq discord.GetDiscordUserRequest
	if err := ec.Bind(&discordReq); err != nil {
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	res, err := ox.UserSvc.RemoveRoleDiscordToUser(ctx, userUUID, discordReq.Username)
	if err != nil {
		slog.Error("RemoveRoleDiscordUser - something went wrong", err)
		return ec.JSON(http.StatusBadRequest, res)
	}

	return ec.JSON(http.StatusOK, res)
}
