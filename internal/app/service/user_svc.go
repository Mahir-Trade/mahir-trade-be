package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"mahir-trade-be/internal/app/models"
	"mahir-trade-be/internal/app/repo/discord"
	"mahir-trade-be/internal/app/repo/google"
	"mahir-trade-be/internal/app/repo/postgres"
	"mahir-trade-be/internal/app/repo/sendgrid"
	"mahir-trade-be/internal/app/service/utils"
	"mahir-trade-be/pkg/middleware"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"go.uber.org/dig"
)

type (
	LoginReq struct {
		Identity string `json:"identity" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	GoogleLoginReq struct {
		State string `json:"state" validate:"required"`
		Code  string `json:"code" validate:"required"`
	}

	ForgotPasswordReq struct {
		Email string `json:"email" validate:"required,email"`
	}

	JWTData struct {
		Email   string `json:"email"`
		UserID  int64  `json:"user_id"`
		Usename string `json:"username"`
	}

	LoginResponse struct {
		Token  string    `json:"token"`
		Expire time.Time `json:"expire"`
	}

	ResetPasswordRequest struct {
		QueryOne             string `json:"q1" validate:"required"`
		QueryTwo             string `json:"q2" validate:"required"`
		Password             string `json:"password" validate:"required,min=8,max=20,uppercase,lowercase,number,specialchar"`
		PasswordConfirmation string `json:"password_confirmation" validate:"required,min=8,max=20,uppercase,lowercase,number,specialchar"`
	}

	UserVerificationReq struct {
		Email       string `json:"email" validate:"required"`
		ExpiredTime string `json:"expired_time" validate:"required"`
	}

	UserSvc interface {
		UserRegistration(ctx context.Context, req models.UserRegistrationRequest) (resp models.DefaultResponse, err error)
		UserLogin(ctx context.Context, req LoginReq) (resp models.DefaultResponse, err error)
		AssignRoleDiscordToUser(ctx context.Context, code string) (resp models.DefaultResponse, err error)
		RemoveRoleDiscordToUser(ctx context.Context) (resp models.DefaultResponse, err error)
		InviteDiscordUserToGuild(ctx context.Context, code, redirectURI string) (resp models.DefaultResponse, err error)
		ConnectDiscordAccountAndAssignRole(ctx context.Context, code string) (resp models.DefaultResponse, err error)
		ConnectDiscordAccountAndRemoveRole(ctx context.Context, code string) (resp models.DefaultResponse, err error)
		LoginWithGoogle(ctx context.Context) (url string, err error)
		CallbackGoogle(ctx context.Context, req GoogleLoginReq) (resp LoginResponse, err error)
		GetDetailUser(ctx context.Context) (resp models.DefaultResponse, err error)
		GetDetailUserForBO(ctx context.Context, userID int64) (resp models.DefaultResponse, err error)
		UpdateMembership(ctx context.Context) (err error)
		ForgotPasswordUser(ctx context.Context, email string) (resp models.DefaultResponse, err error)
		ResetPasswordUser(ctx context.Context, req ResetPasswordRequest) (resp models.DefaultResponse, err error)
		UserVerification(ctx context.Context, req UserVerificationReq) (resp models.DefaultResponse, err error)
	}

	UserSvcImpl struct {
		dig.In

		UserRepo           postgres.UserRepo
		DiscordRepo        discord.DiscordRepo
		DiscordAccountrepo postgres.DiscordAccountRepo
		GoogleRepo         google.GoogleRepo
		UserMembershipRepo postgres.UserMembershipRepo
		EmailTemplateRepo  postgres.EmailTemplateRepo
		SendgridRepo       sendgrid.SendgridRepo
	}
)

func NewUserSvc(impl UserSvcImpl) UserSvc {
	return &impl
}

func (u *UserSvcImpl) UserRegistration(ctx context.Context, req models.UserRegistrationRequest) (resp models.DefaultResponse, err error) {

	const EMAIL_VERIFICATION = "email_verification"
	{
		resp.Code = http.StatusCreated
		resp.Message = "please check your email for verification"
	}

	if req.Password != req.PasswordConfirmation {
		resp.Code = http.StatusBadRequest
		resp.Message = "password and password confirmation must be same"
		resp.Error = errors.New(resp.Message).Error()

		return resp, errors.New(resp.Message)
	}

	{
		req.Email = strings.ToLower(strings.TrimSpace(req.Email))
		req.PhoneNumber = strings.TrimSpace(req.PhoneNumber)
		req.Username = strings.ToLower(strings.TrimSpace(req.Username))
		req.Password, err = utils.HashPassword(req.Password)
		if err != nil {
			slog.ErrorContext(ctx, fmt.Sprintf("[service][HashPassword] err : %v", err))
			resp.Code = http.StatusInternalServerError
			resp.Message = utils.ErrorInternalServer
			resp.Error = err.Error()

			return resp, err
		}
	}

	_, err = u.UserRepo.CreateUser(ctx, models.User{
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		Username:    req.Username,
		Password:    req.Password,
	})
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("[service][UserRegistration] err : %v", err.Error()))
		resp.Code = http.StatusBadRequest
		resp.Message = utils.ErrorBadRequest
		resp.Error = err.Error()

		return resp, err
	}

	// Send Email Confirmation
	encryptedEmail, err := utils.EncryptAES256CBC(req.Email, os.Getenv("JWT_ENCRYPT_KEY"), os.Getenv("JWT_ENCRYPT_IV"))
	if err != nil {
		resp.Code = http.StatusInternalServerError
		resp.Message = "internal server error"
		resp.Error = err.Error()
		slog.ErrorContext(ctx, fmt.Sprintf("[service][UserRegistration][EncryptAES256CBC] err : %v", err))

		return
	}

	expirationVerification := time.Now().Add(time.Hour * 4).Format(time.RFC3339)
	encryptedExpiration, err := utils.EncryptAES256CBC(expirationVerification, os.Getenv("JWT_ENCRYPT_KEY"), os.Getenv("JWT_ENCRYPT_IV"))
	if err != nil {
		resp.Code = http.StatusInternalServerError
		resp.Message = "internal server error"
		resp.Error = err.Error()
		slog.ErrorContext(ctx, fmt.Sprintf("[service][UserRegistration][EncryptAES256CBC] err : %v", err))

		return
	}

	encryptedEmail = url.QueryEscape(encryptedEmail)
	encryptedExpiration = url.QueryEscape(encryptedExpiration)

	values := map[string]interface{}{
		"username":    req.Username,
		"expiredLink": expirationVerification,
		"verifyLink":  fmt.Sprintf("%s?verified=true&email=%s&expiredTime=%s", os.Getenv("EMAIL_VERIFICATION_URL"), encryptedEmail, encryptedExpiration),
	}

	emailTemplate, err := u.EmailTemplateRepo.GetByKey(ctx, EMAIL_VERIFICATION)
	if err != nil {
		resp.Code = http.StatusInternalServerError
		resp.Message = "internal server error"
		resp.Error = err.Error()
		slog.ErrorContext(ctx, fmt.Sprintf("[service][UserRegistration][GetByKey] err : %v", err))

		return
	}

	parsedHtml, err := utils.MappingValuesToTemplate(values, emailTemplate)
	if err != nil {
		resp.Code = http.StatusInternalServerError
		resp.Message = "internal server error"
		resp.Error = err.Error()
		slog.ErrorContext(ctx, fmt.Sprintf("[service][UserRegistration][MappingValuesToTemplate] err : %v", err))

		return
	}

	sendgridReq := models.SendgridSendEmailRequest{
		From:          os.Getenv("SENDGRID_SENDER_EMAIL"),
		To:            req.Email,
		Subject:       "Mahir Trade Email Verification",
		Body:          parsedHtml,
		SenderName:    os.Getenv("SENDGRID_SENDER_NAME"),
		RecepientName: req.Username,
	}

	go func() {
		err = u.SendgridRepo.SendEmail(ctx, sendgridReq)
		if err != nil {
			slog.ErrorContext(ctx, fmt.Sprintf("[service][UserRegistration][SendEmail] err : %v", err))
		}
	}()

	return resp, nil
}

func (u *UserSvcImpl) UserLogin(ctx context.Context, req LoginReq) (resp models.DefaultResponse, err error) {

	{
		resp.Code = http.StatusOK
		resp.Message = "success"
		resp.Data = struct{}{}
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

	if user.VerifiedAt == nil {
		slog.ErrorContext(ctx, fmt.Sprintf("[service][UserLogin][VerifiedAt] err : %v", err))
		err = fmt.Errorf("please verify your email first")
		resp.Code = http.StatusUnauthorized
		resp.Message = "please verify your email first"
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
		UserID:  user.UserID,
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
		Token  string    `json:"token"`
		Expire time.Time `json:"expire"`
	}{
		Token:  token,
		Expire: exp,
	}

	return
}

func (u *UserSvcImpl) AssignRoleDiscordToUser(ctx context.Context, code string) (resp models.DefaultResponse, err error) {
	{
		resp.Code = http.StatusOK
		resp.Message = "success"
		resp.Data = struct{}{}
	}

	userData, ok := ctx.Value(middleware.UserData).(middleware.UserCtxReq)
	if !ok {
		resp.Code = http.StatusBadRequest
		resp.Message = "bad request"
		resp.Error = errors.New("something went wrong, we will fix it soon").Error()
		return
	}

	redirectURI := os.Getenv("DISCORD_REDIRECT_URI_ASSIGN_ROLE")
	resp, err = u.InviteDiscordUserToGuild(ctx, code, redirectURI)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("[service][AssignRoleDiscordToUser][InviteDiscordUserToGuild] err : %v", err))
		return resp, err
	}

	userMembership, err := u.UserMembershipRepo.GetUserMembershipByUserID(ctx, int64(userData.UserID))
	if err != nil {
		resp.Code = http.StatusBadRequest
		resp.Message = "something went wrong"
		resp.Error = err.Error()
		slog.ErrorContext(ctx, fmt.Sprintf("[service][AssignRoleDiscordToUser][GetUserMembershipByUserID] err : %v", err))

		return
	}

	if !userMembership.IsMembershipActive || userMembership.ID == 0 {
		resp.Code = http.StatusBadRequest
		resp.Message = "user membership is not active"
		resp.Error = errors.New("user membership is not active").Error()
		slog.ErrorContext(ctx, fmt.Sprintf("[service][AssignRoleDiscordToUser][GetUserMembershipByUserID] err : %v", err))

		return
	}

	discordAccount, err := u.DiscordAccountrepo.GetDiscordAccountByUserID(ctx, userData.UserID)
	if err != nil {
		resp.Code = http.StatusNotFound
		resp.Message = "something went wrong"
		resp.Error = err.Error()
		slog.ErrorContext(ctx, fmt.Sprintf("[service][AssignRoleDiscordToUser][GetDiscordAccountByUserID] err : %v", err))

		return
	}

	if discordAccount.ID != 0 {
		resp.Code = http.StatusBadRequest
		resp.Message = "discord account already registered"
		resp.Error = errors.New("discord account already registered").Error()
		slog.ErrorContext(ctx, fmt.Sprintf("[service][AssignRoleDiscordToUser][GetDiscordAccountByUserID] err: %v, userID: %d", resp.Error, userData.UserID))

		return
	}

	addRoleDiscordReq := discord.DiscordRoleRequest{UserID: resp.Data.(discord.DiscordUser).ID}
	err = u.DiscordRepo.AddRoleToMember(addRoleDiscordReq)
	if err != nil {
		resp.Code = http.StatusInternalServerError
		resp.Message = "internal server error"
		resp.Error = err.Error()
		slog.ErrorContext(ctx, fmt.Sprintf("[service][AssignRoleDiscordToUser][AddRoleToMember] err : %v", err))

		return
	}

	_, err = u.DiscordAccountrepo.CreateDiscordAccount(ctx, models.DiscordAccount{
		UserID:           int64(userData.UserID),
		DiscordAccountID: resp.Data.(discord.DiscordUser).ID,
		Username:         resp.Data.(discord.DiscordUser).Username,
		Email:            resp.Data.(discord.DiscordUser).Email,
	})
	if err != nil {
		resp.Code = http.StatusInternalServerError
		resp.Message = "internal server error"
		resp.Error = err.Error()
		slog.ErrorContext(ctx, fmt.Sprintf("[service][AssignRoleDiscordToUser][CreateDiscordAccount] err : %v", err))

		return
	}

	return
}

func (u *UserSvcImpl) RemoveRoleDiscordToUser(ctx context.Context) (resp models.DefaultResponse, err error) {
	{
		resp.Code = http.StatusOK
		resp.Message = "success"
		resp.Data = struct{}{}
	}

	userData, ok := ctx.Value(middleware.UserData).(middleware.UserCtxReq)
	if !ok {
		resp.Code = http.StatusBadRequest
		resp.Message = "bad request"
		resp.Error = errors.New("something went wrong, we will fix it soon").Error()
		return
	}

	discordAccount, err := u.DiscordAccountrepo.GetDiscordAccountByUserID(ctx, int64(userData.UserID))
	if err != nil || discordAccount.ID == 0 {
		resp.Code = http.StatusNotFound
		resp.Message = "discord account not found"
		resp.Error = err.Error()
		slog.ErrorContext(ctx, fmt.Sprintf("[service][RemoveRoleDiscordToUser][GetDiscordAccountByUserID] err : %v", err))

		return resp, err
	}

	removeRoleDiscordReq := discord.DiscordRoleRequest{
		UserID: discordAccount.DiscordAccountID,
	}

	err = u.DiscordRepo.RemoveRoleFromMember(removeRoleDiscordReq)
	if err != nil {
		resp.Code = http.StatusInternalServerError
		resp.Message = "internal server error"
		resp.Error = err.Error()
		slog.ErrorContext(ctx, fmt.Sprintf("[service][RemoveRoleDiscordToUser][RemoveRoleFromMember] err : %v", err))

		return
	}

	err = u.DiscordAccountrepo.DeleteDiscordAccountByUserID(ctx, discordAccount.ID)
	if err != nil {
		resp.Code = http.StatusInternalServerError
		resp.Message = "internal server error"
		resp.Error = err.Error()
		slog.ErrorContext(ctx, fmt.Sprintf("[service][RemoveRoleDiscordToUser][DeleteDiscordAccountByUserID] err : %v", err))

		return
	}

	resp.Data = discordAccount

	return
}

func (u *UserSvcImpl) InviteDiscordUserToGuild(ctx context.Context, code, redirectURI string) (resp models.DefaultResponse, err error) {
	{
		resp.Code = http.StatusOK
		resp.Message = "success"
		resp.Data = struct{}{}
	}

	tokenResp, err := u.DiscordRepo.ExchangeCodeForToken(code, redirectURI)
	if err != nil {
		resp.Code = http.StatusBadRequest
		resp.Message = "failed to exchange code"
		resp.Error = err.Error()
		slog.ErrorContext(ctx, fmt.Sprintf("[service][ConnectDiscordAccount][ExchangeCodeForToken] err : %v", err))

		return resp, err
	}

	if tokenResp.Error != "" {
		resp.Code = http.StatusBadRequest
		resp.Message = tokenResp.Error
		resp.Error = errors.New(tokenResp.ErrorDescription).Error()
		err = errors.New(tokenResp.ErrorDescription)
		slog.ErrorContext(ctx, fmt.Sprintf("[service][ConnectDiscordAccount][ExchangeCodeForToken] err : %v", tokenResp.ErrorDescription))

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

	guildID := os.Getenv("DISCORD_GUILD_ID")
	err = u.DiscordRepo.InviteUserToGuild(discordUser.ID, guildID, tokenResp.AccessToken)
	if err != nil {
		resp.Code = http.StatusBadRequest
		resp.Message = "failed to invite user to server"
		resp.Error = err.Error()
		slog.ErrorContext(ctx, fmt.Sprintf("[service][ConnectDiscordAccount][GetUserDataByAccessToken] err : %v", err))

		return resp, err
	}

	resp.Data = discordUser

	return
}

func (u *UserSvcImpl) ConnectDiscordAccountAndAssignRole(ctx context.Context, code string) (resp models.DefaultResponse, err error) {
	{
		resp.Code = http.StatusOK
		resp.Message = "success"
		resp.Data = struct{}{}
	}

	redirectURI := os.Getenv("DISCORD_REDIRECT_URI_ASSIGN_ROLE")
	tokenResp, err := u.DiscordRepo.ExchangeCodeForToken(code, redirectURI)
	if err != nil {
		resp.Code = http.StatusBadRequest
		resp.Message = "failed to exchange code"
		resp.Error = err.Error()
		slog.ErrorContext(ctx, fmt.Sprintf("[service][ConnectDiscordAccountAndAssignRole][ExchangeCodeForToken] err : %v", err))

		return resp, err
	}

	if tokenResp.Error != "" {
		resp.Code = http.StatusBadRequest
		resp.Message = tokenResp.Error
		resp.Error = tokenResp.ErrorDescription
		slog.ErrorContext(ctx, fmt.Sprintf("[service][ConnectDiscordAccountAndAssignRole][ExchangeCodeForToken] err : %v", err))

		return resp, err
	}

	discordUser, err := u.DiscordRepo.GetUserDataByAccessToken(tokenResp.AccessToken)
	if err != nil {
		resp.Code = http.StatusBadRequest
		resp.Message = "failed to get discord user data"
		resp.Error = err.Error()
		slog.ErrorContext(ctx, fmt.Sprintf("[service][ConnectDiscordAccountAndAssignRole][GetUserDataByAccessToken] err : %v", err))

		return resp, err
	}

	guildID := os.Getenv("DISCORD_GUILD_ID")
	err = u.DiscordRepo.InviteUserToGuild(discordUser.ID, guildID, tokenResp.AccessToken)
	if err != nil {
		resp.Code = http.StatusBadRequest
		resp.Message = "failed to invite user to server"
		resp.Error = err.Error()
		slog.ErrorContext(ctx, fmt.Sprintf("[service][ConnectDiscordAccountAndAssignRole][InviteUserToGuild] err : %v", err))

		return resp, err
	}

	addRoleDiscordReq := discord.DiscordRoleRequest{UserID: discordUser.ID}

	err = u.DiscordRepo.AddRoleToMember(addRoleDiscordReq)
	if err != nil {
		resp.Code = http.StatusInternalServerError
		resp.Message = "internal server error"
		resp.Error = err.Error()
		slog.ErrorContext(ctx, fmt.Sprintf("[service][ConnectDiscordAccountAndAssignRole][AddRoleToMember] err : %v", err))

		return
	}

	return
}

func (u *UserSvcImpl) ConnectDiscordAccountAndRemoveRole(ctx context.Context, code string) (resp models.DefaultResponse, err error) {
	{
		resp.Code = http.StatusOK
		resp.Message = "success"
		resp.Data = struct{}{}
	}

	redirectURI := os.Getenv("DISCORD_REDIRECT_URI_REMOVE_ROLE")
	tokenResp, err := u.DiscordRepo.ExchangeCodeForToken(code, redirectURI)
	if err != nil {
		resp.Code = http.StatusBadRequest
		resp.Message = "failed to exchange code"
		resp.Error = err.Error()
		slog.ErrorContext(ctx, fmt.Sprintf("[service][ConnectDiscordAccountAndRemoveRole][ExchangeCodeForToken] err : %v", err))

		return resp, err
	}

	if tokenResp.Error != "" {
		resp.Code = http.StatusBadRequest
		resp.Message = tokenResp.Error
		resp.Error = tokenResp.ErrorDescription
		slog.ErrorContext(ctx, fmt.Sprintf("[service][ConnectDiscordAccountAndRemoveRole][ExchangeCodeForToken] err : %v", err))

		return resp, err
	}

	discordUser, err := u.DiscordRepo.GetUserDataByAccessToken(tokenResp.AccessToken)
	if err != nil {
		resp.Code = http.StatusBadRequest
		resp.Message = "failed to get discord user data"
		resp.Error = err.Error()
		slog.ErrorContext(ctx, fmt.Sprintf("[service][ConnectDiscordAccountAndRemoveRole][GetUserDataByAccessToken] err : %v", err))

		return resp, err
	}

	guildID := os.Getenv("DISCORD_GUILD_ID")
	err = u.DiscordRepo.InviteUserToGuild(discordUser.ID, guildID, tokenResp.AccessToken)
	if err != nil {
		resp.Code = http.StatusBadRequest
		resp.Message = "failed to invite user to server"
		resp.Error = err.Error()
		slog.ErrorContext(ctx, fmt.Sprintf("[service][ConnectDiscordAccountAndRemoveRole][InviteUserToGuild] err : %v", err))

		return resp, err
	}

	removeRoleDiscordReq := discord.DiscordRoleRequest{UserID: discordUser.ID}

	err = u.DiscordRepo.RemoveRoleFromMember(removeRoleDiscordReq)
	if err != nil {
		resp.Code = http.StatusInternalServerError
		resp.Message = "internal server error"
		resp.Error = err.Error()
		slog.ErrorContext(ctx, fmt.Sprintf("[service][ConnectDiscordAccountAndRemoveRole][RemoveRoleFromMember] err : %v", err))

		return
	}

	return
}

func (u *UserSvcImpl) LoginWithGoogle(ctx context.Context) (url string, err error) {
	url, err = u.GoogleRepo.Login(ctx)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("[service][LoginWithGoogle] err : %v", err))
		return
	}
	slog.InfoContext(ctx, fmt.Sprintf("[service][LoginWithGoogle] url : %v", url))
	return
}

func (u *UserSvcImpl) CallbackGoogle(ctx context.Context, req GoogleLoginReq) (resp LoginResponse, err error) {

	decodeString, err := url.QueryUnescape(req.Code)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("[service][CallbackGoogle][QueryUnescape] err : %v", err))
		return
	}

	userInfo, err := u.GoogleRepo.Callback(ctx, google.GoogleCallbackRequest{
		State: req.State,
		Code:  decodeString,
	})
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("[service][CallbackGoogle] err : %v", err))
		return
	}

	userEmail := strings.ToLower(strings.TrimSpace(userInfo.Email))
	userName := strings.ToLower(strings.TrimSpace(userInfo.Name))

	slog.InfoContext(ctx, fmt.Sprintf("[service][CallbackGoogle] userEmail : %v, userName : %v", userEmail, userName))

	user, err := u.UserRepo.FindUserByEmailOrUsername(ctx, userEmail)
	if err != nil {
		slog.InfoContext(ctx, "[service][CallbackGoogle][FindUserByEmailOrUsername], Try to create new user")
	}

	if user.UserID != 0 {
		token, exp, errSign := utils.Sign(JWTData{
			Email:   user.Email,
			UserID:  user.UserID,
			Usename: user.Username,
		})
		if errSign != nil {
			slog.ErrorContext(ctx, fmt.Sprintf("[service][CallbackGoogle][Sign] err : %v", errSign))
			return
		}
		resp = LoginResponse{
			Token:  token,
			Expire: exp,
		}

		return
	}

	passwordHash, err := utils.HashPassword(utils.GenerateRandomPassword(12))
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("[service][CallbackGoogle][HashPassword] err : %v", err))
		return
	}

	userReq := models.User{
		Email:    userEmail,
		Username: userName,
		Password: passwordHash,
	}

	userId, err := u.UserRepo.CreateUser(ctx, userReq)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("[service][CallbackGoogle][CreateUser] err : %v", err))
		return
	}

	err = u.UserRepo.SetUserVerified(ctx, int64(userId))
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("[service][CallbackGoogle][SetUserVerified] err : %v", err))
		return
	}

	token, exp, errSign := utils.Sign(JWTData{
		Email:   userReq.Email,
		UserID:  int64(userId),
		Usename: userReq.Username,
	})

	if errSign != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("[service][CallbackGoogle][Sign] err : %v", errSign))
		return
	}

	resp = LoginResponse{
		Token:  token,
		Expire: exp,
	}

	return
}

func (u *UserSvcImpl) GetDetailUser(ctx context.Context) (resp models.DefaultResponse, err error) {
	{
		resp.Code = http.StatusOK
		resp.Message = "success"
	}

	var userDetailResp models.UserDetailResponse

	userData, ok := ctx.Value(middleware.UserData).(middleware.UserCtxReq)
	if !ok {
		resp.Code = http.StatusBadRequest
		resp.Message = "bad request"
		resp.Error = errors.New(utils.ErrorInternalServer).Error()
		slog.ErrorContext(ctx, fmt.Sprintf("[service][GetDetailUser] err : %v", errors.New(utils.ErrorInternalServer)))

		return resp, errors.New(utils.ErrorInternalServer)
	}

	user, err := u.UserRepo.GetUserByID(ctx, int64(userData.UserID))
	if err != nil {
		resp.Code = http.StatusNotFound
		resp.Message = "user not found"
		resp.Error = err.Error()
		slog.ErrorContext(ctx, fmt.Sprintf("[service][GetDetailUser][GetUserByID] err : %v", err))

		return resp, err
	}

	userMembership, err := u.UserMembershipRepo.GetUserMembershipByUserID(ctx, int64(userData.UserID))
	if err != nil {
		resp.Code = http.StatusBadRequest
		resp.Message = "something went wrong"
		resp.Error = err.Error()
		slog.ErrorContext(ctx, fmt.Sprintf("[service][GetDetailUser][GetUserMembershipByUserID] err : %v", err))

		return
	}

	discordAccount, err := u.DiscordAccountrepo.GetDiscordAccountByUserID(ctx, int64(userData.UserID))
	if err != nil {
		resp.Code = http.StatusNotFound
		resp.Message = "discord account not found"
		resp.Error = err.Error()
		slog.ErrorContext(ctx, fmt.Sprintf("[service][GetDetailUser][GetDiscordAccountByUserID] err : %v", err))

		return
	}

	userDetailResp = models.UserDetailResponse{
		UserID:              user.UserID,
		UUID:                user.UUID,
		Email:               user.Email,
		PhoneNumber:         user.PhoneNumber,
		Username:            user.Username,
		IsActive:            user.IsActive,
		MembershipExpiredAt: userMembership.ExclusiveExpiredAt,
	}

	if !userMembership.IsMembershipActive || userMembership.ID == 0 {
		userDetailResp.IsMembershipActive = false
	} else {
		userDetailResp.IsMembershipActive = true
	}

	if discordAccount.ID > 0 {
		userDetailResp.DiscordUsername = discordAccount.Username
	}

	resp.Data = userDetailResp

	return
}

func (u *UserSvcImpl) GetDetailUserForBO(ctx context.Context, userID int64) (resp models.DefaultResponse, err error) {
	{
		resp.Code = http.StatusOK
		resp.Message = "success"
	}

	user, err := u.UserRepo.GetUserByID(ctx, userID)
	if err != nil {
		resp.Code = http.StatusNotFound
		resp.Message = "user not found"
		resp.Error = err.Error()
		slog.ErrorContext(ctx, fmt.Sprintf("[service][GetDetailUserForBO][GetUserByID] err : %v", err))

		return resp, err
	}

	user.Password = ""
	resp.Data = user

	return
}

func (u *UserSvcImpl) UpdateMembership(ctx context.Context) (err error) {
	slog.InfoContext(ctx, "[service][UpdateMembership][cron] start update membership")

	userMemberships, err := u.UserMembershipRepo.GetUserMembershipExpired(ctx)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("[service][UpdateMembership][GetUserMembershipExpired] err : %v", err))
		return
	}

	for _, userMembership := range userMemberships {
		err = u.UserMembershipRepo.UpdateUserMembershipByUserID(ctx, models.UserMembership{
			UserID:             userMembership.UserID,
			IsMembershipActive: false,
			ExpiredAt:          userMembership.ExpiredAt,
			Status:             models.MembershipStatusExpired,
			UpdatedBy:          "CRONJOB",
		})
		if err != nil {
			slog.ErrorContext(ctx, fmt.Sprintf("[service][UpdateMembership][UpdateUserMembership] err : %v", err))
			break
		}

		discordAccount, err := u.DiscordAccountrepo.GetDiscordAccountByUserID(ctx, userMembership.UserID)
		if err != nil {
			slog.ErrorContext(ctx, fmt.Sprintf("[service][UpdateMembership][GetDiscordAccountByUserID] err : %v", err))
			break
		}

		if discordAccount.ID == 0 {
			continue
		}

		err = u.DiscordRepo.RemoveRoleFromMember(discord.DiscordRoleRequest{UserID: discordAccount.DiscordAccountID})
		if err != nil {
			slog.ErrorContext(ctx, fmt.Sprintf("[service][UpdateMembership][RemoveRoleFromMember] userID: %d, err : %v", userMembership.UserID, err))
			break
		}

		err = u.DiscordAccountrepo.DeleteDiscordAccountByUserID(ctx, discordAccount.ID)
		if err != nil {
			slog.ErrorContext(ctx, fmt.Sprintf("[service][UpdateMembership][DeleteDiscordAccountByUserID] userID: %d, err : %v", userMembership.UserID, err))
			break
		}
	}

	slog.InfoContext(ctx, "[service][UpdateMembership][cron] finish update membership")
	return
}

func (u *UserSvcImpl) ForgotPasswordUser(ctx context.Context, email string) (resp models.DefaultResponse, err error) {

	const FORGOT_PASSWORD = "forgot_password"
	{
		resp.Code = http.StatusOK
		resp.Message = "success"
		resp.Data = struct{}{}
		resp.Error = struct{}{}
	}

	email = strings.ToLower(strings.TrimSpace(email))
	user, err := u.UserRepo.FindUserByEmailOrUsername(ctx, email)
	if err != nil {
		resp.Code = http.StatusOK
		resp.Message = "please check your email"
		resp.Error = err.Error()
		slog.ErrorContext(ctx, fmt.Sprintf("[service][ForgotPasswordUser][FindUserByEmailOrUsername] err : %v", err))

		return resp, err
	}

	resp.Message = "please check your email"

	encryptedEmail, err := utils.EncryptAES256CBC(email, os.Getenv("JWT_ENCRYPT_KEY"), os.Getenv("JWT_ENCRYPT_IV"))
	if err != nil {
		resp.Code = http.StatusInternalServerError
		resp.Message = "internal server error"
		resp.Error = err.Error()
		slog.ErrorContext(ctx, fmt.Sprintf("[service][ForgotPasswordUser][EncryptAES256CBC] err : %v", err))

		return
	}
	encryptedExpiration, err := utils.EncryptAES256CBC(time.Now().Add(time.Hour*4).Format(time.RFC3339), os.Getenv("JWT_ENCRYPT_KEY"), os.Getenv("JWT_ENCRYPT_IV"))
	if err != nil {
		resp.Code = http.StatusInternalServerError
		resp.Message = "internal server error"
		resp.Error = err.Error()
		slog.ErrorContext(ctx, fmt.Sprintf("[service][ForgotPasswordUser][EncryptAES256CBC] err : %v", err))

		return
	}

	encryptedEmail = url.QueryEscape(encryptedEmail)
	encryptedExpiration = url.QueryEscape(encryptedExpiration)

	values := map[string]interface{}{
		"username":   user.Username,
		"verifyLink": fmt.Sprintf("%s?q1=%s&q2=%s", os.Getenv("FORGOT_PASSWORD_URL"), encryptedEmail, encryptedExpiration),
	}

	emailTemplate, err := u.EmailTemplateRepo.GetByKey(ctx, FORGOT_PASSWORD)
	if err != nil {
		resp.Code = http.StatusInternalServerError
		resp.Message = "internal server error"
		resp.Error = err.Error()
		slog.ErrorContext(ctx, fmt.Sprintf("[service][ForgotPasswordUser][GetByKey] err : %v", err))

		return
	}

	parsedHtml, err := utils.MappingValuesToTemplate(values, emailTemplate)
	if err != nil {
		resp.Code = http.StatusInternalServerError
		resp.Message = "internal server error"
		resp.Error = err.Error()
		slog.ErrorContext(ctx, fmt.Sprintf("[service][ForgotPasswordUser][MappingValuesToTemplate] err : %v", err))

		return
	}

	sendgridReq := models.SendgridSendEmailRequest{
		From:          os.Getenv("SENDGRID_SENDER_EMAIL"),
		To:            email,
		Subject:       "Mahir Trade Password Recovery",
		Body:          parsedHtml,
		SenderName:    os.Getenv("SENDGRID_SENDER_NAME"),
		RecepientName: user.Username,
	}

	go func() {
		err = u.SendgridRepo.SendEmail(ctx, sendgridReq)
		if err != nil {
			slog.ErrorContext(ctx, fmt.Sprintf("[service][ForgotPasswordUser][SendEmail] err : %v", err))
		}
	}()

	return
}

func (u *UserSvcImpl) ResetPasswordUser(ctx context.Context, req ResetPasswordRequest) (resp models.DefaultResponse, err error) {
	{
		resp.Code = http.StatusOK
		resp.Message = "success"
		resp.Data = struct{}{}
	}
	dataEmail, err := url.QueryUnescape(req.QueryOne)
	if err != nil {
		resp.Code = http.StatusBadRequest
		resp.Message = "invalid token"
		resp.Error = err.Error()
		slog.ErrorContext(ctx, fmt.Sprintf("[service][ResetPasswordUser][QueryUnescape][QueryOne] err : %v", err))
		return
	}
	dataExpiration, err := url.QueryUnescape(req.QueryTwo)
	if err != nil {
		resp.Code = http.StatusBadRequest
		resp.Message = "invalid token"
		resp.Error = err.Error()
		slog.ErrorContext(ctx, fmt.Sprintf("[service][ResetPasswordUser][QueryUnescape][QueryTwo] err : %v", err))
		return
	}

	email, err := utils.DecryptAES256CBC(dataEmail, os.Getenv("JWT_ENCRYPT_KEY"), os.Getenv("JWT_ENCRYPT_IV"))
	if err != nil {
		resp.Code = http.StatusBadRequest
		resp.Message = "invalid token"
		resp.Error = err.Error()
		slog.ErrorContext(ctx, fmt.Sprintf("[service][ResetPasswordUser][DecryptAES256CBC] err : %v", err))

		return resp, err
	}

	expiration, err := utils.DecryptAES256CBC(dataExpiration, os.Getenv("JWT_ENCRYPT_KEY"), os.Getenv("JWT_ENCRYPT_IV"))
	if err != nil {
		resp.Code = http.StatusBadRequest
		resp.Message = "invalid token"
		resp.Error = err.Error()
		slog.ErrorContext(ctx, fmt.Sprintf("[service][ResetPasswordUser][DecryptAES256CBC] err : %v", err))

		return resp, err
	}

	timeParsed, err := time.Parse(time.RFC3339, expiration)
	if err != nil {
		resp.Code = http.StatusBadRequest
		resp.Message = "invalid token"
		resp.Error = err.Error()
		slog.ErrorContext(ctx, fmt.Sprintf("[service][ResetPasswordUser][ParseTime] err : %v", err))

		return resp, err
	}

	if time.Now().After(timeParsed) {
		resp.Code = http.StatusBadRequest
		resp.Message = "illegal token or token expired"
		resp.Error = fmt.Errorf("illegal token or token expired").Error()
		slog.ErrorContext(ctx, fmt.Sprintf("[service][ResetPasswordUser][CheckTokenExpired] err : %v", err))

		return resp, err
	}

	user, err := u.UserRepo.FindUserByEmailOrUsername(ctx, email)
	if err != nil {
		resp.Code = http.StatusBadRequest
		resp.Message = "illegal token or token expired"
		resp.Error = err.Error()
		slog.ErrorContext(ctx, fmt.Sprintf("[service][ResetPasswordUser][FindUserByEmailOrUsername] err : %v", err))

		return resp, err
	}

	if req.Password != req.PasswordConfirmation {
		resp.Code = http.StatusBadRequest
		resp.Message = "password and password confirmation must be same"
		resp.Error = errors.New(resp.Message).Error()

		return resp, errors.New(resp.Message)
	}

	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		resp.Code = http.StatusInternalServerError
		resp.Message = "internal server error"
		resp.Error = err.Error()
		slog.ErrorContext(ctx, fmt.Sprintf("[service][ResetPasswordUser][HashPassword] err : %v", err))

		return
	}

	isUpdated, err := u.UserRepo.UpdatePassword(ctx, passwordHash, user.Username, user.UserID)
	if err != nil {
		resp.Code = http.StatusInternalServerError
		resp.Message = "internal server error"
		resp.Error = err.Error()
		slog.ErrorContext(ctx, fmt.Sprintf("[service][ResetPasswordUser][UpdatePassword] err : %v", err))

		return
	}

	if !isUpdated {
		resp.Code = http.StatusInternalServerError
		resp.Message = "internal server error"
		resp.Error = errors.New("internal server error").Error()
		slog.ErrorContext(ctx, fmt.Sprintf("[service][ResetPasswordUser][CheckUpdatePassword] err : %v", err))
		return
	}

	return
}

func (u *UserSvcImpl) UserVerification(ctx context.Context, req UserVerificationReq) (resp models.DefaultResponse, err error) {
	{
		resp.Code = http.StatusOK
		resp.Message = "success"
	}

	dataEmail, err := url.QueryUnescape(req.Email)
	if err != nil {
		resp.Code = http.StatusBadRequest
		resp.Message = "invalid token"
		resp.Error = err.Error()
		slog.ErrorContext(ctx, fmt.Sprintf("[service][UserVerification][QueryUnescape][QueryOne] err : %v", err))
		return
	}

	dataExpiration, err := url.QueryUnescape(req.ExpiredTime)
	if err != nil {
		resp.Code = http.StatusBadRequest
		resp.Message = "invalid token"
		resp.Error = err.Error()
		slog.ErrorContext(ctx, fmt.Sprintf("[service][UserVerification][QueryUnescape][QueryTwo] err : %v", err))
		return
	}

	email, err := utils.DecryptAES256CBC(dataEmail, os.Getenv("JWT_ENCRYPT_KEY"), os.Getenv("JWT_ENCRYPT_IV"))
	if err != nil {
		resp.Code = http.StatusBadRequest
		resp.Message = "invalid token"
		resp.Error = err.Error()
		slog.ErrorContext(ctx, fmt.Sprintf("[service][ResetPasswordUser][DecryptAES256CBC] err : %v", err))

		return resp, err
	}

	expiration, err := utils.DecryptAES256CBC(dataExpiration, os.Getenv("JWT_ENCRYPT_KEY"), os.Getenv("JWT_ENCRYPT_IV"))
	if err != nil {
		resp.Code = http.StatusBadRequest
		resp.Message = "invalid token"
		resp.Error = err.Error()
		slog.ErrorContext(ctx, fmt.Sprintf("[service][ResetPasswordUser][DecryptAES256CBC] err : %v", err))

		return resp, err
	}

	timeParsed, err := time.Parse(time.RFC3339, expiration)
	if err != nil {
		resp.Code = http.StatusBadRequest
		resp.Message = "invalid token"
		resp.Error = err.Error()
		slog.ErrorContext(ctx, fmt.Sprintf("[service][ResetPasswordUser][ParseTime] err : %v", err))

		return resp, err
	}

	if time.Now().After(timeParsed) {
		resp.Code = http.StatusBadRequest
		resp.Message = "illegal token or token expired"
		resp.Error = fmt.Errorf("illegal token or token expired").Error()
		slog.ErrorContext(ctx, fmt.Sprintf("[service][ResetPasswordUser][CheckTokenExpired] err : %v", err))

		return resp, err
	}

	user, err := u.UserRepo.FindUserByEmailOrUsername(ctx, email)
	if err != nil {
		resp.Code = http.StatusBadRequest
		resp.Message = "illegal token or token expired"
		resp.Error = err.Error()
		slog.ErrorContext(ctx, fmt.Sprintf("[service][ResetPasswordUser][FindUserByEmailOrUsername] err : %v", err))

		return resp, err
	}

	if user.VerifiedAt == nil {
		err = u.UserRepo.SetUserVerified(ctx, user.UserID)
		if err != nil {
			resp.Code = http.StatusInternalServerError
			resp.Message = "internal server error"
			resp.Error = err.Error()
			slog.ErrorContext(ctx, fmt.Sprintf("[service][UserVerification][SetUserVerified] err : %v", err))

			return resp, err
		}
	}

	token, exp, err := utils.Sign(JWTData{
		Email:   req.Email,
		UserID:  int64(user.UserID),
		Usename: user.Username,
	})
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("[service][UserRegistration][Sign] err : %v", err))
		resp.Code = http.StatusInternalServerError
		resp.Message = utils.ErrorInternalServer
		resp.Error = err.Error()

		return
	}

	resp.Data = struct {
		Token  string    `json:"token"`
		Expire time.Time `json:"expire"`
	}{
		Token:  token,
		Expire: exp,
	}

	return
}
