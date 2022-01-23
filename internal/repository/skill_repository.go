package repository

import (
	"strings"
	"time"

	"github.com/10hourlabs/tentn/ent"
	"github.com/10hourlabs/tentn/ent/predicate"
	"github.com/10hourlabs/tentn/ent/skill"
	"github.com/10hourlabs/tentn/util/collection"
	"github.com/google/uuid"
)

type SkillQuerier interface {
	GetAll() ([]*ent.Skill, error)
	GetAllForTalent(talentID int) ([]*ent.Skill, error)
	GetByUUID(id uuid.UUID) (*ent.Skill, error)
	Create(p SkillParams) (*ent.Skill, error)
	Update(id uuid.UUID, p SkillParams) (*ent.Skill, []error)
	DeleteByUUID(id uuid.UUID) error
}

type SkillRepository struct{}

type SkillParams struct {
	TalentUUID uuid.UUID `json:"talent_uuid" validate:"required"`

	// Talent can specify years of experience in decimal where 1.5 equals 1 and a half year
	YearsOfExperience float32 `json:"years_of_experience" validate:"gte=1.0"`
	Preferred         bool    `json:"preferred"`

	// Talent should add details on this specific skills
	// how they have used them in the past, things they've built or done with it
	Note string `json:"note" validate:"required"`

	Name string `json:"name" validate:"required"`
}

func NewSkillRepository() *SkillRepository {
	return &SkillRepository{}
}

func (*SkillRepository) Filter(prd ...predicate.Skill) ([]*ent.Skill, error) {
	skills, err := dBConn.Skill.Query().
		Where(prd...).
		All(dBContext)
	if err != nil {
		return nil, err
	}
	return skills, nil
}

func (*SkillRepository) GetAll() ([]*ent.Skill, error) {
	records, err := dBConn.Skill.Query().
		Where(skill.DeletedAtIsNil()).
		All(dBContext)
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (*SkillRepository) GetAllByTalentUUID(TalentUUID uuid.UUID) ([]*ent.Skill, error) {
	a, err := NewTalentRepository().GetByUUID(TalentUUID)
	if err != nil {
		return nil, err
	}
	records, err := a.QuerySkills().All(dBContext)
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (*SkillRepository) GetAllForTalent(talentID int) ([]*ent.Skill, error) {
	records, err := dBConn.Skill.Query().
		Where(skill.And(
			skill.TalentID(talentID),
			skill.DeletedAtIsNil())).
		All(dBContext)
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (*SkillRepository) GetByUUID(id uuid.UUID) (*ent.Skill, error) {
	record, err := dBConn.Skill.Query().
		Where(skill.UUIDEQ(id)).
		Only(dBContext)
	if err != nil {
		return nil, err
	}
	if record.DeletedAt != nil {
		return nil, ErrRecordDeleted
	}
	return record, nil
}

func (*SkillRepository) Create(p SkillParams) (*ent.Skill, error) {
	err := validateParams(p)
	if err != nil {
		return nil, err
	}
	a, err := NewTalentRepository().GetByUUID(p.TalentUUID)
	if err != nil {
		return nil, err
	}
	record, err := dBConn.Skill.
		Create().
		SetTalentID(a.ID).
		SetName(strings.ToLower(p.Name)).
		SetPreferred(p.Preferred).
		SetYearsOfExperience(p.YearsOfExperience).
		SetNote(p.Note).
		Save(dBContext)
	if err != nil {
		return nil, err
	}
	return record, err
}

func (r *SkillRepository) Update(id uuid.UUID, p SkillParams) (*ent.Skill, []error) {
	err := validateParams(p, "TalentUUID")
	if err != nil {
		return nil, []error{err}
	}
	record, err := r.GetByUUID(id)
	if err != nil {
		return nil, []error{err}
	}

	var vldErrs []error
	bldr := record.Update()

	// Set and Validate YearsOfExperience if provided
	if vldErr := setNillableYearsOfExperience(&p.YearsOfExperience, func(v *float32) error {
		err := validateParams(p, "YearsOfExperience")
		if err != nil {
			return err
		}
		bldr.SetYearsOfExperience(p.YearsOfExperience)
		return nil
	}); vldErr != nil {
		vldErrs = append(vldErrs, vldErr)
	}

	// Set and Validate Name if provided
	if vldErr := setNillableStringField(p.Name, func(v string) error {
		err := validateParams(p, "Name")
		if err != nil {
			return err
		}
		bldr.SetName(strings.ToLower(v))
		return nil
	}); vldErr != nil {
		vldErrs = append(vldErrs, vldErr)
	}

	// Set and Validate Preferred if provided
	if vldErr := setNillableBoolField(p.Preferred, func(v bool) error {
		err := validateParams(p, "Preferred")
		if err != nil {
			return err
		}
		bldr.SetPreferred(v)
		return nil
	}); vldErr != nil {
		vldErrs = append(vldErrs, vldErr)
	}

	// Set and Validate Note if provided
	if vldErr := setNillableStringField(p.Note, func(v string) error {
		err := validateParams(p, "Note")
		if err != nil {
			return err
		}
		bldr.SetNote(v)
		return nil
	}); vldErr != nil {
		vldErrs = append(vldErrs, vldErr)
	}

	// Return all validation errors at once
	// this prevents the client from making several round trips to the server
	if collection.HasAny(vldErrs) {
		return nil, vldErrs
	}

	record, err = bldr.Save(dBContext)
	if err != nil {
		return nil, []error{err}
	}

	return record, nil
}

func (r *SkillRepository) DeleteByUUID(id uuid.UUID) error {
	record, err := r.GetByUUID(id)
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
