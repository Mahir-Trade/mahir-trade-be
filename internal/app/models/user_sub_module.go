package models

type (
	UserSubModule struct {
		ID          int64  `json:"id,omitempty"`
		UUID        string `json:"uuid,omitempty"`
		UserID      int64  `json:"user_id"`
		SubModuleID int64  `json:"sub_module_id"`
		CreatedBy   string `json:"created_by,omitempty"`
		UpdatedBy   string `json:"updated_by,omitempty"`
		CreatedAt   string `json:"created_at,omitempty"`
		UpdatedAt   string `json:"updated_at,omitempty"`
	}
)
