package payment

import (
	"github.com/10hourlabs/tenlog"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type JobPaymentManager interface {
	Pay(c echo.Context) (string, error)
	GenerateLink(jobID uuid.UUID) (string, error)
}

type PaymentService struct{}

func NewPaymentService(payment_type string) JobPaymentManager {
	if payment_type == "" {
		tenlog.Error("Unable to process any payment driver")
	}
	payments := map[string]JobPaymentManager{
		"stripe": NewStripePayment(),
	}
	return payments[payment_type]
}
