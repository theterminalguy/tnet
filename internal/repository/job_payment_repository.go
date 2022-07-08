package repository

import (
	"errors"
	"log"

	"github.com/10hourlabs/tentn/ent"
	"github.com/10hourlabs/tentn/ent/jobpayment"
	"github.com/google/uuid"
)

type JobPaymentRepository struct{}

type JobPaymentParams struct {
	Amount      float32   `json:"amount,omitempty"`
	Status      string    `json:"status,omitempty"`
	RefId       string    `json:"ref_id,omitempty"`
	Message     string    `json:"message" validate:"required"`
	Currency    string    `json:"currency,omitempty"`
	PaymentLink string    `json:"payment_link,omitempty"`
	JobID       uuid.UUID `json:"job_id" validate:"required"`
	Payload     []string  `json:"payload,omitempty"`
}

func NewJobPaymentRepository() *JobPaymentRepository {
	return &JobPaymentRepository{}
}

func (*JobPaymentRepository) GetAll() ([]*ent.JobPayment, error) {
	return nil, nil
}

func (*JobPaymentRepository) GetByJobID(id uuid.UUID) (*ent.JobPayment, error) {
	record, err := dBConn.JobPayment.Query().
		Where(jobpayment.JobID(id)).
		Only(dBContext)

	if err != nil {
		return record, err
	}
	if record.DeletedAt != nil {
		return nil, ErrRecordDeleted
	}
	return record, nil
}

func (*JobPaymentRepository) Create(p JobPaymentParams) (*ent.JobPayment, error) {
	// var record *ent.Client
	err := ValidateParams(p)
	if err != nil {
		log.Panic(err)
	}
	// Store Generated payment link
	if p.PaymentLink != "" {
		record, err := dBConn.JobPayment.
			Create().
			SetPaymentLink(p.PaymentLink).
			SetStatus(jobpayment.Status(p.Status)).
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

func (px *JobPaymentRepository) Update(id uuid.UUID, p JobPaymentParams) (*ent.JobPayment, error) {
	return nil, nil
}

func (*JobPaymentRepository) DeleteByID(id uuid.UUID) error {
	return nil
}
