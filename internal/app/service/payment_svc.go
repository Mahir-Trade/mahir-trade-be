package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log/slog"
	"mahir-trade-be/internal/app/models"
	"mahir-trade-be/internal/app/repo/midtrans"
	"mahir-trade-be/internal/app/repo/postgres"
	"mahir-trade-be/pkg/middleware"
	"net/http"
	"strings"
	"time"

	"go.uber.org/dig"
)

type (
	PaymentSvc interface {
		GeneratePaymentLink(ctx context.Context, packageID int64) (resp models.DefaultResponse, err error)
	}

	PaymentSvcImpl struct {
		dig.In

		PackageRepo postgres.PackageRepo
		UserRepo    postgres.UserRepo
		Midtrans    midtrans.MidtransRepo
	}
)

func NewPaymentSvc(impl PaymentSvcImpl) PaymentSvc {
	return &impl
}

func (p *PaymentSvcImpl) GeneratePaymentLink(ctx context.Context, packageID int64) (resp models.DefaultResponse, err error) {
	{
		resp.Code = http.StatusOK
		resp.Message = "Success"
	}

	packageData, err := p.PackageRepo.GetPackageByID(ctx, packageID)
	if err != nil {
		resp.Code = http.StatusBadRequest
		resp.Message = "bad request"
		resp.Error = err.Error()
		slog.ErrorContext(ctx, "[service][GeneratePaymentLink] while GetPackageByID err : ", err)

		return resp, err
	}

	if packageData.ID == 0 {
		resp.Code = http.StatusNotFound
		resp.Message = "not found"
		resp.Error = "package not found"
		slog.ErrorContext(ctx, "[service][GeneratePaymentLink] package not found")

		return resp, err
	}

	currUser, ok := ctx.Value(middleware.UserData).(middleware.UserCtxReq)
	if !ok {
		resp.Code = http.StatusUnauthorized
		resp.Message = "unauthorized"
		resp.Error = "user not found"
		slog.ErrorContext(ctx, "[service][GeneratePaymentLink] user not found")

		return resp, err
	}

	user, err := p.UserRepo.GetUserByID(ctx, currUser.UserID)
	if err != nil {
		resp.Code = http.StatusBadRequest
		resp.Message = "bad request"
		resp.Error = err.Error()
		slog.ErrorContext(ctx, "[service][GeneratePaymentLink] while GetUserByID err : ", err)

		return resp, err
	}

	paymentCode, err := generatePaymentCode(user.PhoneNumber)
	if err != nil {
		resp.Code = http.StatusInternalServerError
		resp.Message = "internal server error"
		resp.Error = err.Error()
		slog.ErrorContext(ctx, "[service][GeneratePaymentLink] while generatePaymentCode err : ", err)

		return resp, err
	}

	req := midtrans.GeneratePaymentLinkRequest{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:     paymentCode,
			GrossAmount: packageData.Price,
		},
		ItemDetails: []midtrans.ItemDetails{
			{
				Name:     fmt.Sprintf("Package %d month", packageData.DurationInMonth),
				Price:    packageData.Price,
				Quantity: 1,
			},
		},
		CustomerDetails: midtrans.CustomerDetails{
			FirstName: user.Fullname,
			Email:     currUser.Email,
			Phone:     user.PhoneNumber,
		},
	}

	paymentURL, err := p.Midtrans.GeneratePaymentLink(ctx, req)
	if err != nil {
		resp.Code = http.StatusInternalServerError
		resp.Message = "internal server error"
		resp.Error = err.Error()
		slog.ErrorContext(ctx, "[service][GeneratePaymentLink] while GeneratePaymentLink err : ", err)

		return resp, err
	}

	resp.Data = paymentURL

	// TODO: save to orders table

	return
}

func generatePaymentCode(phoneNumber string) (string, error) {
	randomBytes := make([]byte, 5)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	randomString := strings.ToUpper(hex.EncodeToString(randomBytes))
	phoneLast5 := phoneNumber[len(phoneNumber)-5:]
	secondsInDay := time.Now().Hour()*3600 + time.Now().Minute()*60 + time.Now().Second()
	timeComponent := fmt.Sprintf("%05d", secondsInDay)
	paymentCode := fmt.Sprintf("MT-%s%s%s", randomString[:3], phoneLast5[:2], timeComponent[:5])

	return paymentCode, nil
}
