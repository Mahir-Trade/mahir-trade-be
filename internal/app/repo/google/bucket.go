package google

import (
	"context"
	"fmt"
	"log/slog"
	"mahir-trade-be/internal/app/infra"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	transcoder "cloud.google.com/go/video/transcoder/apiv1"
	"cloud.google.com/go/video/transcoder/apiv1/transcoderpb"
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
		BucketName      string `json:"bucket_name"`
	}

	URLParserResponse struct {
		URL        string `json:"url"`
		Host       string `json:"host"`
		Path       string `json:"path"`
		BucketName string `json:"bucket_name"`
	}

	BucketRepo interface {
		PresignedURL(ctx context.Context, bucketName, privateUrl string) (url string, err error)
		UploadFile(ctx context.Context, form FileUpload, file multipart.File) (url string, err error)
		URLParser(fileURL string) (resp URLParserResponse)
		StartTranscodingJob(bucketName, filename string)
	}

	BucketRepoImpl struct {
		dig.In

		GoogleCfg *infra.GoogleCfg
	}
)

func NewBucketRepo(impl BucketRepoImpl) BucketRepo {
	return &impl
}

func (b *BucketRepoImpl) PresignedURL(ctx context.Context, bucketName, privateUrl string) (signedUrl string, err error) {

	urlParsed, err := url.Parse(privateUrl)
	if err != nil {
		slog.ErrorContext(ctx, "[service][GetSubModuleByID] error while parse url err: %v", err)
		return
	}

	decodePath := urlParsed.Path
	if decodePath, err = url.QueryUnescape(decodePath); err != nil {
		slog.ErrorContext(ctx, "[service][GetSubModuleByID] error while decode url err: %v", err)
		return
	}
	cleanUrl := urlParsed.Scheme + "://" + urlParsed.Host + decodePath

	objectName := cleanUrl[len(fmt.Sprintf("https://%s/%s/", urlParsed.Host, bucketName)):]
	temporaryAccess := time.Now().Add(1 * time.Hour)

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

	signedUrl, err = storage.SignedURL(bucketName, objectName, &opts)
	if err != nil {
		slog.ErrorContext(ctx, "[repo][google][PresignedURL][SignedURL]", err)
		return
	}
	return
}

func (b *BucketRepoImpl) UploadFile(ctx context.Context, form FileUpload, file multipart.File) (url string, err error) {
	cmd := exec.Command("gsutil", "cp", form.LocalFilePath, fmt.Sprintf("gs://%s/%s", form.BucketName, form.Filename))
	output, err := cmd.CombinedOutput()
	if err != nil {
		slog.ErrorContext(ctx, "[repo][google][UploadFile][CombinedOutput] err :", err)
		return url, fmt.Errorf("failed to upload file to GCS: %v, output: %s", err, string(output))
	}

	cmdContentType := exec.Command("gsutil", "setmeta", "-h", fmt.Sprintf("Content-Type:%s", form.FileContentType), fmt.Sprintf("gs://%s/%s", form.BucketName, form.Filename))
	outputContentType, err := cmdContentType.CombinedOutput()
	if err != nil {
		slog.ErrorContext(ctx, "[repo][google][UploadFile][Command] err :", err)
		return url, fmt.Errorf("failed to set content type: %v, output: %s", err, string(outputContentType))
	}

	url, err = b.PresignedURL(ctx, form.BucketName, fmt.Sprintf("https://storage.cloud.google.com/%s/%s", form.BucketName, form.Filename))
	if err != nil {
		slog.ErrorContext(ctx, "[repo][google][UploadFile][PresignedURL]", err)
		err = fmt.Errorf("something went wrong, we will fix it soon")
		return
	}

	return
}

func (b *BucketRepoImpl) URLParser(fileURL string) (resp URLParserResponse) {
	urlParsed, err := url.Parse(fileURL)
	if err != nil {
		slog.Error("[repo][google][URLParser] error while parse url err: %v", err)
		return
	}

	decodePath := urlParsed.Path
	if decodePath, err = url.QueryUnescape(decodePath); err != nil {
		slog.Error("[repo][google][URLParser] error while decode url err: %v", err)
		return
	}
	cleanUrl := urlParsed.Scheme + "://" + urlParsed.Host + decodePath

	resp.URL = cleanUrl
	resp.Host = urlParsed.Host
	resp.Path = decodePath
	resp.BucketName = strings.Split(decodePath, "/")[1]

	return resp
}

func (b *BucketRepoImpl) StartTranscodingJob(bucketName, filename string) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	// Create a new Transcoder client
	client, err := transcoder.NewClient(ctx, option.WithCredentialsFile(os.Getenv("GOOGLE_SERVICE_ACCOUNT_FILE_PATH")))
	if err != nil {
		slog.ErrorContext(ctx, "[repo][google][StartTranscodingJob][CreateClient] err :", err)
		return
	}
	defer client.Close()

	// Define input and output paths
	inputURI := fmt.Sprintf("gs://%s/%s", bucketName, filename)
	outputURI := fmt.Sprintf("gs://%s/", fmt.Sprintf("%s_transcoded", bucketName))

	fileName240p := strings.Replace(filename, ".mp4", "-240p", -1)
	fileName360p := strings.Replace(filename, ".mp4", "-360p", -1)
	fileName480p := strings.Replace(filename, ".mp4", "-480p", -1)
	fileName720p := strings.Replace(filename, ".mp4", "-720p", -1)

	parent := "projects/mahir-trade-429013/locations/asia-southeast1"
	audioStream := &transcoderpb.ElementaryStream{
		Key: "audio",
		ElementaryStream: &transcoderpb.ElementaryStream_AudioStream{
			AudioStream: &transcoderpb.AudioStream{
				Codec:      "aac", // Or other supported codec
				BitrateBps: 128000,
				// Other audio stream fields as needed
			},
		},
	}

	// Define job configuration for different resolutions
	job := &transcoderpb.CreateJobRequest{
		Parent: parent,
		Job: &transcoderpb.Job{
			InputUri:  inputURI,
			OutputUri: outputURI,
			JobConfig: &transcoderpb.Job_Config{
				Config: &transcoderpb.JobConfig{
					ElementaryStreams: []*transcoderpb.ElementaryStream{
						audioStream,
						{
							Key: fileName240p,
							ElementaryStream: &transcoderpb.ElementaryStream_VideoStream{
								VideoStream: &transcoderpb.VideoStream{
									CodecSettings: &transcoderpb.VideoStream_H264{
										H264: &transcoderpb.VideoStream_H264CodecSettings{
											BitrateBps:   500000,
											FrameRate:    30,
											HeightPixels: 240,
											WidthPixels:  426,
										},
									},
								},
							},
						},
						{
							Key: fileName360p,
							ElementaryStream: &transcoderpb.ElementaryStream_VideoStream{
								VideoStream: &transcoderpb.VideoStream{
									CodecSettings: &transcoderpb.VideoStream_H264{
										H264: &transcoderpb.VideoStream_H264CodecSettings{
											BitrateBps:   800000,
											FrameRate:    30,
											HeightPixels: 360,
											WidthPixels:  640,
										},
									},
								},
							},
						},
						{
							Key: fileName480p,
							ElementaryStream: &transcoderpb.ElementaryStream_VideoStream{
								VideoStream: &transcoderpb.VideoStream{
									CodecSettings: &transcoderpb.VideoStream_H264{
										H264: &transcoderpb.VideoStream_H264CodecSettings{
											BitrateBps:   1000000,
											FrameRate:    30,
											HeightPixels: 480,
											WidthPixels:  854,
										},
									},
								},
							},
						},
						{
							Key: fileName720p,
							ElementaryStream: &transcoderpb.ElementaryStream_VideoStream{
								VideoStream: &transcoderpb.VideoStream{
									CodecSettings: &transcoderpb.VideoStream_H264{
										H264: &transcoderpb.VideoStream_H264CodecSettings{
											BitrateBps:   2500000,
											FrameRate:    30,
											HeightPixels: 720,
											WidthPixels:  1280,
										},
									},
								},
							},
						},
					},
					MuxStreams: []*transcoderpb.MuxStream{
						{
							Key:               fileName240p,
							ElementaryStreams: []string{fileName240p, "audio"},
							Container:         "mp4",
						},
						{
							Key:               fileName360p,
							ElementaryStreams: []string{fileName360p, "audio"},
							Container:         "mp4",
						},
						{
							Key:               fileName480p,
							ElementaryStreams: []string{fileName480p, "audio"},
							Container:         "mp4",
						},
						{
							Key:               fileName720p,
							ElementaryStreams: []string{fileName720p, "audio"},
							Container:         "mp4",
						},
					},
				},
			},
		},
	}

	// Create the transcoding job
	slog.InfoContext(ctx, "[repo][google][StartTranscodingJob][CreateJob] job start for filename: %s", filename)
	resp, err := client.CreateJob(ctx, job)
	if err != nil {
		slog.ErrorContext(ctx, "[repo][google][StartTranscodingJob][CreateJob] err :", err)
		return
	}

	// Check job status periodically
	jobName := resp.GetName()
	for {
		getJobResp, err := client.GetJob(ctx, &transcoderpb.GetJobRequest{Name: jobName})
		if err != nil {
			slog.ErrorContext(ctx, "[repo][google][StartTranscodingJob][GetJob] err :", err)
			return
		}

		if getJobResp.State != transcoderpb.Job_SUCCEEDED {
			if getJobResp.State == transcoderpb.Job_FAILED {
				slog.Error("Job failed: %v", getJobResp.GetError().Message)
				return
			}
			time.Sleep(5 * time.Second)
			continue
		}

		// Job succeeded
		slog.Info("Job succeeded: %s", jobName)
		break
	}

	return
}
