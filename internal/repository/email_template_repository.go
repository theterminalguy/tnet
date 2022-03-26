package repository

import (
	"time"

	"github.com/10hourlabs/tentn/ent"
	"github.com/10hourlabs/tentn/ent/emailtemplate"
	"github.com/google/uuid"
)

type EmailTemplateRepository struct{}

type EmailTemplateParams struct {
	UserID  uuid.UUID            `json:"user_id" validate:"required"`
	Subject string               `json:"subject" validate:"required"`
	Body    string               `json:"body" validate:"required"`
	Status  emailtemplate.Status `json:"status" validate:"required"`
	Cc      []string             `json:"cc"`
	Bcc     []string             `json:"bcc"`
}

func NewEmailTemplateRepository() *EmailTemplateRepository {
	return &EmailTemplateRepository{}
}

func (*EmailTemplateRepository) GetAll() ([]*ent.EmailTemplate, error) {
	records, err := dBConn.EmailTemplate.Query().
		Where(emailtemplate.DeletedAtIsNil()).
		All(dBContext)
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (*EmailTemplateRepository) GetEmailForRecruiter(userid uuid.UUID, status string) (*ent.EmailTemplate, error) {
	record, err := dBConn.EmailTemplate.Query().
		Where(
			emailtemplate.UserIDEQ(userid),
			emailtemplate.StatusEQ(emailtemplate.Status(status)),
		).
		Only(dBContext)
	if err != nil {
		return nil, err
	}
	if record.DeletedAt != nil {
		return nil, ErrRecordDeleted
	}
	return record, nil
}

func (*EmailTemplateRepository) GetByID(id uuid.UUID) (*ent.EmailTemplate, error) {
	record, err := dBConn.EmailTemplate.Query().
		Where(emailtemplate.ID(id)).
		Only(dBContext)
	if err != nil {
		return nil, err
	}
	if record.DeletedAt != nil {
		return nil, ErrRecordDeleted
	}
	return record, nil
}

func (*EmailTemplateRepository) Create(p EmailTemplateParams) (*ent.EmailTemplate, error) {
	err := validateParams(p)
	if err != nil {
		return nil, err
	}
	record, err := dBConn.EmailTemplate.
		Create().
		SetBcc(p.Bcc).
		SetCc(p.Cc).
		SetBody(p.Body).
		SetUserID(p.UserID).
		SetStatus(p.Status).
		SetSubject(p.Subject).
		SetFrom("balogun.tobi@10hourlabs.com").
		Save(dBContext)
	if err != nil {
		return nil, err
	}
	return record, err
}

func (r *EmailTemplateRepository) Update(id uuid.UUID, p EmailTemplateParams) (*ent.EmailTemplate, error) {
	err := validateParams(p)
	if err != nil {
		return nil, err
	}
	record, err := r.GetByID(id)
	if err != nil {
		return nil, err
	}
	_, err = record.Update().
		// TODO: set other fields here
		Save(dBContext)
	if err != nil {
		return nil, err
	}
	return record, nil
}

func (r *EmailTemplateRepository) DeleteByID(id uuid.UUID) error {
	record, err := r.GetByID(id)
	if err != nil {
		return err
	}
	_, err = record.Update().
		SetDeletedAt(time.Now()).
		Save(dBContext)
	if err != nil {
		return err
	}
	return nil
}
