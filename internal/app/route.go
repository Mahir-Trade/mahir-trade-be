package app

import (
	"log/slog"
	"mahir-trade-be/internal/app/controller"
	"mahir-trade-be/pkg/middleware"
	"net/http"

	"github.com/bwmarrin/discordgo"
	"github.com/labstack/echo/v4"
)

func setRoute(
	e *echo.Echo,

	authCtrl controller.AuthCtrl,
	groupCtrl controller.GroupCtrl,
	moduleCtrl controller.ModuleCtrlImpl,
	adminCtrl controller.AdminCtrl,
	middleware middleware.MiddleWare,
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
		discord.Use(middleware.AuthAdminOrUser("user"))
		discord.POST("/assign-role", authCtrl.AssignRoleDiscordToUser)
		discord.POST("/remove-role", authCtrl.RemoveRoleDiscordUser)

		// add endpoint to get user discord
		discord.GET("/user", func(c echo.Context) error {
			session, err := discordgo.New("Bot " + "MTIzNTk5NzI1OTg4ODU5MDkxOA.GwqhuY.3MtFKt5OdaBRNOxxemxvAd7x6nXjNor5i-IovM")
			if err != nil {
				slog.Error("Error creating Discord session: ", err)
				return err
			}
			defer session.Close()

			users, err := session.GuildMembersSearch("688746990196228256", "a", 1)
			if err != nil {
				slog.Error("Error getting guild members: ", err)
				return err
			}

			return c.JSON(http.StatusOK, users)
		})
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
