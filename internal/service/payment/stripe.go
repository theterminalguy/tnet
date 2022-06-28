package payment

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/10hourlabs/tentn/logger"
	"github.com/labstack/echo/v4"
	"github.com/stripe/stripe-go/webhook"
)

type StripePayment struct {
	Id          string
	Object      string
	Api_version string
}

func NewStripePayment() *StripePayment {
	return &StripePayment{}
}

func (*StripePayment) Pay(req echo.Context) error {
	const MaxBodyBytes = int64(65536)
	req.Request().Body = http.MaxBytesReader(req.Response().Writer, req.Request().Body, MaxBodyBytes)
	payload, err := ioutil.ReadAll(req.Request().Body)
	if err != nil {
		// fmt.Fprintf(os.Stderr, "Error reading request body: %v\n", err)
		// w.WriteHeader(http.StatusServiceUnavailable)
		return err
	}
	// This is your Stripe CLI webhook secret for testing your endpoint locally.
	endpointSecret := "whsec_395e25c1e34329e5499a86a52d6a7b2654805443129e957820bf364937493766"
	// Pass the request body and Stripe-Signature header to ConstructEvent, along
	// with the webhook signing key.
	event, err := webhook.ConstructEvent(payload, req.Request().Header.Get("Stripe-Signature"),
		endpointSecret)

	if err != nil {
		// fmt.Fprintf(os.Stderr, "Error verifying webhook signature: %v\n", err)
		// w.WriteHeader(http.StatusBadRequest) // Return a 400 error on a bad signature
		return err
	}

	json, _ := json.Marshal(event)
	logger.Debug(string(json))

	fmt.Println(string(payload))

	// Unmarshal the event data into an appropriate struct depending on its Type
	switch event.Type {
	case "payment_intent.succeeded":
		// Then define and call a function to handle the event payment_intent.succeeded
		// ... handle other event types

	default:
		fmt.Fprintf(os.Stderr, "Unhandled event type: %s\n", event.Type)
	}

	return nil
}
