package models

type (
	TransactionDetails struct {
		OrderID     string  `json:"order_id"`
		GrossAmount float64 `json:"gross_amount"`
		// PaymentLinkID string  `json:"payment_link_id"`
	}

	Expiry struct {
		StartTime string `json:"start_time"`
		Duration  int    `json:"duration"`
		Unit      string `json:"unit"`
	}

	ItemDetails struct {
		// ID           string  `json:"id"`
		Name     string  `json:"name"`
		Price    float64 `json:"price"`
		Quantity int     `json:"quantity"`
		// Brand        string  `json:"brand"`
		// Category     string  `json:"category"`
		// MerchantName string  `json:"merchant_name"`
	}

	CustomerDetails struct {
		FirstName string `json:"first_name"`
		// LastName  string `json:"last_name"`
		Email string `json:"email"`
		Phone string `json:"phone"`
		// Notes     string `json:"notes"`
	}

	MidtransGeneratePaymentLinkRequest struct {
		TransactionDetails TransactionDetails `json:"transaction_details"`
		UsageLimit         int                `json:"usage_limit"`
		Expiry             Expiry             `json:"expiry"`
		EnabledPayments    []string           `json:"enabled_payments"`
		ItemDetails        []ItemDetails      `json:"item_details"`
		CustomerDetails    CustomerDetails    `json:"customer_details"`
	}

	MidtransGeneratePaymentLinkResponse struct {
		OrderID       string   `json:"order_id,omitempty"`
		PaymentURL    string   `json:"payment_url,omitempty"`
		ErrorMessages []string `json:"error_messages,omitempty"`
	}

	// callback
	MidtransCalbackMetadata struct {
		ExtraInfo struct {
			PaymentLinkID   string `json:"payment_link_id"`
			PaymentLinkType string `json:"payment_link_type"`
		} `json:"extra_info"`
	}

	MidtransCallbackVaNumbers struct {
		Bank     string `json:"bank"`
		VaNumber string `json:"va_number"`
	}

	MidtransCallbackPaymentAmount struct {
		PaidAt string `json:"paid_at"`
		Amount string `json:"amount"`
	}

	MidtransCallbackRequest struct {
		Currency          string                          `json:"currency,omitempty"`
		CustomField1      string                          `json:"custom_field1,omitempty"`
		CustomField2      string                          `json:"custom_field2,omitempty"`
		CustomField3      string                          `json:"custom_field3,omitempty"`
		ExpiryTime        string                          `json:"expiry_time,omitempty"`
		FraudStatus       string                          `json:"fraud_status,omitempty"`
		GrossAmount       string                          `json:"gross_amount,omitempty"`
		MerchantID        string                          `json:"merchant_id,omitempty"`
		Metadata          MidtransCalbackMetadata         `json:"metadata,omitempty"`
		OrderID           string                          `json:"order_id,omitempty"`
		PaymentAmounts    []MidtransCallbackPaymentAmount `json:"payment_amounts,omitempty"`
		PaymentType       string                          `json:"payment_type,omitempty"`
		SettlementTime    string                          `json:"settlement_time,omitempty"`
		SignatureKey      string                          `json:"signature_key,omitempty"`
		StatusCode        string                          `json:"status_code,omitempty"`
		StatusMessage     string                          `json:"status_message,omitempty"`
		TransactionID     string                          `json:"transaction_id,omitempty"`
		TransactionStatus string                          `json:"transaction_status,omitempty"`
		TransactionTime   string                          `json:"transaction_time,omitempty"`
		VaNumbers         []MidtransCallbackVaNumbers     `json:"va_numbers,omitempty"`

		// additional for credit card
		MaskedCard             string `json:"masked_card,omitempty"`
		Eci                    string `json:"eci,omitempty"`
		ChannelResponseMessage string `json:"channel_response_message,omitempty"`
		ChannelResponseCode    string `json:"channel_response_code,omitempty"`
		CardType               string `json:"card_type,omitempty"`
		Bank                   string `json:"bank,omitempty"`
		ApprovalCode           string `json:"approval_code,omitempty"`

		// additional for QRIS
		Issuer   string `json:"issuer,omitempty"`
		Acquirer string `json:"acquirer,omitempty"`

		// additional for Permata VA
		PermataVaNumber string `json:"permata_va_number,omitempty"`
	}

	MidtransCheckStatusResponse struct {
		MaskedCard             string `json:"masked_card"`
		ApprovalCode           string `json:"approval_code"`
		Bank                   string `json:"bank"`
		Eci                    string `json:"eci"`
		ChannelResponseCode    string `json:"channel_response_code"`
		ChannelResponseMessage string `json:"channel_response_message"`
		TransactionTime        string `json:"transaction_time"`
		GrossAmount            string `json:"gross_amount"`
		Currency               string `json:"currency"`
		OrderID                string `json:"order_id"`
		PaymentType            string `json:"payment_type"`
		SignatureKey           string `json:"signature_key"`
		StatusCode             string `json:"status_code"`
		TransactionID          string `json:"transaction_id"`
		TransactionStatus      string `json:"transaction_status"`
		FraudStatus            string `json:"fraud_status"`
		SettlementTime         string `json:"settlement_time"`
		StatusMessage          string `json:"status_message"`
		MerchantID             string `json:"merchant_id"`
		CardType               string `json:"card_type"`
		ThreeDsVersion         string `json:"three_ds_version"`
		ChallengeCompletion    bool   `json:"challenge_completion"`
	}
)
