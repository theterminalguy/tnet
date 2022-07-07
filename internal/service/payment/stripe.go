package payment

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/10hourlabs/tentn/internal/repository"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stripe/stripe-go/webhook"
)

type Data struct {
	Object struct {
		ID                 string   `json:"id"`
		Object             string   `json:"object"`
		AmountSubtotal     int      `json:"amount_subtotal"`
		AmountTotal        int      `json:"amount_total"`
		ClientReferenceID  string   `json:"client_reference_id"`
		Currency           string   `json:"currency"`
		ExpiresAt          int      `json:"expires_at"`
		Livemode           bool     `json:"livemode"`
		Locale             string   `json:"locale"`
		Mode               string   `json:"mode"`
		PaymentIntent      string   `json:"payment_intent"`
		PaymentLink        string   `json:"payment_link"`
		PaymentMethodTypes []string `json:"payment_method_types"`
		PaymentStatus      string   `json:"payment_status"`
		Metadata           struct {
			JobID string `json:"jobID"`
		} `json:"metadata"`
		Status     string `json:"status"`
		SuccessURL string `json:"success_url"`
	} `json:"object"`
}

type StripePaymentResponse struct {
	ID       string `json:"id"`
	Object   string `json:"object"`
	Created  int    `json:"created"`
	Data     Data   `json:"data"`
	Livemode bool   `json:"livemode"`
	Type     string `json:"type"`
}

type StripePayment struct {
	repo *repository.PaymentRepository
}

func NewStripePayment() *StripePayment {
	return &StripePayment{
		repo: repository.NewPaymentRepository(),
	}
}

func (p *StripePayment) GenerateLink(jobcollectionID uuid.UUID, recruiterID uuid.UUID) (string, error) {
	return "", nil
}

func (p *StripePayment) Pay(req echo.Context) (string, error) {
	const MaxBodyBytes = int64(65536)
	req.Request().Body = http.MaxBytesReader(req.Response().Writer, req.Request().Body, MaxBodyBytes)
	payload, err := ioutil.ReadAll(req.Request().Body)
	if err != nil {
		return "", err
	}

	endpointSecret := os.Getenv("STRIPE_ENDPOINT_SECRET")
	if endpointSecret == "" {
		return "", errors.New("invalid endpoint secret")
	}
	event, err := webhook.ConstructEvent(payload, req.Request().Header.Get("Stripe-Signature"),
		endpointSecret)

	if err != nil {
		return "", err
	}

	var data repository.PaymentParams
	var response = &StripePaymentResponse{}

	if event.Type == "checkout.session.completed" {
		j, err := json.Marshal(event)
		if err != nil {
			return "", err
		}
		err = json.Unmarshal(j, response)
		if err != nil {
			return "", err
		}

		if response.Data.Object.Metadata.JobID == "" {
			return "error occured while processing payment", nil
		}

		data = repository.PaymentParams{
			Amount:   float32(response.Data.Object.AmountTotal) / 100, // stripe amount is in cent
			Status:   response.Data.Object.PaymentStatus,
			RefId:    response.Data.Object.ID,
			Message:  "Successful",
			Currency: response.Data.Object.Currency,
			JobID:    uuid.MustParse(response.Data.Object.Metadata.JobID),
			Payload:  []string{string(j)},
		}

		rID := uuid.MustParse(response.Data.Object.Metadata.JobID)
		_, err = p.repo.Update(rID, data)
		if err != nil {
			return "", err
		}
		return "successful", nil
	}
	return "payment failed", nil
}
