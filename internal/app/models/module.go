package models

import "database/sql"

type (
	Module struct {
		ID           int64          `json:"id,omitempty"`
		UUID         string         `json:"uuid,omitempty"`
		GroupID      sql.NullInt64  `json:"group_id,omitempty"`
		ModuleName   string         `json:"module_name" validate:"required"`
		ThumbnailUrl sql.NullString `json:"thumbnail_url,omitempty"`
		Tag          sql.NullString `json:"tag,omitempty"`
		CreatedBy    string         `json:"created_by,omitempty"`
		UpdatedBy    string         `json:"updated_by,omitempty"`
		CreatedAt    string         `json:"created_at,omitempty"`
		UpdatedAt    string         `json:"updated_at,omitempty"`
	}

	GetModulesRequest struct {
		Limit int64 `json:"limit"`
		Page  int64 `json:"page"`
	}
)
