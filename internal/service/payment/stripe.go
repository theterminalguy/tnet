package payment

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/10hourlabs/tenlog"
	"github.com/10hourlabs/tentn/internal/repository"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

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
		return "", err
	}

	result, _ := ioutil.ReadAll(resp.Body)
	defer req.Body.Close()
	var g PaymentLinkResponse
	err = json.Unmarshal(result, &g)
	if err != nil {
		return "", nil
	}

	generateLink := repository.JobPaymentParams{
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
