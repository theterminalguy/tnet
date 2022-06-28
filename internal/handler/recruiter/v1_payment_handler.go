package recruiter

import (
	"fmt"
	"net/http"
	"os"

	"github.com/10hourlabs/tentn/internal/service/payment"
	"github.com/labstack/echo/v4"
)

type PaymentHandler struct{}

func NewV1PaymentHandler() *PaymentHandler {
	return &PaymentHandler{}
}

func (*PaymentHandler) CreateOne(c echo.Context) error {
	driver := os.Getenv("PAYMENT_DRIVER")
	pay := payment.NewPaymentService(driver)
	resp := pay.Pay(c)
	fmt.Println(resp)
	return c.String(http.StatusOK, "Hello world")
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
