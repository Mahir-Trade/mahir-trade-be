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
	packageCtrl controller.PackageCtrl,
	reportCtrl controller.ReportCtrl,
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

	packageRoute := base.Group("/packages")
	{
		packageRoute.POST("", packageCtrl.CreatePackage)
		packageRoute.GET("/:id", packageCtrl.GetPackageByID)
		packageRoute.GET("", packageCtrl.GetPackages)
		packageRoute.PUT("/:id", packageCtrl.UpdatePackage)
		packageRoute.DELETE("/:id", packageCtrl.DeletePackage)
	}

	report := base.Group("/reports")
	{
		report.POST("", reportCtrl.CreateReport)
		report.GET("/:id", reportCtrl.GetReportByID)
		report.GET("", reportCtrl.GetReports)
		report.PUT("/:id", reportCtrl.UpdateReport)
		report.DELETE("/:id", reportCtrl.DeleteReport)
	}
}
