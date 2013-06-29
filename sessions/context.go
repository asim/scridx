package sessions

import (
  "net/http"
)

type Context struct {
  W http.ResponseWriter
  R *http.Request
  Session *Session
  UserSession *UserSession
  Data *map[string]interface{}
}

// CTX replaces ResponseWriter and Request
func NewContext(w http.ResponseWriter, r *http.Request) (*Context, error) {
  var s *Session
  var u *UserSession
  var err error

  s, err = GetSession(w, r)
  if err != nil {
    return nil, err
  }

  u, err = GetUserSession(w, r)
  if err != nil {
    return nil, err
  }

  data := setup(w, r, u, s)

  return &Context{
    W: w,
    R: r,
    Session: s,
    UserSession: u,
    Data: data,
  }, nil
}

// Used once when a new request context is initialized
func setup(w http.ResponseWriter, r *http.Request, u *UserSession, s *Session) *map[string]interface{} {
  info := s.GetFlash(w, r)
  csrf := s.GetOrSetCSRF(w, r)
  user := u.GetAuthToken()

  data := map[string]interface{}{
    "info": info,
    "csrf": csrf,
    "_user": user,
  }

  return &data
}
