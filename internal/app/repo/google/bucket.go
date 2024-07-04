package google

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"go.uber.org/dig"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

type (
	BucketRepo interface {
		PresignedURL(ctx context.Context, privateUrl string) (url string, err error)
	}

	BucketRepoImpl struct {
		dig.In
	}
)

func NewBucketRepo(impl BucketRepoImpl) BucketRepo {
	return &impl
}

func (b *BucketRepoImpl) PresignedURL(ctx context.Context, privateUrl string) (url string, err error) {
	bucketName := os.Getenv("GOOGLE_BUCKET_NAME")
	objectName := privateUrl[len("https://storage.cloud.google.com/mahir_trade_video/"):]
	temporaryAccess := time.Now().Add(6 * time.Hour)

	client, err := storage.NewClient(ctx, option.WithCredentialsFile(os.Getenv("GOOGLE_SERVICE_ACCOUNT_FILE_PATH")))
	if err != nil {
		slog.ErrorContext(ctx, "[repo][google][PresignedURL][NewClient]", err)
		return
	}
	defer func() {
		if err := client.Close(); err != nil {
			slog.ErrorContext(ctx, "[repo][google][PresignedURL][Close]", err)
		}
	}()

	jsonKey, err := os.ReadFile(os.Getenv("GOOGLE_SERVICE_ACCOUNT_FILE_PATH"))
	if err != nil {
		slog.ErrorContext(ctx, "[repo][google][PresignedURL][ReadFile]", err)
		err = fmt.Errorf("something went wrong, we will fix it soon")
		return
	}

	conf, err := google.JWTConfigFromJSON(jsonKey)
	if err != nil {
		return "", fmt.Errorf("google.JWTConfigFromJSON: %v", err)
	}

	opts := storage.SignedURLOptions{
		GoogleAccessID: conf.Email,
		PrivateKey:     conf.PrivateKey,
		Method:         http.MethodGet,
		Expires:        temporaryAccess,
	}

	url, err = storage.SignedURL(bucketName, objectName, &opts)
	if err != nil {
		slog.ErrorContext(ctx, "[repo][google][PresignedURL][SignedURL]", err)
		return
	}
	return

}
