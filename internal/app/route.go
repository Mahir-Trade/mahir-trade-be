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
	moduleCtrl controller.ModuleCtrlImpl,
	adminCtrl controller.AdminCtrl,
	middleware middleware.MiddleWare,
	packageCtrl controller.PackageCtrl,
) {
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	base := e.Group("/v1")

	users := base.Group("/users")
	{
		users.POST("/register", authCtrl.UserRegistration)
		users.POST("/login", authCtrl.UserLogin)
		users.GET("/login/google", authCtrl.LoginWithGoogle)
		users.GET("/login/google/callback", authCtrl.CallbackGoogle)
	}

	groups := base.Group("/groups", middleware.AuthAdminOrUser("admin", "user"))
	{
		groups.GET("", groupCtrl.GetGroups)
		groups.GET("/:id", groupCtrl.GetGroupByID)

		groups.POST("", groupCtrl.CreateGroup, middleware.AuthAdminOrUser("admin"))
		groups.PUT("/:id", groupCtrl.UpdateGroup, middleware.AuthAdminOrUser("admin"))
		groups.DELETE("/:id", groupCtrl.DeleteGroup, middleware.AuthAdminOrUser("admin"))
	}

	discord := base.Group("/discord")
	{
		// discord callback
		discord.GET("/account", authCtrl.InviteDiscordUserToGuild)
		discord.GET("/account/add-role", authCtrl.ConnectDiscordAccountAndAssignRole)
		discord.GET("/account/remove-role", authCtrl.ConnectDiscordAccountAndRemoveRole)

		// internal
		discord.POST("/connect-role", authCtrl.AssignRoleDiscordToUser, middleware.AuthAdminOrUser("user"))
		discord.POST("/remove-role", authCtrl.RemoveRoleDiscordUser, middleware.AuthAdminOrUser("user"))
	}

	packageRoute := base.Group("/packages", middleware.AuthAdminOrUser("admin", "user"))
	{
		packageRoute.GET("/:id", packageCtrl.GetPackageByID)
		packageRoute.GET("", packageCtrl.GetPackages)
		packageRoute.POST("", packageCtrl.CreatePackage, middleware.AuthAdminOrUser("admin"))
		packageRoute.PUT("/:id", packageCtrl.UpdatePackage, middleware.AuthAdminOrUser("admin"))
		packageRoute.DELETE("/:id", packageCtrl.DeletePackage, middleware.AuthAdminOrUser("admin"))
	}

	modules := base.Group("/modules", middleware.AuthAdminOrUser("admin", "user"))
	{
		modules.GET("/:module_id", moduleCtrl.GetModuleByID)
		modules.POST("", moduleCtrl.CreateModule, middleware.AuthAdminOrUser("admin"))
		modules.PATCH("/:module_id", moduleCtrl.UpdateModule, middleware.AuthAdminOrUser("admin"))
	}

	admins := base.Group("/admins")
	{
		admins.POST("/register", adminCtrl.AdminRegistration)
		admins.POST("/login", adminCtrl.AdminLogin)
	}
}
