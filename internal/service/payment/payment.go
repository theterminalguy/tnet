package payment

import (
	"log"

	"github.com/labstack/echo/v4"
)

type IPayment interface {
	Pay(c echo.Context) error
}

type PaymentService struct{}

func NewPaymentService(payment_type string) IPayment {
	if payment_type == "" {
		log.Panic("Unable to process any payment driver")
	}
	payments := map[string]IPayment{
		"stripe": NewStripePayment(),
	}
	return payments[payment_type]
}
