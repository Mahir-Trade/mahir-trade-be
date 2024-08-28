package models

type (
	SendgridSendEmailRequest struct {
		From          string     `json:"from" validate:"required,email"`
		To            string     `json:"to" validate:"required,email"`
		SenderName    string     `json:"sender_name" validate:"required"`
		RecepientName string     `json:"recepient_name" validate:"required"`
		Subject       string     `json:"subject" validate:"required"`
		Body          string     `json:"body" validate:"required"`
		Attach        Attachment `json:"attach"`
	}

	Attachment struct {
		Content     string `json:"content,omitempty"`
		Type        string `json:"type,omitempty"`
		Filename    string `json:"filename,omitempty"`
		Disposition string `json:"disposition,omitempty"`
		FileName    string `json:"file_name,omitempty"`
		ContentID   string `json:"content_id,omitempty"`
	}
)
