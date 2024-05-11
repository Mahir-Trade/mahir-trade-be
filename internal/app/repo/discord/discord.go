package discord

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log/slog"
	"net/http"
	"net/url"
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

	TokenResponse struct {
		AccessToken      string `json:"access_token"`
		TokenType        string `json:"token_type"`
		ExpiresIn        int    `json:"expires_in"`
		Scope            string `json:"scope"`
		RefreshToken     string `json:"refresh_token"`
		Error            string `json:"error,omitempty"`
		ErrorDescription string `json:"error_description,omitempty"`
	}

	DiscordRepo interface {
		AddRoleToMember(req DiscordRoleRequest) error
		RemoveRoleFromMember(req DiscordRoleRequest) error
		GetDiscordUser(req GetDiscordUserRequest) (user DiscordUser, err error)
		ExchangeCodeForToken(code, redirectURI string) (resp TokenResponse, err error)
		GetUserDataByAccessToken(accessToken string) (user DiscordUser, err error)
		InviteUserToGuild(userId, guildId, accessToken string) (err error)
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

func (d *DiscordRepoImpl) ExchangeCodeForToken(code, redirectURI string) (resp TokenResponse, err error) {
	tokenURL := os.Getenv("DISCORD_BASE_URL") + "/v10/oauth2/token"

	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", redirectURI)

	payload := bytes.NewBufferString(data.Encode())

	req, err := http.NewRequest("POST", tokenURL, payload)
	if err != nil {
		slog.Error("[repo][discord][ExchangeCodeForToken] Error create new request: ", err)
		return resp, err
	}

	clientID := os.Getenv("DISCORD_CLIENT_ID")
	clientSecret := os.Getenv("DISCORD_CLIENT_SECRET")

	req.SetBasicAuth(clientID, clientSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		slog.Error("[repo][discord][ExchangeCodeForToken] Error sending request: ", err)
		return resp, err
	}
	defer res.Body.Close()

	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		slog.Error("[repo][discord][ExchangeCodeForToken] Error reading response body: ", err)
		return resp, err
	}

	var tokenResp TokenResponse
	err = json.Unmarshal(bodyBytes, &tokenResp)
	if err != nil {
		slog.Error("[repo][discord][ExchangeCodeForToken] Error unmarshal response body: ", err)
		return resp, err
	}

	return tokenResp, nil
}

func (d *DiscordRepoImpl) GetUserDataByAccessToken(accessToken string) (user DiscordUser, err error) {
	userURL := os.Getenv("DISCORD_BASE_URL") + "/users/@me"

	req, err := http.NewRequest("GET", userURL, nil)
	if err != nil {
		slog.Error("[repo][discord][GetUserDataByAccessToken] Error create new request: ", err)
		return user, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		slog.Error("[repo][discord][GetUserDataByAccessToken] Error sending request: ", err)
		return user, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		slog.Error("[repo][discord][GetUserDataByAccessToken] Error reading response body: ", err)
		return user, err
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		slog.Error("[repo][discord][GetUserDataByAccessToken] Error unmarshal response body: ", err)
		return user, err
	}

	return user, nil
}

func (d *DiscordRepoImpl) InviteUserToGuild(userId, guildId, accessToken string) (err error) {
	inviteGuildURL := os.Getenv("DISCORD_BASE_URL") + fmt.Sprintf("/guilds/%s/members/%s", guildId, userId)

	data := map[string]string{
		"access_token": accessToken,
	}

	payload, err := json.Marshal(data)
	if err != nil {
		slog.Error("[repo][discord][InviteUserToGuild] Error marshaling JSON:", err)
		return err
	}

	req, err := http.NewRequest("PUT", inviteGuildURL, bytes.NewBuffer(payload))
	if err != nil {
		slog.Error("[repo][discord][InviteUserToGuild] Error create new request: ", err)
		return err
	}

	botToken := os.Getenv("DISCORD_TOKEN")

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bot "+botToken)

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		slog.Error("[repo][discord][InviteUserToGuild] Error sending request: ", err)
		return err
	}

	return nil
}
