package models

type (
	Admin struct {
		AdminID   int64  `json:"admin_id,omitempty"`
		UUID      string `json:"uuid,omitempty"`
		Email     string `json:"email" validate:"required,email"`
		Username  string `json:"username" validate:"required"`
		Password  string `json:"password" validate:"required"`
		CreatedAt string `json:"created_at,omitempty"`
		UpdatedAt string `json:"updated_at,omitempty"`
	}
)
