package discord

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/dig"
)

type (
	DiscordRoleRequest struct {
		UserID string
	}

	GetDiscordUserRequest struct {
		Username string
	}

	DiscordUser struct {
		ID       string
		Email    string
		Username string
	}

	DiscordRepo interface {
		AddRoleToMember(req DiscordRoleRequest) error
		RemoveRoleFromMember(req DiscordRoleRequest) error
		GetDiscordUser(req GetDiscordUserRequest) (user DiscordUser, err error)
	}

	DiscordRepoImpl struct {
		dig.In
	}
)

func NewDiscordRepo(impl DiscordRepoImpl) DiscordRepo {
	return &impl
}

func (d *DiscordRepoImpl) AddRoleToMember(req DiscordRoleRequest) error {
	session, err := discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))
	if err != nil {
		slog.Error("[repo][discord][AddRoleToMember] Error creating Discord session: ", err)
		return err
	}
	defer session.Close()

	guildId := os.Getenv("DISCORD_GUILD_ID")
	roleID := os.Getenv("DISCORD_ROLE_ID")

	err = session.GuildMemberRoleAdd(guildId, req.UserID, roleID)
	if err != nil {
		slog.Error("[repo][discord][AddRoleToMember] Error adding role to member: ", err)
		return err
	}

	return nil
}

func (d *DiscordRepoImpl) RemoveRoleFromMember(req DiscordRoleRequest) error {
	session, err := discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))
	if err != nil {
		slog.Error("[repo][discord][RemoveRoleFromMember] Error creating Discord session: ", err)
		return err
	}
	defer session.Close()

	guildId := os.Getenv("DISCORD_GUILD_ID")
	roleID := os.Getenv("DISCORD_ROLE_ID")

	err = session.GuildMemberRoleRemove(guildId, req.UserID, roleID)
	if err != nil {
		slog.Error("[repo][discord][RemoveRoleFromMember] Error removing role from member: ", err)
		return err
	}

	return nil
}

func (d *DiscordRepoImpl) GetDiscordUser(req GetDiscordUserRequest) (user DiscordUser, err error) {
	session, err := discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))
	if err != nil {
		slog.Error("[repo][discord][GetDiscordUser] Error creating Discord session: ", err)
		return user, err
	}
	defer session.Close()

	guildId := os.Getenv("DISCORD_GUILD_ID")

	users, err := session.GuildMembersSearch(guildId, req.Username, 1)
	if err != nil {
		slog.Error("[repo][discord][GetDiscordUser] Error getting guild members: ", err)
		return user, err
	}

	if len(users) == 0 {
		errMsg := "user not found"
		slog.Error(fmt.Sprintf("[repo][discord][GetDiscordUser] %s", errMsg))
		return user, errors.New(errMsg)
	}

	user.ID = users[0].User.ID
	user.Email = users[0].User.Email
	user.Username = users[0].User.Username

	return user, nil
}
