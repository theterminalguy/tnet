package entsm

import (
	"encoding/base32"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	repo "github.com/theterminalguy/tnet/internal/repository"
)

var (

	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key = []byte(os.Getenv("SESSION_KEY"))
	// store = sessions.NewCookieStore(key)
	// store = sessions.NewFilesystemStore(".", key)
	store = NewEntSystemStore(key)
)

const DefaultSessionName = "__tentn_user_session"

func GetSessionStore() *EntSystemStore {
	return store
}

type EntSystemStore struct {
	Codecs      []securecookie.Codec
	Options     *sessions.Options // default configuration
	SessionRepo *repo.SessionRepository
}

func NewEntSystemStore(keyPairs ...[]byte) *EntSystemStore {
	cs := &EntSystemStore{
		Codecs: securecookie.CodecsFromPairs(keyPairs...),
		Options: &sessions.Options{
			Path:   "/",
			MaxAge: 86400 * 30, // 30 -days
		},
		SessionRepo: repo.NewSessionRepository(),
	}
	cs.MaxAge(cs.Options.MaxAge)
	return cs
}

func (s *EntSystemStore) MaxAge(age int) {
	s.Options.MaxAge = age
	// Set the maxAge for each securecookie instance.
	for _, codec := range s.Codecs {
		if sc, ok := codec.(*securecookie.SecureCookie); ok {
			sc.MaxAge(age)
		}
	}
}

func (s *EntSystemStore) Get(r *http.Request, name string) (*sessions.Session, error) {
	return sessions.GetRegistry(r).Get(s, name)
}

func (s *EntSystemStore) New(r *http.Request, name string) (*sessions.Session, error) {
	session := sessions.NewSession(s, name)
	opts := *s.Options
	session.Options = &opts
	session.IsNew = true
	var err error
	if c, errCookie := r.Cookie(name); errCookie == nil {
		err = securecookie.DecodeMulti(name, c.Value, &session.ID, s.Codecs...)
		if err == nil {
			err = s.load(session)
			if err == nil {
				session.IsNew = false
			}
		}
	}
	return session, err
}

func (s *EntSystemStore) Save(r *http.Request, w http.ResponseWriter,
	session *sessions.Session) error {
	// Delete if max-age is <= 0
	if session.Options.MaxAge <= 0 {
		if err := s.erase(session); err != nil {
			return err
		}
		http.SetCookie(w, sessions.NewCookie(session.Name(), "", session.Options))
		return nil
	}
	if session.ID == "" {
		session.ID = strings.TrimRight(
			base32.StdEncoding.EncodeToString(
				securecookie.GenerateRandomKey(32)), "=")
	}
	if err := s.save(session); err != nil {
		return err
	}
	encoded, err := securecookie.EncodeMulti(session.Name(), session.ID,
		s.Codecs...)
	if err != nil {
		return err
	}
	http.SetCookie(w, sessions.NewCookie(session.Name(), encoded, session.Options))
	return nil
}

// save writes encoded session data to db.
func (s *EntSystemStore) save(session *sessions.Session) error {
	encoded, err := securecookie.EncodeMulti(session.Name(), session.Values,
		s.Codecs...)
	if err != nil {
		return err
	}
	teamID, _ := session.Values["team_id"].(string)
	params := repo.SessionRepositoryParams{
		SessionID: session.ID,
		Encoded:   encoded,
		TeamID:    teamID,
	}
	_, err = s.SessionRepo.CreateSession(params)
	if err != nil {
		return err
	}
	return nil
}

// load reads a session record from db and decodes its content into session.Values.
func (s *EntSystemStore) load(session *sessions.Session) error {
	record, err := s.SessionRepo.GetBySessionID(session.ID)
	if err != nil {
		return err
	}
	if err = securecookie.DecodeMulti(session.Name(), record.Encoded,
		&session.Values, s.Codecs...); err != nil {
		return err
	}
	return nil
}

// delete session record from db
func (s *EntSystemStore) erase(session *sessions.Session) error {
	err := s.SessionRepo.DeleteSession(session.ID)
	if err != nil {
		return err
	}
	return nil
}
