package google

import (
	"context"
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

	GoogleRepo interface {
		Login(ctx context.Context) (url string, err error)
		Callback(ctx context.Context, req GoogleCallbackRequest) (err error)
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

func (g *GoogleRepoImpl) Callback(ctx context.Context, req GoogleCallbackRequest) (err error) {
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

	slog.InfoContext(ctx, "[repo][google][Callback] response:", string(userData), err)

	return
}
