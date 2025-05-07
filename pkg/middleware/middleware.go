package middleware

import (
	"context"
	"encoding/json"
	"log/slog"
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
		Role     string `json:"role"`
	}

	MiddleWareImpl struct {
		dig.In

		UserRepo           postgres.UserRepo
		AdminRepo          postgres.AdminRepo
		UserMembershipRepo postgres.UserMembershipRepo
	}

	MiddleWare interface {
		AuthAdminOrUser(role ...string) func(next echo.HandlerFunc) echo.HandlerFunc
		CheckMembershipActive() echo.MiddlewareFunc
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
				slog.ErrorContext(ctx, "[Middleware][AuthAdminOrUser] len role is 0")
				errResponse.Message = "Unauthorized"
				return c.JSON(http.StatusUnauthorized, errResponse)
			}

			validRoles := map[string]bool{}
			for _, r := range role {
				if r == "admin" || r == "user" {
					validRoles[r] = true
				} else {
					errResponse.Message = "Unauthorized"
					return c.JSON(http.StatusUnauthorized, errResponse)
				}
			}

			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				slog.ErrorContext(ctx, "[Middleware][AuthAdminOrUser] missing authorization header")
				errResponse.Message = "Missing Authorization header"
				return c.JSON(http.StatusUnauthorized, errResponse)
			}

			splitHeader := strings.Split(authHeader, " ")
			if len(splitHeader) != 2 || splitHeader[0] != "Bearer" {
				slog.ErrorContext(ctx, "[Middleware][AuthAdminOrUser] invalid authorization header format")
				errResponse.Message = "Invalid Authorization header format"
				return c.JSON(http.StatusUnauthorized, errResponse)
			}

			token, err := utils.Verify(splitHeader[1])
			if err != nil || token.Valid() != nil {
				slog.ErrorContext(ctx, "[Middleware][AuthAdminOrUser] token invalid", err)
				errResponse.Message = "Unauthorized"
				if err != nil {
					errResponse.Error = err.Error()
				}
				return c.JSON(http.StatusUnauthorized, errResponse)
			}

			var userCtx UserCtxReq
			bt, err := json.Marshal(token.Data)
			if err != nil {
				slog.ErrorContext(ctx, "[Middleware][AuthAdminOrUser] error marshal token data", err)
				errResponse.Message = "Unauthorized"
				errResponse.Error = err.Error()
				return c.JSON(http.StatusUnauthorized, errResponse)
			}
			if err := json.Unmarshal(bt, &userCtx); err != nil {
				slog.ErrorContext(ctx, "[Middleware][AuthAdminOrUser] error unmarshal token data", err)
				errResponse.Message = "Unauthorized"
				errResponse.Error = err.Error()
				return c.JSON(http.StatusUnauthorized, errResponse)
			}

			if userCtx.UserID == 0 {
				slog.ErrorContext(ctx, "[Middleware][AuthAdminOrUser] empty user id")
				errResponse.Message = "Unauthorized"
				errResponse.Error = "Invalid user id"
				return c.JSON(http.StatusUnauthorized, errResponse)
			}

			if validRoles["admin"] {
				adminData, err := m.AdminRepo.FindByUsername(ctx, userCtx.Username)
				if err == nil {
					userCtx.UserID = adminData.AdminID
					userCtx.Email = adminData.Email
					userCtx.Username = adminData.Username
					userCtx.Role = "admin"
					goto SET_CONTEXT
				}
			}

			if validRoles["user"] {
				userData, err := m.UserRepo.FindUserByEmailOrUsername(ctx, userCtx.Email)
				if err == nil {
					userCtx.UserID = int64(userData.UserID)
					userCtx.Email = userData.Email
					userCtx.Username = userData.Username
					userCtx.Role = "user"
					goto SET_CONTEXT
				}
			}

			slog.ErrorContext(ctx, "[Middleware][AuthAdminOrUser] user role not found")
			errResponse.Message = "Unauthorized"
			errResponse.Error = "Invalid role"
			return c.JSON(http.StatusUnauthorized, errResponse)

		SET_CONTEXT:
			ctx = context.WithValue(ctx, UserData, userCtx)
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}
}

func (m *MiddleWareImpl) CheckMembershipActive() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			errResponse := models.DefaultResponse{
				Code:    http.StatusForbidden,
				Message: "Access denied. Membership required.",
				Data:    struct{}{},
			}

			ctx := c.Request().Context()
			userCtxVal := ctx.Value(UserData)
			if userCtxVal == nil {
				slog.ErrorContext(ctx, "[Middleware][CheckMembershipActive] missing user context")
				errResponse.Message = "Unauthorized"
				return c.JSON(http.StatusUnauthorized, errResponse)
			}

			userCtx, ok := userCtxVal.(UserCtxReq)
			if !ok {
				slog.ErrorContext(ctx, "[Middleware][CheckMembershipActive] invalid user context type")
				errResponse.Message = "Unauthorized"
				return c.JSON(http.StatusUnauthorized, errResponse)
			}

			if userCtx.Role == "admin" {
				return next(c)
			}

			membership, err := m.UserMembershipRepo.GetUserMembershipByUserID(ctx, userCtx.UserID)
			if err != nil {
				slog.ErrorContext(ctx, "[Middleware][CheckMembershipActive] error get membership data", err)
				errResponse.Message = "Internal server error"
				return c.JSON(http.StatusInternalServerError, errResponse)
			}

			if !membership.IsMembershipActive {
				slog.WarnContext(ctx, "[Middleware][CheckMembershipActive] user is not active member", "user_id", userCtx.UserID)
				return c.JSON(http.StatusForbidden, errResponse)
			}

			return next(c)
		}
	}
}
