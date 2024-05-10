package app

import (
	"mahir-trade-be/internal/app/controller"
	"mahir-trade-be/pkg/middleware"
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

	discord := base.Group("/discord")
	{
		// discord callback
		discord.GET("/account", authCtrl.InviteDiscordUserToGuild)
		discord.GET("/account/add-role", authCtrl.ConnectDiscordAccountAndAssignRole)
		discord.GET("/account/remove-role", authCtrl.ConnectDiscordAccountAndRemoveRole)

		// internal
		discord.POST("/connect-role", middleware.AuthMiddleware(authCtrl.AssignRoleDiscordToUser))
		discord.POST("/remove-role", middleware.AuthMiddleware(authCtrl.RemoveRoleDiscordUser))
	}
}
