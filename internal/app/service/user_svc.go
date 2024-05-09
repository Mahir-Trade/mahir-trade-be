package service

import (
	"context"
	"fmt"
	"log/slog"
	"mahir-trade-be/internal/app/models"
	"mahir-trade-be/internal/app/repo/discord"
	"mahir-trade-be/internal/app/repo/postgres"
	"mahir-trade-be/internal/app/service/utils"
	"net/http"
	"os"
	"strings"

	"go.uber.org/dig"
)

type (
	LoginReq struct {
		Identity string `json:"identity" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	JWTData struct {
		Email   string `json:"email"`
		UserID  string `json:"user_id"`
		Usename string `json:"username"`
	}

	UserSvc interface {
		UserRegistration(ctx context.Context, req models.User) (err error)
		UserLogin(ctx context.Context, req LoginReq) (resp models.DefaultResponse, err error)
		AssignRoleDiscordToUser(ctx context.Context, userUUID, usernameDiscord string) (resp models.DefaultResponse, err error)
		RemoveRoleDiscordToUser(ctx context.Context, userUUID, usernameDiscord string) (resp models.DefaultResponse, err error)
		ConnectDiscordAccount(ctx context.Context, code string) (resp models.DefaultResponse, err error)
	}

	UserSvcImpl struct {
		dig.In

		UserRepo    postgres.UserRepo
		DiscordRepo discord.DiscordRepo
	}
)

func NewUserSvc(impl UserSvcImpl) UserSvc {
	return &impl
}

func (u *UserSvcImpl) UserRegistration(ctx context.Context, req models.User) (err error) {

	{
		req.Email = strings.ToLower(strings.TrimSpace(req.Email))
		req.PhoneNumber = strings.TrimSpace(req.PhoneNumber)
		req.Username = strings.ToLower(strings.TrimSpace(req.Username))
		req.Password, err = utils.HashPassword(req.Password)
		if err != nil {
			slog.ErrorContext(ctx, fmt.Sprintf("[service][HashPassword] err : %v", err))
			return err
		}
	}

	_, err = u.UserRepo.CreateUser(ctx, req)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("[service][UserRegistration] err : %v", err.Error()))
		return err
	}

	return nil
}

func (u *UserSvcImpl) UserLogin(ctx context.Context, req LoginReq) (resp models.DefaultResponse, err error) {

	{
		resp.Code = http.StatusOK
		resp.Message = "success"
		req.Identity = strings.ToLower(strings.TrimSpace(req.Identity))
	}

	user, err := u.UserRepo.FindUserByEmailOrUsername(ctx, req.Identity)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("[service][UserLogin][FindUserByEmail] err : %v", err))
		err = fmt.Errorf("invalid email or password")
		resp.Code = http.StatusUnauthorized
		resp.Message = "invalid email or password"
		resp.Error = err
		return
	}

	if err = utils.VerifyPassword(req.Password, user.Password); err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("[service][UserLogin][VerifyPassword] err : %v", err))
		err = fmt.Errorf("invalid email or password")
		resp.Code = http.StatusUnauthorized
		resp.Message = "invalid email or password"
		resp.Error = err
		return
	}

	token, exp, err := utils.Sign(JWTData{
		Email:   user.Email,
		UserID:  user.UUID,
		Usename: user.Username,
	})

	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("[service][UserLogin][Sign] err : %v", err))
		err = fmt.Errorf("internal server error, we will fix it soon")
		resp.Code = http.StatusInternalServerError
		resp.Message = "internal server error, we will fix it soon"
		return
	}

	resp.Data = struct {
		Token  string `json:"token"`
		Expire int64  `json:"expire"`
	}{
		Token:  token,
		Expire: exp,
	}

	return
}

func (u *UserSvcImpl) AssignRoleDiscordToUser(ctx context.Context, userUUID, usernameDiscord string) (resp models.DefaultResponse, err error) {
	{
		resp.Code = http.StatusOK
		resp.Message = "success"
	}

	user, err := u.UserRepo.GetUserByUUID(ctx, userUUID)
	if err != nil {
		resp.Code = http.StatusNotFound
		resp.Message = "user not found"
		resp.Error = err.Error()
		slog.ErrorContext(ctx, fmt.Sprintf("[service][AssignRoleDiscordToUser][GetUserByID] err : %v", err))

		return resp, err
	}

	if user.UserID == 0 {
		resp.Code = http.StatusNotFound
		resp.Message = "user not found"
		resp.Error = "user not found"
		slog.ErrorContext(ctx, fmt.Sprintf("[service][AssignRoleDiscordToUser] err : user not found"))

		return resp, fmt.Errorf("user not found")
	}

	userDiscordReq := discord.GetDiscordUserRequest{
		Username: usernameDiscord,
	}

	userDiscord, err := u.DiscordRepo.GetDiscordUser(userDiscordReq)
	if err != nil {
		resp.Code = http.StatusNotFound
		resp.Message = "discord user not found"
		resp.Error = err.Error()
		slog.ErrorContext(ctx, fmt.Sprintf("[service][AssignRoleDiscordToUser][GetDiscordUser] err : %v", err))

		return resp, err
	}

	addRoleDiscordReq := discord.DiscordRoleRequest{
		UserID: userDiscord.ID,
	}

	err = u.DiscordRepo.AddRoleToMember(addRoleDiscordReq)
	if err != nil {
		resp.Code = http.StatusInternalServerError
		resp.Message = "internal server error"
		resp.Error = err.Error()
		slog.ErrorContext(ctx, fmt.Sprintf("[service][AssignRoleDiscordToUser][AddRoleToMember] err : %v", err))

		return
	}

	resp.Data = userDiscord

	return
}

func (u *UserSvcImpl) RemoveRoleDiscordToUser(ctx context.Context, userUUID, usernameDiscord string) (resp models.DefaultResponse, err error) {
	{
		resp.Code = http.StatusOK
		resp.Message = "success"
	}

	user, err := u.UserRepo.GetUserByUUID(ctx, userUUID)
	if err != nil {
		resp.Code = http.StatusNotFound
		resp.Message = "user not found"
		resp.Error = err.Error()
		slog.ErrorContext(ctx, fmt.Sprintf("[service][RemoveRoleDiscordToUser][GetUserByID] err : %v", err))

		return resp, err
	}

	if user.UserID == 0 {
		resp.Code = http.StatusNotFound
		resp.Message = "user not found"
		resp.Error = "user not found"
		slog.ErrorContext(ctx, fmt.Sprintf("[service][RemoveRoleDiscordToUser] err : user not found"))

		return resp, fmt.Errorf("user not found")
	}

	userDiscordReq := discord.GetDiscordUserRequest{
		Username: usernameDiscord,
	}

	userDiscord, err := u.DiscordRepo.GetDiscordUser(userDiscordReq)
	if err != nil {
		resp.Code = http.StatusNotFound
		resp.Message = "discord user not found"
		resp.Error = err.Error()
		slog.ErrorContext(ctx, fmt.Sprintf("[service][RemoveRoleDiscordToUser][GetDiscordUser] err : %v", err))

		return resp, err
	}

	removeRoleDiscordReq := discord.DiscordRoleRequest{
		UserID: userDiscord.ID,
	}

	err = u.DiscordRepo.RemoveRoleFromMember(removeRoleDiscordReq)
	if err != nil {
		resp.Code = http.StatusInternalServerError
		resp.Message = "internal server error"
		resp.Error = err.Error()
		slog.ErrorContext(ctx, fmt.Sprintf("[service][RemoveRoleDiscordToUser][RemoveRoleFromMember] err : %v", err))

		return
	}

	resp.Data = userDiscord

	return
}

func (u *UserSvcImpl) ConnectDiscordAccount(ctx context.Context, code string) (resp models.DefaultResponse, err error) {
	{
		resp.Code = http.StatusOK
		resp.Message = "success"
	}

	tokenResp, err := u.DiscordRepo.ExchangeCodeForToken(code)
	if err != nil {
		resp.Code = http.StatusBadRequest
		resp.Message = "failed to exchange code"
		resp.Error = err.Error()
		slog.ErrorContext(ctx, fmt.Sprintf("[service][ConnectDiscordAccount][ExchangeCodeForToken] err : %v", err))

		return resp, err
	}

	discordUser, err := u.DiscordRepo.GetUserDataByAccessToken(tokenResp.AccessToken)
	if err != nil {
		resp.Code = http.StatusBadRequest
		resp.Message = "failed to get discord user data"
		resp.Error = err.Error()
		slog.ErrorContext(ctx, fmt.Sprintf("[service][ConnectDiscordAccount][GetUserDataByAccessToken] err : %v", err))

		return resp, err
	}

	// TODO: save user to table discord_accounts

	guildID := os.Getenv("DISCORD_GUILD_ID")
	err = u.DiscordRepo.InviteUserToGuild(discordUser.ID, guildID, tokenResp.AccessToken)
	if err != nil {
		resp.Code = http.StatusBadRequest
		resp.Message = "failed to invite user to server"
		resp.Error = err.Error()
		slog.ErrorContext(ctx, fmt.Sprintf("[service][ConnectDiscordAccount][GetUserDataByAccessToken] err : %v", err))

		return resp, err
	}

	return
}