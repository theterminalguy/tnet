package scope

import (
	"github.com/10hourlabs/tentn/ent"
	"github.com/10hourlabs/tentn/ent/education"
	"github.com/10hourlabs/tentn/ent/jobapplication"
	"github.com/10hourlabs/tentn/ent/portfoliolink"
	"github.com/10hourlabs/tentn/ent/skill"
	"github.com/10hourlabs/tentn/ent/workexperience"
	repo "github.com/10hourlabs/tentn/internal/repository"
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

func (t *TalentScope) GetWorkExperiences() ([]*ent.WorkExperience, error) {
	return t.Talent.QueryWorkExperiences().All(repo.GetDBContext())
}

func (t *TalentScope) GetWorkExperienceByUUID(uuid uuid.UUID) (*ent.WorkExperience, error) {
	return t.Talent.QueryWorkExperiences().
		Where(workexperience.UUIDEQ(uuid)).
		First(repo.GetDBContext())
}

func (t *TalentScope) UpdateWorkExperience(uuid uuid.UUID, params repo.WorkExperienceParams) (*ent.WorkExperience, []error) {
	params.TalentUUID = t.Talent.UUID
	record, err := t.GetWorkExperienceByUUID(uuid)
	if err != nil {
		return nil, []error{err}
	}
	return repo.NewWorkExperienceRepository().Update(record.UUID, params)
}

func (t *TalentScope) DeleteWorkExperience(uuid uuid.UUID) error {
	record, err := t.GetWorkExperienceByUUID(uuid)
	if err != nil {
		return err
	}
	return repo.NewWorkExperienceRepository().DeleteByUUID(record.UUID)
}

func (t *TalentScope) GetSkills() ([]*ent.Skill, error) {
	return t.Talent.QuerySkills().All(repo.GetDBContext())
}

func (t *TalentScope) GetSkillByUUID(uuid uuid.UUID) (*ent.Skill, error) {
	return t.Talent.QuerySkills().
		Where(skill.UUIDEQ(uuid)).
		First(repo.GetDBContext())
}

func (t *TalentScope) UpdateSkill(uuid uuid.UUID, params repo.SkillParams) (*ent.Skill, []error) {
	params.TalentUUID = t.Talent.UUID
	record, err := t.GetSkillByUUID(uuid)
	if err != nil {
		return nil, []error{err}
	}
	return repo.NewSkillRepository().Update(record.UUID, params)
}

func (t *TalentScope) DeleteSkill(uuid uuid.UUID) error {
	record, err := t.GetSkillByUUID(uuid)
	if err != nil {
		return err
	}
	return repo.NewSkillRepository().DeleteByUUID(record.UUID)
}

func (t *TalentScope) GetPortfolioLinks() ([]*ent.PortfolioLink, error) {
	return t.Talent.QueryPortfoliolinks().All(repo.GetDBContext())
}

func (t *TalentScope) GetPortfolioLinkByUUID(uuid uuid.UUID) (*ent.PortfolioLink, error) {
	return t.Talent.QueryPortfoliolinks().
		Where(portfoliolink.UUIDEQ(uuid)).
		First(repo.GetDBContext())
}

func (t *TalentScope) UpdatePortfolioLink(uuid uuid.UUID, params repo.PortfolioLinkParams) (*ent.PortfolioLink, []error) {
	params.TalentUUID = t.Talent.UUID
	record, err := t.GetPortfolioLinkByUUID(uuid)
	if err != nil {
		return nil, []error{err}
	}
	return repo.NewPortfolioLinkRepository().Update(record.UUID, params)
}

func (t *TalentScope) DeletePortfolioLink(uuid uuid.UUID) error {
	record, err := t.GetPortfolioLinkByUUID(uuid)
	if err != nil {
		return err
	}
	return repo.NewPortfolioLinkRepository().DeleteByUUID(record.UUID)
}

func (t *TalentScope) GetJobApplications() ([]*ent.JobApplication, error) {
	return t.Talent.QueryJobApplications().All(repo.GetDBContext())
}

func (t *TalentScope) GetJobApplicationByUUID(uuid uuid.UUID) (*ent.JobApplication, error) {
	return t.Talent.QueryJobApplications().
		Where(jobapplication.UUIDEQ(uuid)).
		First(repo.GetDBContext())
}

func (t *TalentScope) GetEducations() ([]*ent.Education, error) {
	return t.Talent.QueryEducations().All(repo.GetDBContext())
}

func (t *TalentScope) GetEducationByUUID(uuid uuid.UUID) (*ent.Education, error) {
	return t.Talent.QueryEducations().
		Where(education.UUIDEQ(uuid)).
		First(repo.GetDBContext())
}

func (t *TalentScope) UpdateEducation(uuid uuid.UUID, params repo.EducationParams) (*ent.Education, []error) {
	params.TalentUUID = t.Talent.UUID
	record, err := t.GetEducationByUUID(uuid)
	if err != nil {
		return nil, []error{err}
	}
	return repo.NewEducationRepository().Update(record.UUID, params)
}

func (t *TalentScope) DeleteEducation(uuid uuid.UUID) error {
	record, err := t.GetJobApplicationByUUID(uuid)
	if err != nil {
		return err
	}
	return repo.NewEducationRepository().DeleteByUUID(record.UUID)
}
