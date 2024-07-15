package google

import (
	"context"
	"fmt"
	"log/slog"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"time"

	"cloud.google.com/go/storage"
	"go.uber.org/dig"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

type (
	FileUpload struct {
		Filename        string `json:"filename"`
		Size            int64  `json:"size"`
		LocalFilePath   string `json:"local_file_path"`
		FileContentType string `json:"file_content_type"`
	}

	BucketRepo interface {
		PresignedURL(ctx context.Context, privateUrl string) (url string, err error)
		UploadFile(ctx context.Context, form FileUpload, file multipart.File) (url string, err error)
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

func (b *BucketRepoImpl) UploadFile(ctx context.Context, form FileUpload, file multipart.File) (url string, err error) {
	bucketName := os.Getenv("GOOGLE_BUCKET_NAME")

	cmd := exec.Command("gsutil", "cp", form.LocalFilePath, fmt.Sprintf("gs://%s/%s", bucketName, form.Filename))
	output, err := cmd.CombinedOutput()
	if err != nil {
		slog.ErrorContext(ctx, "[repo][google][UploadFile][CombinedOutput] err :", err)
		return url, fmt.Errorf("failed to upload file to GCS: %v, output: %s", err, string(output))
	}

	cmdContentType := exec.Command("gsutil", "setmeta", "-h", fmt.Sprintf("Content-Type:%s", form.FileContentType), fmt.Sprintf("gs://%s/%s", bucketName, form.Filename))
	outputContentType, err := cmdContentType.CombinedOutput()
	if err != nil {
		slog.ErrorContext(ctx, "[repo][google][UploadFile][Command] err :", err)
		return url, fmt.Errorf("failed to set content type: %v, output: %s", err, string(outputContentType))
	}

	url, err = b.PresignedURL(ctx, fmt.Sprintf("https://storage.cloud.google.com/%s/%s", bucketName, form.Filename))
	if err != nil {
		slog.ErrorContext(ctx, "[repo][google][UploadFile][PresignedURL]", err)
		err = fmt.Errorf("something went wrong, we will fix it soon")
		return
	}

	return
}
