package models

type (
	User struct {
		UserID      int    `json:"user_id,omitempty"`
		UUID        string `json:"uuid,omitempty"`
		Email       string `json:"email" validate:"required,email"`
		Fullname    string `json:"fullname" validate:"required"`
		PhoneNumber string `json:"phone_number,omitempty"`
		Username    string `json:"username" validate:"required"`
		Password    string `json:"password" validate:"required,min=8,max=20,uppercase,lowercase,number,specialchar"`
	}
)
