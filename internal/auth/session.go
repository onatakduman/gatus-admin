package auth

import (
	"net/http"

	"github.com/gorilla/sessions"
)

const (
	SessionName   = "admin-session"
	AuthedKey     = "authed"
	FlashErrorKey = "flash-error"
	FlashOkKey    = "flash-ok"
)

type Store struct {
	cs *sessions.CookieStore
}

func NewStore(secret []byte, secure bool) *Store {
	cs := sessions.NewCookieStore(secret)
	cs.Options = &sessions.Options{
		Path:     "/admin",
		MaxAge:   60 * 60 * 24,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteStrictMode,
	}
	return &Store{cs: cs}
}

func (s *Store) Get(r *http.Request) *sessions.Session {
	sess, _ := s.cs.Get(r, SessionName)
	return sess
}

func (s *Store) IsAuthed(r *http.Request) bool {
	v, ok := s.Get(r).Values[AuthedKey].(bool)
	return ok && v
}

func (s *Store) SetAuthed(r *http.Request, w http.ResponseWriter, v bool) error {
	sess := s.Get(r)
	sess.Values[AuthedKey] = v
	return sess.Save(r, w)
}

func (s *Store) AddFlash(r *http.Request, w http.ResponseWriter, key, msg string) error {
	sess := s.Get(r)
	sess.AddFlash(msg, key)
	return sess.Save(r, w)
}

func (s *Store) PopFlashes(r *http.Request, w http.ResponseWriter, key string) []string {
	sess := s.Get(r)
	raw := sess.Flashes(key)
	_ = sess.Save(r, w)
	out := make([]string, 0, len(raw))
	for _, v := range raw {
		if sv, ok := v.(string); ok {
			out = append(out, sv)
		}
	}
	return out
}
