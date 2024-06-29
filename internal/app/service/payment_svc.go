package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log/slog"
	"mahir-trade-be/internal/app/models"
	"mahir-trade-be/internal/app/repo/midtrans"
	"mahir-trade-be/internal/app/repo/postgres"
	"mahir-trade-be/internal/app/service/utils"
	"mahir-trade-be/pkg/middleware"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go.uber.org/dig"
)

type (
	PaymentSvc interface {
		GeneratePaymentLink(ctx context.Context, packageID int64) (resp models.DefaultResponse, err error)
		MidtransPaymentLinkNotification(ctx context.Context, req models.MidtransCallbackRequest) error
	}

	PaymentSvcImpl struct {
		dig.In

		PackageRepo        postgres.PackageRepo
		UserRepo           postgres.UserRepo
		Midtrans           midtrans.MidtransRepo
		OrderRepo          postgres.OrderRepo
		TransactionRepo    postgres.TransactionRepo
		UserMembershipRepo postgres.UserMembershipRepo
		GeneralLogRepo     postgres.GeneralLogRepo
	}
)

func NewPaymentSvc(impl PaymentSvcImpl) PaymentSvc {
	return &impl
}

const (
	formattedTime = "2006-01-02 15:04:05"
)

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

	req := models.MidtransGeneratePaymentLinkRequest{
		TransactionDetails: models.TransactionDetails{
			OrderID:     paymentCode,
			GrossAmount: packageData.Price,
		},
		ItemDetails: []models.ItemDetails{
			{
				Name:     fmt.Sprintf("Package %d month", packageData.DurationInMonth),
				Price:    packageData.Price,
				Quantity: 1,
			},
		},
		CustomerDetails: models.CustomerDetails{
			FirstName: user.Username,
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

	orderReq := models.Order{
		UserID:      user.UserID,
		PackageID:   packageData.ID,
		Status:      utils.StatusOrderPending,
		PaymentCode: paymentCode,
		PaymentURL:  paymentURL.PaymentURL,
		CreatedBy:   user.Email,
	}

	_, err = p.OrderRepo.CreateOrder(ctx, orderReq)
	if err != nil {
		resp.Code = http.StatusBadRequest
		resp.Message = "bad request"
		resp.Error = err.Error()
		slog.ErrorContext(ctx, "[service][GeneratePaymentLink] while CreateOrder err : ", err)

		return resp, err
	}

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

func (p *PaymentSvcImpl) MidtransPaymentLinkNotification(ctx context.Context, req models.MidtransCallbackRequest) error {
	slog.InfoContext(ctx, fmt.Sprintf("[service][MidtransPaymentLinkNotification] start process, request : %+v", req))

	if req.OrderID == "" {
		slog.ErrorContext(ctx, "[service][MidtransPaymentLinkNotification] order id is required")
		return fmt.Errorf("order id is required")
	}

	splittedOrderID := strings.Split(req.OrderID, "-")
	if len(splittedOrderID) != 3 {
		slog.ErrorContext(ctx, "[service][MidtransPaymentLinkNotification] invalid order id")
		return fmt.Errorf("invalid order id")
	}

	orderID := splittedOrderID[0] + "-" + splittedOrderID[1]

	order, err := p.OrderRepo.GetOrderByPaymentCode(ctx, orderID)
	if err != nil {
		slog.ErrorContext(ctx, "[service][MidtransPaymentLinkNotification] while GetOrderByPaymentCode err : ", err)
		return err
	}

	if order.ID == 0 {
		slog.ErrorContext(ctx, "[service][MidtransPaymentLinkNotification] order not found")
		return fmt.Errorf("order not found")
	}

	rawData, err := json.Marshal(req)
	if err != nil {
		slog.ErrorContext(ctx, "[service][MidtransPaymentLinkNotification] while Marshal err : ", err)
		return err
	}

	generalLogReq := postgres.GeneralLog{
		UserID:  int(order.UserID),
		RawBody: fmt.Sprintf("MidtransPaymentLinkNotification - %v", string(rawData)),
	}

	err = p.GeneralLogRepo.CreateGeneralLog(ctx, generalLogReq)
	if err != nil {
		slog.ErrorContext(ctx, "[service][MidtransPaymentLinkNotification] while CreateGeneralLog err : ", err)
		return err
	}

	if order.Status == utils.StatusOrderSuccess {
		slog.InfoContext(ctx, "[service][MidtransPaymentLinkNotification] order already success")
		return nil
	}

	orderTransaction, err := p.Midtrans.MidtransCheckStatusTransaction(ctx, req.OrderID)
	if err != nil {
		slog.ErrorContext(ctx, "[service][MidtransPaymentLinkNotification] while MidtransCheckStatusTransaction err : ", err)
		return err
	}

	if orderTransaction.TransactionStatus != req.TransactionStatus {
		slog.ErrorContext(ctx, "[service][MidtransPaymentLinkNotification] transaction status not match")
		return fmt.Errorf("transaction status not match")
	} else {
		amountFloat, err := strconv.ParseFloat(req.GrossAmount, 64)
		if err != nil {
			slog.ErrorContext(ctx, "[service][MidtransPaymentLinkNotification] while ParseFloat err : ", err)
			return err
		}

		if req.TransactionStatus == utils.MidtransStatusSettlement || req.TransactionStatus == utils.MidtransStatusCapture {
			order.Status = utils.StatusOrderSuccess

			transactionReq := models.Transaction{
				OrderID:        order.ID,
				WebhookID:      req.TransactionID,
				Amount:         amountFloat,
				SettlementDate: req.SettlementTime,
				CreatedBy:      "SYSTEM",
			}

			if transactionReq.SettlementDate == "" {
				transactionReq.SettlementDate = req.TransactionTime
			}

			_, err = p.TransactionRepo.CreateTransaction(ctx, transactionReq)
			if err != nil {
				slog.ErrorContext(ctx, "[service][MidtransPaymentLinkNotification] while CreateTransaction err : ", err)
				return err
			}

			packageData, err := p.PackageRepo.GetPackageByID(ctx, order.PackageID)
			if err != nil {
				slog.ErrorContext(ctx, "[service][MidtransPaymentLinkNotification] while GetPackageByID err : ", err)
				return err
			}

			userMembership, err := p.UserMembershipRepo.GetUserMembershipByUserID(ctx, order.UserID)
			if err != nil {
				slog.ErrorContext(ctx, "[service][MidtransPaymentLinkNotification] while GetUserMembershipByUserID err : ", err)
				return err
			}

			if userMembership.ID > 0 {
				currExpiredAt, err := time.Parse(time.RFC3339, userMembership.ExpiredAt)
				if err != nil {
					slog.ErrorContext(ctx, "[service][MidtransPaymentLinkNotification] while Parse expiredAt err : ", err)
					return err
				}

				userMembershipReq := models.UserMembership{
					UserID:             order.UserID,
					PackageID:          order.PackageID,
					ExpiredAt:          currExpiredAt.AddDate(0, int(packageData.DurationInMonth), 0).Format(formattedTime),
					IsMembershipActive: true,
					CreatedBy:          "SYSTEM",
				}

				err = p.UserMembershipRepo.UpdateUserMembershipByUserID(ctx, userMembershipReq)
				if err != nil {
					slog.ErrorContext(ctx, "[service][MidtransPaymentLinkNotification] while UpdateUserMembershipByUserID err : ", err)
					return err
				}
			} else {
				userMembershipReq := models.UserMembership{
					UserID:             order.UserID,
					PackageID:          order.PackageID,
					ExpiredAt:          time.Now().AddDate(0, int(packageData.DurationInMonth), 0).Format(formattedTime),
					IsMembershipActive: true,
					CreatedBy:          "SYSTEM",
				}

				_, err = p.UserMembershipRepo.CreateUserMembership(ctx, userMembershipReq)
				if err != nil {
					slog.ErrorContext(ctx, "[service][MidtransPaymentLinkNotification] while CreateUserMembership err : ", err)
					return err
				}
			}

			slog.InfoContext(ctx, fmt.Sprintf("[service][MidtransPaymentLinkNotification] order %s success", order.PaymentCode))
		} else if req.TransactionStatus != utils.MidtransStatusPending {
			order.Status = utils.StatusOrderFailed
		}

		if order.Status != utils.StatusOrderPending {
			err = p.OrderRepo.UpdateOrderStatus(ctx, order)
			if err != nil {
				slog.ErrorContext(ctx, "[service][MidtransPaymentLinkNotification] while UpdateOrderStatus err : ", err)
				return err
			}

			slog.InfoContext(ctx, fmt.Sprintf("[service][MidtransPaymentLinkNotification] order %s status updated to %s", order.PaymentCode, order.Status))
		}
	}

	return nil
}
