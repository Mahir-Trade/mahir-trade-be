package middleware

import (
	"mahir-trade-be/internal/app/models"
	"mahir-trade-be/internal/app/service/utils"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		errResponse := models.DefaultResponse{Code: http.StatusUnauthorized}

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

		c.Set("user_id", token.Data.(map[string]interface{})["user_id"])
		c.Set("email", token.Data.(map[string]interface{})["email"])
		c.Set("username", token.Data.(map[string]interface{})["username"])

		return next(c)
	}
}
