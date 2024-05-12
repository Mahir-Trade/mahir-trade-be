package controller

import (
	"log/slog"
	"mahir-trade-be/internal/app/models"
	"mahir-trade-be/internal/app/service"
	"mahir-trade-be/internal/app/service/utils"
	"net/http"
	"os"

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
		LoginWithGoogle(ec echo.Context) error
		CallbackGoogle(ec echo.Context) error
		AssignRoleDiscordToUser(ec echo.Context) error
		RemoveRoleDiscordUser(ec echo.Context) error
		InviteDiscordUserToGuild(ec echo.Context) error
		ConnectDiscordAccountAndAssignRole(ec echo.Context) error
		ConnectDiscordAccountAndRemoveRole(ec echo.Context) error
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
	code := ec.QueryParam("code")
	res, err := ox.UserSvc.AssignRoleDiscordToUser(ctx, userUUID, code)
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
	res, err := ox.UserSvc.RemoveRoleDiscordToUser(ctx, userUUID)
	if err != nil {
		slog.Error("RemoveRoleDiscordUser - something went wrong", err)
		return ec.JSON(http.StatusBadRequest, res)
	}

	return ec.JSON(http.StatusOK, res)
}

func (ox *AuthCtrlImpl) InviteDiscordUserToGuild(ec echo.Context) error {
	ctx := ec.Request().Context()

	defer func() {
		if r := recover(); r != nil {
			slog.Error("ConnectDiscordAccount - something went wrong", r)
		}
	}()

	code := ec.QueryParam("code")
	redirectURI := os.Getenv("DISCORD_REDIRECT_URI")
	res, err := ox.UserSvc.InviteDiscordUserToGuild(ctx, code, redirectURI)
	if err != nil {
		slog.Error("ConnectDiscordAccount - something went wrong", err)
		return ec.JSON(http.StatusBadRequest, res)
	}

	return ec.JSON(http.StatusOK, res)
}

func (ox *AuthCtrlImpl) ConnectDiscordAccountAndAssignRole(ec echo.Context) error {
	ctx := ec.Request().Context()

	defer func() {
		if r := recover(); r != nil {
			slog.Error("ConnectDiscordAccountAndAssignRole - something went wrong", r)
		}
	}()

	code := ec.QueryParam("code")
	res, err := ox.UserSvc.ConnectDiscordAccountAndAssignRole(ctx, code)
	if err != nil {
		slog.Error("ConnectDiscordAccountAndAssignRole - something went wrong", err)
		return ec.JSON(http.StatusBadRequest, res)
	}

	return ec.JSON(http.StatusOK, res)
}

func (ox *AuthCtrlImpl) ConnectDiscordAccountAndRemoveRole(ec echo.Context) error {
	ctx := ec.Request().Context()

	defer func() {
		if r := recover(); r != nil {
			slog.Error("ConnectDiscordAccountAndRemoveRole - something went wrong", r)
		}
	}()

	code := ec.QueryParam("code")
	res, err := ox.UserSvc.ConnectDiscordAccountAndRemoveRole(ctx, code)
	if err != nil {
		slog.Error("ConnectDiscordAccountAndRemoveRole - something went wrong", err)
		return ec.JSON(http.StatusBadRequest, res)
	}

	return ec.JSON(http.StatusOK, res)
}

func (ox *AuthCtrlImpl) LoginWithGoogle(ec echo.Context) error {
	ctx := ec.Request().Context()

	defer func() {
		if r := recover(); r != nil {
			slog.Error("LoginWithGoogle - something went wrong", r)
		}
	}()

	url, err := ox.UserSvc.LoginWithGoogle(ctx)
	if err != nil {
		slog.Error("LoginWithGoogle - something went wrong", err)
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "bad request",
			Error:   err.Error(),
		})
	}

	return ec.Redirect(http.StatusTemporaryRedirect, url)
}

func (ox *AuthCtrlImpl) CallbackGoogle(ec echo.Context) error {
	ctx := ec.Request().Context()

	defer func() {
		if r := recover(); r != nil {
			slog.Error("CallbackGoogle - something went wrong", r)
		}
	}()

	// var req service.GoogleLoginReq
	// if err := ec.Bind(&req); err != nil {
	// 	slog.ErrorContext(ctx, "[controller][CallbackGoogle]", err)
	// 	return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
	// 		Code:    http.StatusBadRequest,
	// 		Message: "bad request",
	// 		Data:    struct{}{},
	// 		Error:   err.Error(),
	// 	})
	// }

	state := ec.QueryParam("state")
	code := ec.QueryParam("code")

	if state == "" || code == "" {
		return ec.JSON(http.StatusBadRequest, ErrorMessage{
			Indonesian: "Invalid request",
			English:    "Invalid request",
			Error:      "state or code is empty",
		})
	}

	req := service.GoogleLoginReq{
		State: state,
		Code:  code,
	}

	res, err := ox.UserSvc.CallbackGoogle(ctx, req)
	if err != nil {
		slog.ErrorContext(ctx, "[controller][CallbackGoogle]", err)
		return ec.JSON(http.StatusBadRequest, res)
	}

	return ec.JSON(http.StatusOK, res)
}
