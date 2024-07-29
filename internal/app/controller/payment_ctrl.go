package controller

import (
	"log/slog"
	"mahir-trade-be/internal/app/models"
	"mahir-trade-be/internal/app/service"
	"mahir-trade-be/internal/app/service/utils"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/dig"
)

type (
	PaymentCtrl interface {
		CreatePayment(ec echo.Context) error
		PaymentLinkCallback(ec echo.Context) error
	}

	PaymentCtrlImpl struct {
		dig.In

		PaymentSvc service.PaymentSvc
	}
)

func NewPaymentCtrl(impl PaymentCtrlImpl) PaymentCtrl {
	return &impl
}

func (p *PaymentCtrlImpl) CreatePayment(ec echo.Context) error {
	ctx := ec.Request().Context()

	defer func() {
		if r := recover(); r != nil {
			slog.Error("CreateGroup - something went wrong", r)
		}
	}()

	type req struct {
		PackageID int64 `json:"package_id" validate:"required"`
	}

	var request req

	if err := ec.Bind(&request); err != nil {
		slog.Error("CreatePayment - something went wrong, bind error", err)
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   "invalid pagination request",
		})
	}

	validate := utils.Validate
	err := validate.Struct(request)
	if err != nil {
		slog.Error("CreatePayment - something went wrong, validation error", err)
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   "invalid pagination request",
		})
	}

	resp, err := p.PaymentSvc.GeneratePaymentLink(ctx, request.PackageID)
	if err != nil {
		slog.Error("CreatePayment - something went wrong, GeneratePaymentLink error", err)
		return ec.JSON(http.StatusBadRequest, resp)
	}

	return ec.JSON(http.StatusOK, resp)
}

func (p *PaymentCtrlImpl) PaymentLinkCallback(ec echo.Context) error {
	ctx := ec.Request().Context()

	defer func() {
		if r := recover(); r != nil {
			slog.Error("PaymentLinkCallback - something went wrong", r)
		}
	}()

	var req models.MidtransCallbackRequest

	if err := ec.Bind(&req); err != nil {
		slog.Error("PaymentLinkCallback - something went wrong, bind error", err)
		return ec.JSON(http.StatusBadRequest, models.DefaultResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   "invalid pagination request",
		})
	}

	err := p.PaymentSvc.MidtransPaymentLinkNotification(ctx, req)
	if err != nil {
		slog.Error("PaymentLinkCallback - something went wrong, MidtransPaymentLinkNotification error", err)
		return ec.JSON(http.StatusOK, "success")
	}

	return ec.JSON(http.StatusOK, "success")
}
