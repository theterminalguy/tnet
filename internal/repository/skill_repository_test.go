package repository_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/theterminalguy/tnet/ent/schema/userrole"
	"github.com/theterminalguy/tnet/ent/talent"
	repo "github.com/theterminalguy/tnet/internal/repository"
)

func createTestTalent(t *testing.T) (uuid.UUID, uuid.UUID) {
	t.Helper()
	ur := repo.NewUserRepository()
	u, err := ur.Create(repo.UserParams{
		FirstName: "Test",
		LastName:  "User",
		Email:     uuid.New().String() + "@test.com",
		PhotoURL:  "https://example.com/photo.jpg",
		Role:      userrole.Talent,
	})
	require.NoError(t, err)

	tr := repo.NewTalentRepository()
	tal, err := tr.Create(repo.TalentParams{
		UserID:                u.ID,
		FirstName:             "Test",
		LastName:              "Talent",
		Email:                 u.Email,
		PreferredName:         "Testy",
		Pronoun:               "they/them",
		PreferredJobTitle:     "Engineer",
		ProfessionalStartDate: "2020-01-01",
		Phone:                 "+1234567890",
		CountryCode:           "US",
		City:                  "New York",
		JobPreference:         talent.JobPreferenceRemote,
		Available:             true,
		TimeZone:              "GMT",
		State:                 "NY",
	})
	require.NoError(t, err)
	return u.ID, tal.ID
}

func validSkillParams(talentID uuid.UUID) repo.SkillParams {
	return repo.SkillParams{
		TalentID:          talentID,
		Name:              "Go",
		YearsOfExperience: 2.0,
		Preferred:         true,
		Note:              "Built microservices with Go",
	}
}

// ---- GetAll ----

func TestGetAllSkills_ReturnsNoError(t *testing.T) {
	r := repo.NewSkillRepository()
	records, err := r.GetAll()
	require.NoError(t, err)
	assert.NotNil(t, records)
}

func TestGetAllSkills_ReturnsCreatedSkill(t *testing.T) {
	_, talentID := createTestTalent(t)
	r := repo.NewSkillRepository()
	_, err := r.Create(validSkillParams(talentID))
	require.NoError(t, err)

	records, err := r.GetAll()
	require.NoError(t, err)
	assert.True(t, len(records) >= 1)
}

func TestGetAllSkills_ExcludesDeletedSkills(t *testing.T) {
	_, talentID := createTestTalent(t)
	r := repo.NewSkillRepository()
	created, err := r.Create(validSkillParams(talentID))
	require.NoError(t, err)

	err = r.DeleteByID(created.ID)
	require.NoError(t, err)

	records, err := r.GetAll()
	require.NoError(t, err)
	for _, s := range records {
		assert.NotEqual(t, created.ID, s.ID)
	}
}

// ---- GetByID ----

func TestGetSkillByID_HappyPath(t *testing.T) {
	_, talentID := createTestTalent(t)
	r := repo.NewSkillRepository()
	created, err := r.Create(validSkillParams(talentID))
	require.NoError(t, err)

	found, err := r.GetByID(created.ID)
	require.NoError(t, err)
	assert.Equal(t, created.ID, found.ID)
}

func TestGetSkillByID_NotFound(t *testing.T) {
	r := repo.NewSkillRepository()
	_, err := r.GetByID(uuid.New())
	assert.Error(t, err)
}

func TestGetSkillByID_DeletedReturnsError(t *testing.T) {
	_, talentID := createTestTalent(t)
	r := repo.NewSkillRepository()
	created, err := r.Create(validSkillParams(talentID))
	require.NoError(t, err)

	err = r.DeleteByID(created.ID)
	require.NoError(t, err)

	_, err = r.GetByID(created.ID)
	assert.ErrorIs(t, err, repo.ErrRecordDeleted)
}

// ---- Create ----

func TestCreateSkill_HappyPath(t *testing.T) {
	_, talentID := createTestTalent(t)
	r := repo.NewSkillRepository()
	p := validSkillParams(talentID)
	created, err := r.Create(p)
	require.NoError(t, err)
	assert.Equal(t, "go", created.Name) // stored lowercase
}

func TestCreateSkill_MissingRequiredFields(t *testing.T) {
	r := repo.NewSkillRepository()
	_, err := r.Create(repo.SkillParams{})
	assert.Error(t, err)
}

func TestCreateSkill_InvalidTalentIDReturnsError(t *testing.T) {
	r := repo.NewSkillRepository()
	p := validSkillParams(uuid.New()) // non-existent talent
	_, err := r.Create(p)
	assert.Error(t, err)
}

// ---- DeleteByID ----

func TestDeleteSkillByID_HappyPath(t *testing.T) {
	_, talentID := createTestTalent(t)
	r := repo.NewSkillRepository()
	created, err := r.Create(validSkillParams(talentID))
	require.NoError(t, err)

	err = r.DeleteByID(created.ID)
	assert.NoError(t, err)
}

func TestDeleteSkillByID_NotFound(t *testing.T) {
	r := repo.NewSkillRepository()
	err := r.DeleteByID(uuid.New())
	assert.Error(t, err)
}

func TestDeleteSkillByID_AlreadyDeletedReturnsError(t *testing.T) {
	_, talentID := createTestTalent(t)
	r := repo.NewSkillRepository()
	created, err := r.Create(validSkillParams(talentID))
	require.NoError(t, err)

	err = r.DeleteByID(created.ID)
	require.NoError(t, err)

	err = r.DeleteByID(created.ID)
	assert.Error(t, err)
}
