package repository

import (
	"errors"
	"time"

	"github.com/theterminalguy/tnet/util/date"
)

var (
	ErrStartDateEqualEndDate      = errors.New("error: start date equal end date")
	ErrStartDateLesserThanEndDate = errors.New("error: start date should not be lesser than end date")
)

type DatePeriod struct {
	StartDate string `json:"start_date" validate:"datetime=2006-01-02T15:04:05Z07:00"`
	EndDate   string `json:"end_date" validate:"datetime=2006-01-02T15:04:05Z07:00"`
}

func (p *DatePeriod) IsValid(cb func(startdate, enddate *time.Time)) error {

	var sd, ed *time.Time
	if p.StartDate != "" {
		err := ValidateParams(p, "StartDate")
		if err != nil {
			return err
		}
		sd, err = date.JSStringToRFC3339(p.StartDate)
		if err != nil {
			return err
		}
	}

	if p.EndDate != "" {
		err := ValidateParams(p, "EndDate")
		if err != nil {
			return err
		}
		ed, err = date.JSStringToRFC3339(p.EndDate)
		if err != nil {
			return err
		}
	}

	if sd != nil && ed != nil {
		return IsEqual(*sd, *ed)
	}

	cb(sd, ed)
	return nil

}

func (p *DatePeriod) IsValidStrict(cb func(startdate, enddate time.Time)) error {

	var sd, ed *time.Time

	err := ValidateParams(p, "StartDate")
	if err != nil {
		return err
	}
	sd, err = date.JSStringToRFC3339(p.StartDate)
	if err != nil {
		return err
	}

	err = ValidateParams(p, "EndDate")
	if err != nil {
		return err
	}
	ed, err = date.JSStringToRFC3339(p.EndDate)
	if err != nil {
		return err
	}

	if sd != nil && ed != nil {
		return IsEqual(*sd, *ed)
	}

	cb(*sd, *ed)
	return nil

}

// IsEqual checks the given times if they are equal or if the start time is after the end time
func IsEqual(sd, ed time.Time) error {

	if sd.Truncate(24 * time.Hour).Equal(ed.Truncate(24 * time.Hour)) {
		return ErrStartDateEqualEndDate
	}

	if sd.Truncate(24 * time.Hour).After(ed.Truncate(24 * time.Hour)) {
		return ErrStartDateLesserThanEndDate
	}

	return nil
}
