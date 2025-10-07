package recruiter

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/theterminalguy/tnet/internal/service/payment"
)

type PaymentHandler struct{}

func NewV1PaymentHandler() *PaymentHandler {
	return &PaymentHandler{}
}

func (*PaymentHandler) CreateOne(c echo.Context) error {
	driver := os.Getenv("PAYMENT_DRIVER")
	pay := payment.NewPaymentService(driver)
	resp, err := pay.Pay(c)
	if err != nil {
		c.String(http.StatusBadRequest, "Error occured while processing payment")
	}

	return c.String(http.StatusOK, resp)
}

func (*PaymentHandler) ReadAll(c echo.Context) error {
	return nil
}

func (*PaymentHandler) DeleteOne(c echo.Context) error {
	return nil
}

func (*PaymentHandler) ReadByID(c echo.Context) error {
	return nil
}

func (*PaymentHandler) Search(c echo.Context) error {
	return nil
}

func (*PaymentHandler) UpdateByID(c echo.Context) error {
	return nil
}
