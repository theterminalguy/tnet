package repository

import (
	"errors"
	"log"

	"github.com/10hourlabs/tentn/ent"
	"github.com/10hourlabs/tentn/ent/payment"
	"github.com/google/uuid"
)

type PaymentRepository struct{}

type PaymentParams struct {
	Amount      float32   `json:"amount,omitempty"`
	Status      string    `json:"status,omitempty"`
	RefId       string    `json:"ref_id,omitempty"`
	Message     string    `json:"message" validate:"required"`
	Currency    string    `json:"currency,omitempty"`
	PaymentLink string    `json:"payment_link,omitempty"`
	JobID       uuid.UUID `json:"job_id" validate:"required"`
	Payload     []string  `json:"payload,omitempty"`
}

func NewPaymentRepository() *PaymentRepository {
	return &PaymentRepository{}
}

func (*PaymentRepository) GetAll() ([]*ent.Payment, error) {
	return nil, nil
}

func (*PaymentRepository) GetByJobID(id uuid.UUID) (*ent.Payment, error) {
	record, err := dBConn.Payment.Query().
		Where(payment.JobID(id)).
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
			SetStatus(payment.Status(p.Status)).
			SetMessage(p.Message).
			SetJobID(p.JobID).
			Save(dBContext)

		if err != nil {
			return nil, err
		}

		return record, nil
	}

	return nil, errors.New("error occured while processing payment link")
}

func (px *PaymentRepository) Update(id uuid.UUID, p PaymentParams) (*ent.Payment, error) {
	return nil, nil
}

func (*PaymentRepository) DeleteByID(id uuid.UUID) error {
	return nil
}
