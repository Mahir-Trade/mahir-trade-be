package google

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"go.uber.org/dig"
	"golang.org/x/oauth2"
)

type (
	GoogleCallbackRequest struct {
		State string `json:"state"`
		Code  string `json:"code"`
	}

	GoogleCallbackResponse struct {
		ID            string `json:"id"`
		Email         string `json:"email"`
		Name          string `json:"name"`
		Picture       string `json:"picture"`
		VerifiedEmail bool   `json:"verified_email"`
	}

	GoogleRepo interface {
		Login(ctx context.Context) (url string, err error)
		Callback(ctx context.Context, req GoogleCallbackRequest) (dataInfo GoogleCallbackResponse, err error)
	}

	GoogleRepoImpl struct {
		dig.In

		*oauth2.Config
	}
)

func NewGoogleRepo(impl GoogleRepoImpl) GoogleRepo {
	return &impl
}

func (g *GoogleRepoImpl) Login(ctx context.Context) (url string, err error) {
	url = g.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	return
}

func (g *GoogleRepoImpl) Callback(ctx context.Context, req GoogleCallbackRequest) (dataInfo GoogleCallbackResponse, err error) {
	if req.State != "state-token" {
		err = fmt.Errorf("invalid credentials")
		slog.ErrorContext(ctx, "[repo][google][Callback][state]", err)
		return
	}

	token, err := g.Exchange(ctx, req.Code)
	if err != nil {
		slog.ErrorContext(ctx, "[repo][google][Callback][Exchange]", err)
		return
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		slog.ErrorContext(ctx, "[repo][google][Callback][http.Get]", err)
		return
	}
	defer resp.Body.Close()

	userData, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.ErrorContext(ctx, "[repo][google][Callback][io.ReadAll]", err)
		return
	}
	err = json.Unmarshal(userData, &dataInfo)
	if err != nil {
		slog.ErrorContext(ctx, "[repo][google][Callback][json.Unmarshal]", err)
		return
	}

	return
}
