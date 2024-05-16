package midtrans

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"

	"go.uber.org/dig"
)

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

	GeneratePaymentLinkRequest struct {
		TransactionDetails TransactionDetails `json:"transaction_details"`
		UsageLimit         int                `json:"usage_limit"`
		Expiry             Expiry             `json:"expiry"`
		EnabledPayments    []string           `json:"enabled_payments"`
		ItemDetails        []ItemDetails      `json:"item_details"`
		CustomerDetails    CustomerDetails    `json:"customer_details"`
	}

	GeneratePaymentLinkResponse struct {
		OrderID       string   `json:"order_id,omitempty"`
		PaymentURL    string   `json:"payment_url,omitempty"`
		ErrorMessages []string `json:"error_messages,omitempty"`
	}

	MidtransRepo interface {
		GeneratePaymentLink(ctx context.Context, req GeneratePaymentLinkRequest) (res GeneratePaymentLinkResponse, err error)
	}

	MidtransRepoImpl struct {
		dig.In
	}
)

func NewMidtransRepo(impl MidtransRepoImpl) MidtransRepo {
	return &impl
}

func generateTokenAuthorization() string {
	clientSecret := os.Getenv("MIDTRANS_CLIENT_SECRET") + ":"
	encoded := base64.StdEncoding.EncodeToString([]byte(clientSecret))

	return "Basic " + encoded
}

func (m *MidtransRepoImpl) GeneratePaymentLink(ctx context.Context, req GeneratePaymentLinkRequest) (res GeneratePaymentLinkResponse, err error) {
	generatePaymentLinkURL := os.Getenv("MIDTRANS_BASE_URL") + "/v1/payment-links"

	requestPayload := GeneratePaymentLinkRequest{
		TransactionDetails: req.TransactionDetails,
		UsageLimit:         1,
		Expiry: Expiry{
			StartTime: time.Now().Format("2006-01-02 15:04 -0700"),
			Duration:  1,
			Unit:      "days",
		},
		EnabledPayments: []string{
			"credit_card",
			"bca_klikbca",
			"gopay",
			"permata_va",
			"bca_va",
			"bri_va",
			"bni_va",
			"indomaret",
			"shopeepay",
		},
		ItemDetails:     req.ItemDetails,
		CustomerDetails: req.CustomerDetails,
	}

	payload, err := json.Marshal(requestPayload)
	if err != nil {
		slog.ErrorContext(ctx, "[repo][midtrans][Generate PaymentLink] Error marshal request payload: ", err)
		return res, err
	}

	request, err := http.NewRequest("POST", generatePaymentLinkURL, bytes.NewBuffer(payload))
	if err != nil {
		slog.ErrorContext(ctx, "[repo][midtrans][Generate PaymentLink] Error create new request: ", err)
		return res, err
	}

	request.Header.Set("Authorization", generateTokenAuthorization())
	request.Header.Set("Content-Type", "application/json")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		slog.ErrorContext(ctx, "[repo][midtrans][Generate PaymentLink] Error sending request: ", err)
		return res, err
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		slog.ErrorContext(ctx, "[repo][midtrans][Generate PaymentLink] Error reading response body: ", err)
		return res, err
	}

	err = json.Unmarshal(body, &res)
	if err != nil {
		slog.ErrorContext(ctx, "[repo][midtrans][Generate PaymentLink] Error unmarshal response body: ", err)
		return res, err
	}

	return res, nil
}
