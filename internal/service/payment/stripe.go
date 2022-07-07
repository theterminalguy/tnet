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
)

type GeneratePaymentLinkResponse struct {
	URL string `json:"url"`
}

type StripePayment struct {
	repo *repository.PaymentRepository
}

func NewStripePayment() *StripePayment {
	return &StripePayment{
		repo: repository.NewPaymentRepository(),
	}
}

func (p *StripePayment) GenerateLink(jobID uuid.UUID) (string, error) {
	url := "https://api.stripe.com/v1/payment_links"
	priceKey := os.Getenv("STRIPE_PRODUCT_KEY")
	apikey := os.Getenv("STRIPE_API_KEY")

	if priceKey == "" {
		log.Panic("Stripe product key is required")
	}

	data := []byte(fmt.Sprintf("line_items[0][price]=%s&line_items[0][quantity]=1&metadata[jobID]=%s", priceKey, jobID))
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
		Status:      "not_paid",
		Message:     "Pending",
		PaymentLink: g.URL,
		JobID:       jobID,
	}

	r, err := p.repo.Create(generateLink)
	if err != nil {
		return "", err
	}

	return r.PaymentLink, nil
}

func (p *StripePayment) Pay(req echo.Context) (string, error) {
	return "", nil
}
