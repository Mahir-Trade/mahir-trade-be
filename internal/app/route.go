package app

import (
	"mahir-trade-be/internal/app/controller"
	"net/http"

	"github.com/labstack/echo/v4"
)

func setRoute(
	e *echo.Echo,

	authCtrl controller.AuthCtrl,
	groupCtrl controller.GroupCtrl,
) {
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	base := e.Group("/v1")

	users := base.Group("/users")
	{
		users.POST("/register", authCtrl.UserRegistration)
		users.POST("/login", authCtrl.UserLogin)
	}

	groups := base.Group("/groups")
	{
		groups.POST("", groupCtrl.CreateGroup)
		groups.GET("/:id", groupCtrl.GetGroupByID)
		groups.GET("", groupCtrl.GetGroups)
		groups.PUT("/:id", groupCtrl.UpdateGroup)
		groups.DELETE("/:id", groupCtrl.DeleteGroup)
	}
}
