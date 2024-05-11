package models

type (
	Package struct {
		ID              int64   `json:"id,omitempty"`
		Price           float64 `json:"price" validate:"required"`
		DurationInMonth int64   `json:"duration_in_month" validate:"required"`
		Description     string  `json:"description" validate:"required"`
		CreatedBy       string  `json:"created_by,omitempty"`
		UpdatedBy       string  `json:"updated_by,omitempty"`
		CreatedAt       string  `json:"created_at,omitempty"`
		UpdatedAt       string  `json:"updated_at,omitempty"`
	}

	GetPackagesRequest struct {
		Limit int64 `json:"limit"`
		Page  int64 `json:"page"`
	}
)
