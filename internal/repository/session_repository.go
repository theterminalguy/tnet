package repository

import (
	"time"

	"github.com/10hourlabs/tentn/ent"
	"github.com/10hourlabs/tentn/ent/session"
)

type SessionRepository struct{}

type SessionRepositoryParams struct {
	SessionID string `json:"session_id"`
	Encoded   string `json:"encoded"`
	TeamID    string `json:"team_id"`
}

func NewSessionRepository() *SessionRepository {
	return &SessionRepository{}
}

func (s *SessionRepository) GetSessionByTeamID(teamID string) (*ent.Session, error) {
	record, err := dBConn.Session.Query().
		Where(session.TeamID(teamID)).
		Only(dBContext)
	if err != nil {
		return nil, err
	}
	return record, nil
}

func (s *SessionRepository) GetBySessionID(id string) (*ent.Session, error) {
	record, err := dBConn.Session.Query().
		Where(session.SessionIDEQ(id)).
		Only(dBContext)
	if err != nil {
		return nil, err
	}
	if record.DeletedAt != nil {
		return nil, ErrRecordDeleted
	}
	return record, nil
}

func (s *SessionRepository) CreateSession(p SessionRepositoryParams) (*ent.Session, error) {
	record, err := dBConn.Session.Create().
		SetSessionID(p.SessionID).
		SetEncoded(p.Encoded).
		SetTeamID(p.TeamID).
		Save(dBContext)

	if err != nil {
		return nil, err
	}
	return record, nil
}

func (s *SessionRepository) DeleteSession(id string) error {
	record, err := s.GetBySessionID(id)
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
