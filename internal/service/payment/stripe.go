package payment

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/10hourlabs/tenlog"
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

type PaymentLinkResponse struct {
	URL string `json:"url"`
}

type StripePayment struct {
	repo *repository.JobPaymentRepository
}

func NewStripePayment() *StripePayment {
	return &StripePayment{
		repo: repository.NewJobPaymentRepository(),
	}
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
		tenlog.Error("invalid endpoint secret")
	}
	event, err := webhook.ConstructEvent(payload, req.Request().Header.Get("Stripe-Signature"),
		endpointSecret)

	if err != nil {
		tenlog.Error(err)
	}

	var data repository.JobPaymentParams
	var response = &StripePaymentResponse{}

	if event.Type == "checkout.session.completed" {
		j, err := json.Marshal(event)
		if err != nil {
			tenlog.Error(err)
		}
		err = json.Unmarshal(j, response)
		if err != nil {
			tenlog.Error(err)
		}

		if response.Data.Object.Metadata.JobID == "" {
			tenlog.Error("jobID not found in metadata")
		}

		data = repository.JobPaymentParams{
			Amount:   float32(response.Data.Object.AmountTotal) / 100, // stripe amount is in cent
			PaidTo:   time.Now(),
			RefId:    response.Data.Object.ID,
			Message:  "Successful",
			Currency: response.Data.Object.Currency,
			JobID:    uuid.MustParse(response.Data.Object.Metadata.JobID),
			Payload:  []string{string(j)},
		}

		rID := uuid.MustParse(response.Data.Object.Metadata.JobID)
		_, err = p.repo.Update(rID, data)
		if err != nil {
			tenlog.Error(err)
		}
		return "successful", nil
	}
	return "payment failed", nil
}

func (p *StripePayment) GenerateLink(jobID uuid.UUID) (string, error) {
	endpointUrl := "https://api.stripe.com/v1/payment_links"
	priceKey := os.Getenv("STRIPE_PRODUCT_KEY")
	apikey := os.Getenv("STRIPE_API_KEY")

	if priceKey == "" {
		tenlog.Error("Stripe product key is required")
	}

	params := url.Values{}
	params.Add("line_items[0][price]", priceKey)
	params.Add("line_items[0][quantity]", "1")
	params.Add("metadata[jobID]", jobID.String())
	body := strings.NewReader(params.Encode())
	req, err := http.NewRequest("POST", endpointUrl, body)
	if err != nil {
		return "", err
	}
	req.SetBasicAuth(apikey, "")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		tenlog.Error(err)
	}

	result, _ := ioutil.ReadAll(resp.Body)
	defer req.Body.Close()
	var g PaymentLinkResponse
	err = json.Unmarshal(result, &g)
	if err != nil {
		tenlog.Error(err)
	}

	generateLink := repository.JobPaymentParams{
		Message:     "Pending",
		PaymentLink: g.URL,
		JobID:       jobID,
	}

	r, err := p.repo.Create(generateLink)
	if err != nil {
		tenlog.Error(err)
	}

	return r.PaymentLink, nil
}
