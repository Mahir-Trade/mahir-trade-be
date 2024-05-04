package models

type (
	Group struct {
		ID        int64  `json:"id,omitempty"`
		UUID      string `json:"uuid,omitempty"`
		GroupName string `json:"group_name" validate:"required"`
		CreatedAt string `json:"created_at,omitempty"`
		UpdatedAt string `json:"updated_at,omitempty"`
	}

	GetGroupsRequest struct {
		Limit int64 `json:"limit"`
		Page  int64 `json:"page"`
	}
)
