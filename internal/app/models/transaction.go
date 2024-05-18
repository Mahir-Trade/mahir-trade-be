package models

type (
	Transaction struct {
		ID             int64   `json:"id,omitempty"`
		OrderID        int64   `json:"order_id" validate:"required"`
		Amount         float64 `json:"amount" validate:"required"`
		SettlementDate string  `json:"settlement_date,omitempty"`
		WebhookID      string  `json:"webhook_id,omitempty"`
		CreatedBy      string  `json:"created_by,omitempty"`
		UpdatedBy      string  `json:"updated_by,omitempty"`
		CreatedAt      string  `json:"created_at,omitempty"`
		UpdatedAt      string  `json:"updated_at,omitempty"`
	}
)
