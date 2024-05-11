package models

type (
	DiscordAccount struct {
		ID               int64  `json:"id,omitempty"`
		UUID             string `json:"uuid,omitempty"`
		UserID           int64  `json:"user_id,omitempty"`
		DiscordAccountID string `json:"discord_account_id,omitempty"`
		Username         string `json:"username,omitempty"`
		Email            string `json:"email,omitempty"`
		CreatedAt        string `json:"created_at,omitempty"`
		UpdatedAt        string `json:"updated_at,omitempty"`
	}
)
