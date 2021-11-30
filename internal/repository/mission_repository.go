package repository

import (
	"errors"
	"time"

	"github.com/10hourlabs/tentn/ent"
	"github.com/10hourlabs/tentn/ent/mission"
	"github.com/10hourlabs/tentn/util/collection"
	"github.com/10hourlabs/tentn/util/date"
	"github.com/google/uuid"
)

type MissionRepository struct{}

type MissionParams struct {
	PartnerUUID uuid.UUID `json:"partner_uuid" validate:"required"`
	TalentUUID  uuid.UUID `json:"talent_uuid" validate:"required"`
	MissionType string    `json:"mission_type" validate:"required"`
	StartDate   string    `json:"start_date" validate:"datetime=2006-01-02T15:04:05Z07:00"`
	EndDate     string    `json:"end_date"`
}

func NewMissionRepository() *MissionRepository {
	return &MissionRepository{}
}

func (*MissionRepository) GetAll() ([]*ent.Mission, error) {
	records, err := dBConn.Mission.Query().
		Where(mission.DeletedAtIsNil()).
		All(dBContext)
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (*MissionRepository) GetByUUID(id uuid.UUID) (*ent.Mission, error) {
	record, err := dBConn.Mission.Query().
		Where(mission.UUIDEQ(id)).
		Only(dBContext)
	if err != nil {
		return nil, err
	}
	if record.DeletedAt != nil {
		return nil, RecordNotFoundError
	}
	return record, nil
}

func (*MissionRepository) Create(p MissionParams) (*ent.Mission, error) {
	err := validateParams(p)
	if err != nil {
		return nil, err
	}

	a, err := NewTalentRepository().GetByUUID(p.TalentUUID)
	if err != nil {
		return nil, err
	}
	j, err := NewPartnerRepository().GetByUUID(p.PartnerUUID)
	if err != nil {
		return nil, err
	}
	records, err := dBConn.Mission.Query().Where(
		mission.And(
			mission.PartnerID(j.ID),
			mission.TalentID(a.ID),
		)).All(dBContext)
	if collection.HasAny(records) {
		return nil, errors.New("existing mission for partner")
	}

	var sd *time.Time
	var ed *time.Time

	sd, err = date.JSStringToRFC3339(p.StartDate)
	if err != nil {
		return nil, err
	}

	if p.EndDate != "" {
		ed, err = date.JSStringToRFC3339(p.EndDate)
		if err != nil {
			return nil, err
		}

		err = IsEqual(*sd, *ed)
		if err != nil {
			return nil, err
		}
	}

	record, err := dBConn.Mission.
		Create().
		SetPartnerID(a.ID).
		SetTalentID(j.ID).
		SetMissionType(mission.MissionType(p.MissionType)).
		SetStartDate(*sd).
		SetNillableEndDate(ed).
		Save(dBContext)
	if err != nil {
		return nil, err
	}
	return record, err
}

func (r *MissionRepository) Update(id uuid.UUID, p MissionParams) (*ent.Mission, []error) {
	err := validateParams(p, "TalentUUID")
	if err != nil {
		return nil, []error{err}
	}

	err = validateParams(p, "PartnerUUID")
	if err != nil {
		return nil, []error{err}
	}

	record, err := r.GetByUUID(id)
	if err != nil {
		return nil, []error{err}
	}
	var vldErrs []error
	bldr := record.Update()

	sd := date.ToRFC3339(record.StartDate)

	setNillableStringField(p.StartDate, func(v string) error {
		sd = v
		return nil
	})

	period := &DatePeriod{
		StartDate: sd,
		EndDate:   p.EndDate,
	}

	if err = period.IsValid(func(startdate, enddate *time.Time) {
		bldr.SetStartDate(*startdate)
		bldr.SetNillableEndDate(enddate)
	}); err != nil {
		vldErrs = append(vldErrs, err)
	}

	// Set and Validate MissionType if provided
	if vldErr := setNillableStringField(p.MissionType, func(v string) error {
		err := validateParams(p, "MissionType")
		if err != nil {
			return err
		}
		bldr.SetMissionType(mission.MissionType(v))
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

func (r *MissionRepository) DeleteByUUID(id uuid.UUID) error {
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
