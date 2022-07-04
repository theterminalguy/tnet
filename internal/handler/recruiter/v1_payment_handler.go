package recruiter

import (
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
	// currentRecruiter := c.Get(oneword.CurrentRecruiter).(*scope.RecruiterScope)
	driver := os.Getenv("PAYMENT_DRIVER")
	pay := payment.NewPaymentService(driver)
	resp, err := pay.Pay(c)
	if err != nil {
		c.String(http.StatusBadRequest, "Error occured while processing payment")
	}

	// Move generated link to create jd endpoint
	// resp, err := payment.NewStripePayment().GenerateLink(uuid.MustParse("601944a8-3833-4bb4-bf29-4298cc1cc5fb"), uuid.MustParse("601944a8-3833-4bb4-bf29-4298cc1cc5fb"))

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
