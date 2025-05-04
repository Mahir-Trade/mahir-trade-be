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
	subModuleCtrl controller.SubModuleCtrl,
	reportCtrl controller.ReportCtrl,
	paymentCtrl controller.PaymentCtrl,
	cronCtrl controller.SchedulerCtrl,
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
		users.GET("/login/google-callback", authCtrl.CallbackGoogle)
		users.POST("/forgot-password", authCtrl.ForgotPassword)
		users.POST("/reset-password", authCtrl.RequestResetPassword)
		users.POST("/verify-email", authCtrl.SetUserVerified)

		// Dashboard
		users.GET("/detail", authCtrl.GetDetailUser, middleware.AuthAdminOrUser("user"))
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
		discord.POST("/connect-role", authCtrl.AssignRoleDiscordToUser, middleware.AuthAdminOrUser("user"), middleware.CheckMembershipActive())
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

	modules := base.Group("/modules", middleware.AuthAdminOrUser("admin", "user"), middleware.CheckMembershipActive())
	{
		modules.GET("", moduleCtrl.GetModules)
		modules.GET("/:module_id", moduleCtrl.GetModuleByID)
		modules.GET("/group/:group_id", moduleCtrl.GetModulesByGroupID)
		modules.POST("", moduleCtrl.CreateModule, middleware.AuthAdminOrUser("admin"))
		modules.PATCH("/:module_id", moduleCtrl.UpdateModule, middleware.AuthAdminOrUser("admin"))
		modules.DELETE("/:module_id", moduleCtrl.DeleteModule, middleware.AuthAdminOrUser("admin"))
		modules.GET("/user/:module_id", moduleCtrl.GetPercetangeMarkWatchedModulesUser, middleware.AuthAdminOrUser("user"))
	}

	subModules := base.Group("/sub-modules", middleware.AuthAdminOrUser("admin", "user"), middleware.CheckMembershipActive())
	{
		subModules.GET("/:sub_module_id", subModuleCtrl.GetSubModuleByID)
		subModules.GET("", subModuleCtrl.GetSubModules)
		subModules.GET("/module/:module_id", subModuleCtrl.GetSubModulesByModuleID)
		subModules.POST("", subModuleCtrl.CreateSubModule, middleware.AuthAdminOrUser("admin"))
		subModules.PATCH("/:sub_module_id", subModuleCtrl.UpdateSubModule, middleware.AuthAdminOrUser("admin"))
		subModules.DELETE("/:sub_module_id", subModuleCtrl.SoftDeleteSubModule, middleware.AuthAdminOrUser("admin"))
		subModules.POST("/mark-watched", subModuleCtrl.MarkSubModuleAsWatched)
	}

	admins := base.Group("/admins")
	{
		admins.POST("/register", adminCtrl.AdminRegistration)
		admins.POST("/login", adminCtrl.AdminLogin)

		admins.GET("/detail", adminCtrl.GetDetailAdminInfo, middleware.AuthAdminOrUser("admin"))
		admins.GET("/user-detail/:user_id", adminCtrl.GetDetailUserForBO, middleware.AuthAdminOrUser("admin"))
		admins.GET("/users", adminCtrl.GetAllUsers, middleware.AuthAdminOrUser("admin"))

		admins.POST("/users/toggle-expired", adminCtrl.ToggleInactiveUserMembership, middleware.AuthAdminOrUser("admin"))
		admins.POST("/start-membership-program", adminCtrl.StartMembershipProgram, middleware.AuthAdminOrUser("admin"))
	}

	report := base.Group("/reports")
	{
		report.GET("/:id", reportCtrl.GetReportByID)
		report.GET("", reportCtrl.GetReports)
		report.POST("", reportCtrl.CreateReport, middleware.AuthAdminOrUser("admin"))
		report.PUT("/:id", reportCtrl.UpdateReport, middleware.AuthAdminOrUser("admin"))
		report.DELETE("/:id", reportCtrl.DeleteReport, middleware.AuthAdminOrUser("admin"))
	}

	order := base.Group("/orders", middleware.AuthAdminOrUser("admin", "user"))
	{
		order.POST("/create", paymentCtrl.CreatePayment, middleware.AuthAdminOrUser("user"))

	}

	base.POST("/upload", subModuleCtrl.UploadFile, middleware.AuthAdminOrUser("admin"))
	base.POST("/upload/content", reportCtrl.UploadContent, middleware.AuthAdminOrUser("admin"))

	// Public route
	{
		base.POST("/payment-link-callback", paymentCtrl.PaymentLinkCallback)
		base.GET("/cron", paymentCtrl.PaymentLinkCallback)
		base.GET("/package-list", packageCtrl.GetPackages)
	}
}
