package repository

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/theterminalguy/tnet/ent"
	"github.com/theterminalguy/tnet/ent/jobpayment"
)

type JobPaymentRepository struct{}

type JobPaymentParams struct {
	Amount      float32   `json:"amount,omitempty"`
	PaidTo      time.Time `json:"paid_to,omitempty"`
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

func (*JobPaymentRepository) HasPaid(refID string, jobID uuid.UUID) bool {
	record := dBConn.JobPayment.Query().
		Where(jobpayment.RefID(refID)).
		Where(jobpayment.And(jobpayment.JobIDEQ(jobID))).
		Where(jobpayment.And(jobpayment.PaidToNEQ(time.Time{}))).
		CountX(dBContext)

	return record != 0
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
	err := ValidateParams(p)
	if err != nil {
		return nil, err
	}
	// Store Generated payment link
	if p.PaymentLink != "" {
		record, err := dBConn.JobPayment.
			Create().
			SetPaymentLink(p.PaymentLink).
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
	if id == uuid.Nil {
		return nil, errors.New("jobID is missing in the payload")
	}
	if px.HasPaid(p.RefId, id) {
		return nil, errors.New("duplicate transaction")
	}

	client, err := px.GetByJobID(id)
	if err != nil {
		return nil, err
	}

	record, err := client.Update().
		SetRefID(p.RefId).
		SetMessage(p.Message).
		SetCurrency(p.Currency).
		SetPayload(p.Payload).
		SetAmount(float64(p.Amount)).
		SetPaidTo(p.PaidTo).
		Save(dBContext)

	if err != nil {
		return nil, err
	}
	return record, nil
}

func (*JobPaymentRepository) DeleteByID(id uuid.UUID) error {
	return nil
}
