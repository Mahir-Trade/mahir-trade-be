package models

type (
	User struct {
		UserID      int64  `json:"user_id,omitempty"`
		UUID        string `json:"uuid,omitempty"`
		Email       string `json:"email" validate:"required,email"`
		PhoneNumber string `json:"phone_number,omitempty"`
		Username    string `json:"username" validate:"required"`
		IsActive    bool   `json:"is_active,omitempty"`
		Password    string `json:"password,omitempty" validate:"required,min=8,max=20,uppercase,lowercase,number,specialchar"`
		CreatedAt   string `json:"created_at,omitempty"`
		CreatedBy   string `json:"created_by,omitempty"`
		UpdatedAt   string `json:"updated_at,omitempty"`
		UpdatedBy   string `json:"updated_by,omitempty"`
	}

	UserRegistrationRequest struct {
		Email                string `json:"email" validate:"required,email"`
		PhoneNumber          string `json:"phone_number,omitempty"`
		Username             string `json:"username" validate:"required"`
		Password             string `json:"password,omitempty" validate:"required,min=8,max=20,uppercase,lowercase,number,specialchar"`
		PasswordConfirmation string `json:"password_confirmation,omitempty"`
	}
)
