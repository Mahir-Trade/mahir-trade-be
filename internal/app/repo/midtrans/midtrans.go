package midtrans

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"mahir-trade-be/internal/app/models"
	"net/http"
	"os"
	"time"

	"go.uber.org/dig"
)

type (
	MidtransRepo interface {
		GeneratePaymentLink(ctx context.Context, req models.MidtransGeneratePaymentLinkRequest) (res models.MidtransGeneratePaymentLinkResponse, err error)
		MidtransCheckStatusTransaction(ctx context.Context, orderID string) (res models.MidtransCheckStatusResponse, err error)
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

func (m *MidtransRepoImpl) GeneratePaymentLink(ctx context.Context, req models.MidtransGeneratePaymentLinkRequest) (res models.MidtransGeneratePaymentLinkResponse, err error) {
	generatePaymentLinkURL := os.Getenv("MIDTRANS_BASE_URL") + "/v1/payment-links"

	requestPayload := models.MidtransGeneratePaymentLinkRequest{
		TransactionDetails: req.TransactionDetails,
		UsageLimit:         1,
		Expiry: models.Expiry{
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

func (m *MidtransRepoImpl) MidtransCheckStatusTransaction(ctx context.Context, orderID string) (res models.MidtransCheckStatusResponse, err error) {
	checkStatusURL := os.Getenv("MIDTRANS_BASE_URL") + fmt.Sprintf("/v2/%s/status", orderID)

	request, err := http.NewRequest("GET", checkStatusURL, nil)
	if err != nil {
		slog.ErrorContext(ctx, "[repo][midtrans][CheckStatusTransaction] Error create new request: ", err)
		return res, err
	}

	request.Header.Set("Authorization", generateTokenAuthorization())
	request.Header.Set("Content-Type", "application/json")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		slog.ErrorContext(ctx, "[repo][midtrans][CheckStatusTransaction] Error sending request: ", err)
		return res, err
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		slog.ErrorContext(ctx, "[repo][midtrans][CheckStatusTransaction] Error reading response body: ", err)
		return res, err
	}

	err = json.Unmarshal(body, &res)
	if err != nil {
		slog.ErrorContext(ctx, "[repo][midtrans][CheckStatusTransaction] Error unmarshal response body: ", err)
		return res, err
	}

	return res, nil
}
