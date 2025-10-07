package scope

import (
	"github.com/theterminalguy/tentn/ent"
	"github.com/theterminalguy/tentn/ent/education"
	"github.com/theterminalguy/tentn/ent/jobapplication"
	"github.com/theterminalguy/tentn/ent/portfoliolink"
	"github.com/theterminalguy/tentn/ent/skill"
	"github.com/theterminalguy/tentn/ent/workexperience"
	repo "github.com/theterminalguy/tentn/internal/repository"
	"github.com/google/uuid"
)

type TalentScope struct {
	Talent *ent.Talent
}

func NewTalentScope(talent *ent.Talent) *TalentScope {
	return &TalentScope{
		Talent: talent,
	}
}

func (t *TalentScope) GetID() uuid.UUID {
	return t.Talent.ID
}

func (t *TalentScope) GetWorkExperiences() ([]*ent.WorkExperience, error) {
	return t.Talent.QueryWorkExperiences().All(repo.GetDBContext())
}

func (t *TalentScope) GetWorkExperienceByID(id uuid.UUID) (*ent.WorkExperience, error) {
	return t.Talent.QueryWorkExperiences().
		Where(workexperience.ID(id)).
		First(repo.GetDBContext())
}

func (t *TalentScope) UpdateWorkExperience(uuid uuid.UUID, params repo.WorkExperienceParams) (*ent.WorkExperience, []error) {
	params.TalentID = t.Talent.ID
	record, err := t.GetWorkExperienceByID(uuid)
	if err != nil {
		return nil, []error{err}
	}
	return repo.NewWorkExperienceRepository().Update(record.ID, params)
}

func (t *TalentScope) DeleteWorkExperience(uuid uuid.UUID) error {
	record, err := t.GetWorkExperienceByID(uuid)
	if err != nil {
		return err
	}
	return repo.NewWorkExperienceRepository().DeleteByID(record.ID)
}

func (t *TalentScope) GetSkills() ([]*ent.Skill, error) {
	return t.Talent.QuerySkills().All(repo.GetDBContext())
}

func (t *TalentScope) GetSkillByID(id uuid.UUID) (*ent.Skill, error) {
	return t.Talent.QuerySkills().
		Where(skill.ID(id)).
		First(repo.GetDBContext())
}

func (t *TalentScope) UpdateSkill(uuid uuid.UUID, params repo.SkillParams) (*ent.Skill, []error) {
	params.TalentID = t.Talent.ID
	record, err := t.GetSkillByID(uuid)
	if err != nil {
		return nil, []error{err}
	}
	return repo.NewSkillRepository().Update(record.ID, params)
}

func (t *TalentScope) DeleteSkill(id uuid.UUID) error {
	record, err := t.GetSkillByID(id)
	if err != nil {
		return err
	}
	return repo.NewSkillRepository().DeleteByID(record.ID)
}

func (t *TalentScope) GetPortfolioLinks() ([]*ent.PortfolioLink, error) {
	return t.Talent.QueryPortfoliolinks().All(repo.GetDBContext())
}

func (t *TalentScope) GetPortfolioLinkByID(id uuid.UUID) (*ent.PortfolioLink, error) {
	return t.Talent.QueryPortfoliolinks().
		Where(portfoliolink.ID(id)).
		First(repo.GetDBContext())
}

func (t *TalentScope) UpdatePortfolioLink(id uuid.UUID, params repo.PortfolioLinkParams) (*ent.PortfolioLink, []error) {
	params.TalentID = t.Talent.ID
	record, err := t.GetPortfolioLinkByID(id)
	if err != nil {
		return nil, []error{err}
	}
	return repo.NewPortfolioLinkRepository().Update(record.ID, params)
}

func (t *TalentScope) DeletePortfolioLink(id uuid.UUID) error {
	record, err := t.GetPortfolioLinkByID(id)
	if err != nil {
		return err
	}
	return repo.NewPortfolioLinkRepository().DeleteByID(record.ID)
}

func (t *TalentScope) GetJobApplications() ([]*ent.JobApplication, error) {
	return t.Talent.QueryJobApplications().All(repo.GetDBContext())
}

func (t *TalentScope) GetJobApplicationByID(id uuid.UUID) (*ent.JobApplication, error) {
	return t.Talent.QueryJobApplications().
		Where(jobapplication.ID(id)).
		First(repo.GetDBContext())
}

func (t *TalentScope) GetEducations() ([]*ent.Education, error) {
	return t.Talent.QueryEducations().All(repo.GetDBContext())
}

func (t *TalentScope) GetEducationByID(id uuid.UUID) (*ent.Education, error) {
	return t.Talent.QueryEducations().
		Where(education.ID(id)).
		First(repo.GetDBContext())
}

func (t *TalentScope) UpdateEducation(id uuid.UUID, params repo.EducationParams) (*ent.Education, []error) {
	params.TalentID = t.Talent.ID
	record, err := t.GetEducationByID(id)
	if err != nil {
		return nil, []error{err}
	}
	return repo.NewEducationRepository().Update(record.ID, params)
}

func (t *TalentScope) DeleteEducation(id uuid.UUID) error {
	record, err := t.GetJobApplicationByID(id)
	if err != nil {
		return err
	}
	return repo.NewEducationRepository().DeleteByID(record.ID)
}
