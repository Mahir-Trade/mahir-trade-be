package infra

import (
	"go.uber.org/dig"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type (
	Oauth struct {
		dig.Out
		GoogleOauth *oauth2.Config
	}

	OauthCfgs struct {
		dig.In
		Google *GoogleCfg
	}

	GoogleCfg struct {
		ClientID               string `envconfig:"GOOGLE_CLIENT_ID" required:"true"`
		ClientSecret           string `envconfig:"GOOGLE_CLIENT_SECRET" required:"true"`
		RedirectURL            string `envconfig:"GOOGLE_REDIRECT_URL" required:"true"`
		ServiceAccountFilePath string `envconfig:"GOOGLE_SERVICE_ACCOUNT_FILE_PATH" required:"true"`
		VideoBucketName        string `envconfig:"GOOGLE_VIDEO_BUCKET_NAME" required:"true"`
		ImageBucketName        string `envconfig:"GOOGLE_IMAGE_BUCKET_NAME" required:"true"`
		FileBucketName         string `envconfig:"GOOGLE_FILE_BUCKET_NAME" required:"true"`
	}
)

func NewOauth(cfgs OauthCfgs) Oauth {
	return Oauth{
		GoogleOauth: initGoogleAuth(cfgs.Google),
	}
}

func initGoogleAuth(cfg *GoogleCfg) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		Endpoint:     google.Endpoint,
		RedirectURL:  cfg.RedirectURL,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
	}
}
