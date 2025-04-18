package models

type (
	PaginationRequest struct {
		Limit            int64    `query:"limit" validate:"gte=1,lte=100"`
		Page             int64    `query:"page" validate:"gte=1"`
		Search           string   `query:"search" validate:"omitempty"`
		Filters          []string `query:"filters,omitempty"` // optional
		ShowAll          bool     `query:"showAll"`
		MembershipStatus string   `query:"membershipStatus,omitempty"`
	}

	StartMembershipProgramRequest struct {
		EndDate string `json:"end_date" validate:"required"`
	}
)
