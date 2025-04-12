package models

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

type (
	Report struct {
		ID                 int64           `json:"id"`
		ReportName         string          `json:"report_name" validate:"required"`
		ReportThumbnailURL string          `json:"report_thumbnail_url" validate:"required"`
		Contents           json.RawMessage `json:"contents" validate:"required"`
		CreatedBy          string          `json:"created_by"`
		UpdatedBy          string          `json:"updated_by"`
		CreatedAt          string          `json:"created_at"`
		UpdatedAt          string          `json:"updated_at"`
	}
)

func (r *Report) ValidateContents() error {
	contentsStr := strings.TrimSpace(string(r.Contents))
	if len(contentsStr) == 0 || contentsStr == "[]" || contentsStr == "[{}]" {
		return fmt.Errorf("contents field is invalid")
	}

	whitespaceArrayRegex := regexp.MustCompile(`^\[\s*\]$`)
	if whitespaceArrayRegex.MatchString(contentsStr) {
		return fmt.Errorf("contents field contains only whitespace or empty array")
	}

	var contents []map[string]interface{}
	if err := json.Unmarshal(r.Contents, &contents); err != nil {
		return fmt.Errorf("contents field is not valid JSON: %v", err)
	}

	for _, item := range contents {
		if len(item) == 0 {
			return fmt.Errorf("contents field contains empty objects")
		}
	}

	return nil
}
