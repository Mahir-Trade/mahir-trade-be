package utils

const (
	StatusOrderPending = "pending"
	StatusOrderSuccess = "success"
	StatusOrderFailed  = "failed"

	MidtransStatusSettlement = "settlement"
	MidtransStatusPending    = "pending"
	MidtransStatusCapture    = "capture"

	ErrorInternalServer     = "internal server error, we will fix it soon"
	ErrorBadRequest         = "bad request"
	ErrorUnauthorized       = "unauthorized"
	ErrorInvalidRequestBody = "invalid request body"

	MembershipProgramStartDateConfig = "MEMBERSHIP_PROGRAM_START_DATE"
	MembershipProgramEndDateConfig   = "MEMBERSHIP_PROGRAM_END_DATE"
)
