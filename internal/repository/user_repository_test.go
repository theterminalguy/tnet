package repository_test

import (
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/theterminalguy/tnet/ent/schema/userrole"
	repo "github.com/theterminalguy/tnet/internal/repository"
)

func init() {
	os.Setenv("ENV", "test")
	os.Setenv("PLATFORM", "docker")
	os.Setenv("PORT", "8080")
	os.Setenv("APP_HOST", "http://localhost")
	os.Setenv("JWT_SIGNED_SECRET", "test-secret")
}

func newUserRepo() *repo.UserRepository {
	return repo.NewUserRepository()
}

func validUserParams() repo.UserParams {
	return repo.UserParams{
		FirstName: "Jane",
		LastName:  "Doe",
		Email:     uuid.New().String() + "@example.com",
		PhotoURL:  "https://example.com/photo.jpg",
		Role:      userrole.Talent,
	}
}

// ---- GetAll ----

func TestGetAllUsers_ReturnsEmptySlice(t *testing.T) {
	r := newUserRepo()
	records, err := r.GetAll()
	require.NoError(t, err)
	assert.NotNil(t, records)
}

func TestGetAllUsers_ReturnsCreatedUser(t *testing.T) {
	r := newUserRepo()
	p := validUserParams()
	_, err := r.Create(p)
	require.NoError(t, err)

	records, err := r.GetAll()
	require.NoError(t, err)
	assert.True(t, len(records) >= 1)
}

func TestGetAllUsers_ExcludesDeletedUsers(t *testing.T) {
	r := newUserRepo()
	p := validUserParams()
	created, err := r.Create(p)
	require.NoError(t, err)

	err = r.DeleteByID(created.ID)
	require.NoError(t, err)

	records, err := r.GetAll()
	require.NoError(t, err)
	for _, u := range records {
		assert.NotEqual(t, created.ID, u.ID)
	}
}

// ---- GetByID ----

func TestGetUserByID_HappyPath(t *testing.T) {
	r := newUserRepo()
	p := validUserParams()
	created, err := r.Create(p)
	require.NoError(t, err)

	found, err := r.GetByID(created.ID)
	require.NoError(t, err)
	assert.Equal(t, created.ID, found.ID)
}

func TestGetUserByID_NotFound(t *testing.T) {
	r := newUserRepo()
	_, err := r.GetByID(uuid.New())
	assert.Error(t, err)
}

func TestGetUserByID_DeletedReturnsError(t *testing.T) {
	r := newUserRepo()
	p := validUserParams()
	created, err := r.Create(p)
	require.NoError(t, err)

	err = r.DeleteByID(created.ID)
	require.NoError(t, err)

	_, err = r.GetByID(created.ID)
	assert.ErrorIs(t, err, repo.ErrRecordDeleted)
}

// ---- GetByEmail ----

func TestGetUserByEmail_HappyPath(t *testing.T) {
	r := newUserRepo()
	p := validUserParams()
	created, err := r.Create(p)
	require.NoError(t, err)

	found, err := r.GetByEmail(created.Email)
	require.NoError(t, err)
	assert.Equal(t, created.Email, found.Email)
}

func TestGetUserByEmail_NotFound(t *testing.T) {
	r := newUserRepo()
	_, err := r.GetByEmail("nonexistent@example.com")
	assert.Error(t, err)
}

func TestGetUserByEmail_DeletedReturnsError(t *testing.T) {
	r := newUserRepo()
	p := validUserParams()
	created, err := r.Create(p)
	require.NoError(t, err)

	err = r.DeleteByID(created.ID)
	require.NoError(t, err)

	_, err = r.GetByEmail(created.Email)
	assert.ErrorIs(t, err, repo.ErrRecordDeleted)
}

// ---- Create ----

func TestCreateUser_HappyPath(t *testing.T) {
	r := newUserRepo()
	p := validUserParams()
	created, err := r.Create(p)
	require.NoError(t, err)
	assert.Equal(t, p.FirstName, created.FirstName)
	assert.Equal(t, p.Email, created.Email)
}

func TestCreateUser_MissingRequiredFields(t *testing.T) {
	r := newUserRepo()
	_, err := r.Create(repo.UserParams{})
	assert.Error(t, err)
}

func TestCreateUser_DuplicateEmailReturnsError(t *testing.T) {
	r := newUserRepo()
	p := validUserParams()
	_, err := r.Create(p)
	require.NoError(t, err)

	_, err = r.Create(p)
	assert.Error(t, err)
}

// ---- DeleteByID ----

func TestDeleteUserByID_HappyPath(t *testing.T) {
	r := newUserRepo()
	p := validUserParams()
	created, err := r.Create(p)
	require.NoError(t, err)

	err = r.DeleteByID(created.ID)
	assert.NoError(t, err)
}

func TestDeleteUserByID_NotFound(t *testing.T) {
	r := newUserRepo()
	err := r.DeleteByID(uuid.New())
	assert.Error(t, err)
}

func TestDeleteUserByID_AlreadyDeletedReturnsError(t *testing.T) {
	r := newUserRepo()
	p := validUserParams()
	created, err := r.Create(p)
	require.NoError(t, err)

	err = r.DeleteByID(created.ID)
	require.NoError(t, err)

	err = r.DeleteByID(created.ID)
	assert.Error(t, err)
}
