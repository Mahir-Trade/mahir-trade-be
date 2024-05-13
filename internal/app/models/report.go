package models

type (
	Report struct {
		ID                 int64  `json:"id"`
		ReportThumbnailURL string `json:"report_thumbnail_url" validate:"required"`
		ReportFileURL      string `json:"report_file_url" validate:"required"`
		CreatedBy          string `json:"created_by"`
		UpdatedBy          string `json:"updated_by"`
		CreatedAt          string `json:"created_at"`
		UpdatedAt          string `json:"updated_at"`
	}
)
