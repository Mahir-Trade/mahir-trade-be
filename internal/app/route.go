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
		discord.GET("/discord/account", authCtrl.ConnectDiscordAccount)
		discord.POST("/assign-role", middleware.AuthMiddleware(authCtrl.AssignRoleDiscordToUser))

		discord.POST("/remove-role", middleware.AuthMiddleware(authCtrl.RemoveRoleDiscordUser))

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
}
