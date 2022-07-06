package repository

import (
	"errors"
	"fmt"
	"log"

	"github.com/10hourlabs/tentn/ent"
	"github.com/10hourlabs/tentn/ent/payment"
	"github.com/google/uuid"
)

type PaymentRepository struct{}

type PaymentParams struct {
	Amount          float32 `json:"amount,omitempty"`
	Status          string  `json:"status,omitempty"`
	RefId           string  `json:"ref_id,omitempty"`
	UserId          uuid.UUID
	Message         string    `json:"message" validate:"required"`
	Currency        string    `json:"currency,omitempty"`
	PaymentLink     string    `json:"payment_link,omitempty"`
	JobCollectionID uuid.UUID // this is a relationship to job collection entity
	Payload         []string  `json:"payload,omitempty"`
}

func NewPaymentRepository() *PaymentRepository {
	return &PaymentRepository{}
}

func (*PaymentRepository) GetAll() ([]*ent.Payment, error) {
	return nil, nil
}

func (*PaymentRepository) GetByUserID(id uuid.UUID) (*ent.Payment, error) {
	record, err := dBConn.Payment.Query().
		Where(payment.UserID(id)).
		Only(dBContext)

	if err != nil {
		return record, err
	}
	if record.DeletedAt != nil {
		return nil, ErrRecordDeleted
	}
	return record, nil
}

func (*PaymentRepository) Create(p PaymentParams) (*ent.Payment, error) {
	// var record *ent.Client
	err := ValidateParams(p)
	if err != nil {
		log.Panic(err)
	}
	// Store Generated payment link
	if p.PaymentLink != "" {
		record, err := dBConn.Payment.
			Create().
			SetPaymentLink(p.PaymentLink).
			SetUserID(p.UserId).
			SetStatus(payment.Status(p.Status)).
			SetMessage(p.Message).
			Save(dBContext)

		if err != nil {
			return nil, err
		}

		return record, nil
	}

	return nil, errors.New("error occured while processing payment link")
}

func (px *PaymentRepository) Update(id uuid.UUID, p PaymentParams) (*ent.Payment, error) {
	if p.UserId == uuid.Nil {
		return nil, errors.New("recruiterID is missing in the payload")
	}
	client, err := px.GetByUserID(p.UserId)
	if err != nil {
		fmt.Println(client)
		return nil, err
	}

	record, err := client.Update().
		SetRefID(p.RefId).
		SetMessage(p.Message).
		SetCurrency(p.Currency).
		SetPayload(p.Payload).
		SetAmount(float64(p.Amount)).
		SetStatus(payment.Status(p.Status)).
		Save(dBContext)

	if err != nil {
		return nil, err
	}
	return record, nil
}

func (*PaymentRepository) DeleteByID(id uuid.UUID) error {
	return nil
}
