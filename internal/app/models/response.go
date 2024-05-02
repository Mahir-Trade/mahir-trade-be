package models

type (
	DefaultResponse struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    any    `json:"data,omitempty"`
		Error   any    `json:"error,omitempty"`
	}

	DefaultMetaData struct {
		Page        uint `json:"page"`
		TotalPages  uint `json:"totalPages"`
		TotalItems  uint `json:"totalItems"`
		Limit       uint `json:"limit"`
		HasNext     bool `json:"hasNext"`
		HasPrevious bool `json:"hasPrevious"`
	}

	DefaultPaginationResponseData struct {
		Results         interface{} `json:"results"`
		DefaultMetaData `json:"meta"`
	}
)
