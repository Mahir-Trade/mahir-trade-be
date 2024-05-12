package middleware

import (
	"context"
	"encoding/json"
	"mahir-trade-be/internal/app/models"
	"mahir-trade-be/internal/app/repo/postgres"
	"mahir-trade-be/internal/app/service/utils"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"go.uber.org/dig"
)

const (
	UserData userDataKey = "user_data"
)

type (
	UserCtxReq struct {
		UserID   int64  `json:"user_id"`
		Email    string `json:"email"`
		Username string `json:"username"`
	}

	MiddleWareImpl struct {
		dig.In

		UserRepo  postgres.UserRepo
		AdminRepo postgres.AdminRepo
	}

	MiddleWare interface {
		AuthAdminOrUser(role ...string) func(next echo.HandlerFunc) echo.HandlerFunc
	}

	userDataKey string
)

func NewMiddleWare(impl MiddleWareImpl) MiddleWare {
	return &impl
}

func (m *MiddleWareImpl) AuthAdminOrUser(role ...string) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			errResponse := models.DefaultResponse{Code: http.StatusUnauthorized, Data: struct{}{}}

			ctx := c.Request().Context()
			if len(role) == 0 {
				errResponse.Message = "Unauthorized"
				return c.JSON(http.StatusUnauthorized, errResponse)
			}

			mRole := map[string]string{
				"admin": "admin",
				"user":  "user",
			}

			for _, r := range role {
				if _, ok := mRole[r]; !ok {
					errResponse.Message = "Unauthorized"
					return c.JSON(http.StatusUnauthorized, errResponse)
				}
			}
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				errResponse.Message = "Missing Authorization header"
				return c.JSON(http.StatusUnauthorized, errResponse)
			}

			splitHeader := strings.Split(authHeader, " ")
			if len(splitHeader) != 2 {
				errResponse.Message = "Invalid Authorization header format"
				return c.JSON(http.StatusUnauthorized, errResponse)
			}

			if splitHeader[0] != "Bearer" {
				errResponse.Message = "Invalid Authorization scheme"
				return c.JSON(http.StatusUnauthorized, errResponse)
			}

			token, err := utils.Verify(splitHeader[1])
			if err != nil {
				errResponse.Message = "Unauthorized"
				errResponse.Error = err.Error()
				return c.JSON(http.StatusUnauthorized, errResponse)
			}

			err = token.Valid()
			if err != nil {
				errResponse.Message = "Unauthorized"
				errResponse.Error = err.Error()
				return c.JSON(http.StatusUnauthorized, errResponse)
			}

			var userCtx UserCtxReq
			bt, err := json.Marshal(token.Data)
			if err != nil {
				errResponse.Message = "Unauthorized"
				errResponse.Error = err.Error()
				return c.JSON(http.StatusUnauthorized, errResponse)
			}

			err = json.Unmarshal(bt, &userCtx)
			if err != nil {
				errResponse.Message = "Unauthorized"
				errResponse.Error = err.Error()
				return c.JSON(http.StatusUnauthorized, errResponse)
			}

			if userCtx.UserID == 0 {
				errResponse.Message = "Unauthorized"
				errResponse.Error = "Invalid user id"
				return c.JSON(http.StatusUnauthorized, errResponse)
			}

			if len(role) == 1 && role[0] == "admin" {
				adminData, errGetAdmin := m.AdminRepo.FindByUsername(ctx, userCtx.Username)
				if errGetAdmin != nil {
					errResponse.Message = "Unauthorized"
					errResponse.Error = "Invalid role"
					return c.JSON(http.StatusUnauthorized, errResponse)
				}
				{
					userCtx.UserID = adminData.AdminID
					userCtx.Email = adminData.Email
					userCtx.Username = adminData.Username
				}
			} else if len(role) == 1 && role[0] == "user" {
				userData, errGetUser := m.UserRepo.FindUserByEmailOrUsername(ctx, userCtx.Username)
				if errGetUser != nil {
					errResponse.Message = "Unauthorized"
					errResponse.Error = "Invalid role"
					return c.JSON(http.StatusUnauthorized, errResponse)
				}
				{
					userCtx.UserID = int64(userData.UserID)
					userCtx.Email = userData.Email
					userCtx.Username = userData.Username
				}
			}

			ctx = context.WithValue(ctx, UserData, userCtx)

			c.SetRequest(c.Request().WithContext(ctx))

			return next(c)
		}
	}
}
