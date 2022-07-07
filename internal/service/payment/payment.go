package payment

import (
	"log"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type PaymentManager interface {
	Pay(c echo.Context) (string, error)
	GenerateLink(jobID uuid.UUID) (string, error)
}

type PaymentService struct{}

func NewPaymentService(payment_type string) PaymentManager {
	if payment_type == "" {
		log.Panic("Unable to process any payment driver")
	}
	payments := map[string]PaymentManager{
		"stripe": NewStripePayment(),
	}
	return payments[payment_type]
}
