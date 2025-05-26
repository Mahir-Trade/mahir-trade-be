package models

type MembershipStatus string

const (
	MembershipStatusUnknown  MembershipStatus = "UNKNOWN"
	MembershipStatusPreOrder MembershipStatus = "PRE_ORDER"
	MembershipStatusActive   MembershipStatus = "ACTIVE"
	MembershipStatusExpired  MembershipStatus = "EXPIRED"
)

func (s MembershipStatus) IsValid() bool {
	switch s {
	case MembershipStatusUnknown,
		MembershipStatusPreOrder,
		MembershipStatusActive,
		MembershipStatusExpired:
		return true
	default:
		return false
	}
}

type (
	UserMembership struct {
		ID                 int64            `json:"id,omitempty"`
		UserID             int64            `json:"user_id" validate:"required"`
		PackageID          int64            `json:"package_id" validate:"required"`
		ExpiredAt          string           `json:"expired_at" validate:"required"`
		ExclusiveExpiredAt string           `json:"exclusive_expired_at,omitempty"`
		IsMembershipActive bool             `json:"is_membership_active" validate:"required"`
		Status             MembershipStatus `json:"status"`
		CreatedBy          string           `json:"created_by,omitempty"`
		UpdatedBy          string           `json:"updated_by,omitempty"`
		CreatedAt          string           `json:"created_at,omitempty"`
		UpdatedAt          string           `json:"updated_at,omitempty"`
	}
)
