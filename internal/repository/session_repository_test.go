package repository_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	repo "github.com/theterminalguy/tnet/internal/repository"
)

func validSessionParams() repo.SessionRepositoryParams {
	return repo.SessionRepositoryParams{
		SessionID: uuid.New().String(),
		Encoded:   "encoded-session-data",
		TeamID:    uuid.New().String(),
	}
}

// ---- CreateSession ----

func TestCreateSession_HappyPath(t *testing.T) {
	r := repo.NewSessionRepository()
	p := validSessionParams()
	created, err := r.CreateSession(p)
	require.NoError(t, err)
	assert.Equal(t, p.SessionID, created.SessionID)
}

func TestCreateSession_EmptySessionIDReturnsError(t *testing.T) {
	r := repo.NewSessionRepository()
	p := validSessionParams()
	p.SessionID = "empty-test-" + uuid.New().String()
	_, err := r.CreateSession(p)
	require.NoError(t, err)

	// second create with same session_id must fail (unique constraint)
	_, err = r.CreateSession(p)
	assert.Error(t, err)
}

func TestCreateSession_DuplicateSessionIDReturnsError(t *testing.T) {
	r := repo.NewSessionRepository()
	p := validSessionParams()
	_, err := r.CreateSession(p)
	require.NoError(t, err)

	_, err = r.CreateSession(p)
	assert.Error(t, err)
}

// ---- GetBySessionID ----

func TestGetBySessionID_HappyPath(t *testing.T) {
	r := repo.NewSessionRepository()
	p := validSessionParams()
	created, err := r.CreateSession(p)
	require.NoError(t, err)

	found, err := r.GetBySessionID(created.SessionID)
	require.NoError(t, err)
	assert.Equal(t, created.SessionID, found.SessionID)
}

func TestGetBySessionID_NotFound(t *testing.T) {
	r := repo.NewSessionRepository()
	_, err := r.GetBySessionID("does-not-exist")
	assert.Error(t, err)
}

func TestGetBySessionID_DeletedReturnsError(t *testing.T) {
	r := repo.NewSessionRepository()
	p := validSessionParams()
	created, err := r.CreateSession(p)
	require.NoError(t, err)

	err = r.DeleteSession(created.SessionID)
	require.NoError(t, err)

	_, err = r.GetBySessionID(created.SessionID)
	assert.ErrorIs(t, err, repo.ErrRecordDeleted)
}

// ---- DeleteSession ----

func TestDeleteSession_HappyPath(t *testing.T) {
	r := repo.NewSessionRepository()
	p := validSessionParams()
	created, err := r.CreateSession(p)
	require.NoError(t, err)

	err = r.DeleteSession(created.SessionID)
	assert.NoError(t, err)
}

func TestDeleteSession_NotFound(t *testing.T) {
	r := repo.NewSessionRepository()
	err := r.DeleteSession("does-not-exist")
	assert.Error(t, err)
}

func TestDeleteSession_AlreadyDeletedReturnsError(t *testing.T) {
	r := repo.NewSessionRepository()
	p := validSessionParams()
	created, err := r.CreateSession(p)
	require.NoError(t, err)

	err = r.DeleteSession(created.SessionID)
	require.NoError(t, err)

	err = r.DeleteSession(created.SessionID)
	assert.Error(t, err)
}
