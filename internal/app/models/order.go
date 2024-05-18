package models

type (
	Order struct {
		ID          int64  `json:"id,omitempty"`
		UserID      int64  `json:"user_id" validate:"required"`
		PackageID   int64  `json:"package_id" validate:"required"`
		Status      string `json:"status" validate:"required"`
		PaymentCode string `json:"payment_code" validate:"required"`
		PaymentURL  string `json:"payment_url" validate:"required"`
		CreatedBy   string `json:"created_by,omitempty"`
		UpdatedBy   string `json:"updated_by,omitempty"`
		CreatedAt   string `json:"created_at,omitempty"`
		UpdatedAt   string `json:"updated_at,omitempty"`
	}
)
