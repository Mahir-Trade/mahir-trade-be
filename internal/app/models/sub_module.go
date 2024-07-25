package models

import "database/sql"

type (
	SubModule struct {
		ID            int64         `json:"id,omitempty"`
		UUID          string        `json:"uuid,omitempty"`
		ModuleID      sql.NullInt64 `json:"module_id"`
		SubModuleName string        `json:"sub_module_name" validate:"required"`
		Title         string        `json:"title" validate:"required"`
		VideoURL      string        `json:"video_url" validate:"required"`
		Status        string        `json:"status,omitempty"`
		CreatedBy     string        `json:"created_by,omitempty"`
		UpdatedBy     string        `json:"updated_by,omitempty"`
		CreatedAt     string        `json:"created_at,omitempty"`
		UpdatedAt     string        `json:"updated_at,omitempty"`
	}
)
