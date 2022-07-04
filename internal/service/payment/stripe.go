package payment

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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
			RecruiterID string `json:"recruiterID,omitempty"`
		} `json:"metadata,omitempty"`
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

type GeneratePaymentLinkResponse struct {
	ID       string `json:"id,omitempty"`
	Object   string `json:"object,omitempty"`
	Active   bool   `json:"active,omitempty"`
	Livemode bool   `json:"livemode,omitempty"`
	Metadata struct {
	} `json:"metadata,omitempty"`
	OnBehalfOf         interface{} `json:"on_behalf_of,omitempty"`
	PaymentIntentData  interface{} `json:"payment_intent_data,omitempty"`
	PaymentMethodTypes interface{} `json:"payment_method_types,omitempty"`
	SubscriptionData   interface{} `json:"subscription_data,omitempty"`
	URL                string      `json:"url"`
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
	url := "https://api.stripe.com/v1/payment_links"
	priceKey := os.Getenv("STRIPE_PRODUCT_KEY")
	apikey := os.Getenv("STRIPE_API_KEY")

	if priceKey == "" {
		log.Panic("Stripe product key is required")
	}

	data := []byte(fmt.Sprintf("line_items[0][price]=%s&line_items[0][quantity]=1&metadata[recruiterID]=%s", priceKey, recruiterID))
	responseBody := bytes.NewBuffer(data)
	req, err := http.NewRequest("POST", url, responseBody)
	if err != nil {
		return "", err
	}
	req.SetBasicAuth(apikey, "")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	result, _ := ioutil.ReadAll(resp.Body)
	defer req.Body.Close()
	var g GeneratePaymentLinkResponse
	err = json.Unmarshal(result, &g)
	if err != nil {
		return "", nil
	}

	generateLink := repository.PaymentParams{
		Status:          "not_paid",
		UserId:          recruiterID,
		Message:         "Pending",
		PaymentLink:     g.URL,
		JobCollectionID: jobcollectionID,
	}

	r, err := p.repo.Create(generateLink)
	if err != nil {
		return "", err
	}

	return *r.PaymentLink, nil
}

func (p *StripePayment) Pay(req echo.Context) (string, error) {
	const MaxBodyBytes = int64(65536)
	req.Request().Body = http.MaxBytesReader(req.Response().Writer, req.Request().Body, MaxBodyBytes)
	payload, err := ioutil.ReadAll(req.Request().Body)
	if err != nil {
		return "", err
	}
	// This is your Stripe CLI webhook secret for testing your endpoint locally.
	endpointSecret := "whsec_395e25c1e34329e5499a86a52d6a7b2654805443129e957820bf364937493766"
	// Pass the request body and Stripe-Signature header to ConstructEvent, along
	// with the webhook signing key.
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

		data = repository.PaymentParams{
			Amount:   float32(response.Data.Object.AmountTotal) / 100, // stripe amount is in cent
			Status:   response.Data.Object.PaymentStatus,
			RefId:    response.ID,
			UserId:   uuid.MustParse(response.Data.Object.Metadata.RecruiterID),
			Message:  "Successful",
			Currency: response.Data.Object.Currency,
		}

		rID := uuid.MustParse(response.Data.Object.Metadata.RecruiterID)
		_, err = p.repo.Update(rID, data)
		if err != nil {
			return "", err
		}
		return "successful", nil
	}

	//send mail to admin for successful payment.
	if event.Type == "charge.failed" {

		return "payment failed", nil
	}

	return "error occured", nil
}
