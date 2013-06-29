package sessions

import (
  "fmt"
  "log"
  "time"
  "encoding/json"
  "encoding/base64"
  "net/http"
  "crypto/rand"
  "crypto/aes"
  "crypto/sha1"
  mrand "math/rand"
  "github.com/gorilla/sessions"
  "github.com/gorilla/securecookie"
  "github.com/garyburd/redigo/redis"
)

func encode(b []byte) string {
  return base64.StdEncoding.EncodeToString(b)[:64]
}

const (
  SSID = "_s"
  AUTH = "_a"
  EKEY = "mymilkshakebringsalltheboystothe"
  AKEY = "ipitythefoolwhomesswithmeonthisd"
  MAXAGE = 86400 * 14
)

var (
  r = setRedisConn()
  seed = setSeed()
  encryptKey = setEncryptionKey()
  authKey = setAuthKey()
  sessionStore = newSessionStore(0)
  authStore = newSessionStore(0)
)

type Session struct {
  *sessions.Session
}

type UserSession struct {
  *sessions.Session
}

type Token struct {
  Hash string
  Time int64
  Id int64
  Username string
}

func setSeed() bool {
  mrand.Seed(time.Now().UTC().UnixNano())
  return true
}

func setRedisConn() redis.Conn {
  c, err := redis.Dial("tcp", ":6379")
  if err != nil {
    log.Fatal("Error connecting to redis", err)
    return nil
  }

  return c
}

func generateEncryptionKey() []byte {
  return securecookie.GenerateRandomKey(32)
}

func generateCSRFKey() []byte {
  b := make([]byte, 64)
  n, err := rand.Read(b)
  if n != 64 || err != nil {
    log.Print("Error generating CSRF key", err)
    return generateCSRFKey()
  }

  return b
}

func generateAuthKey() []byte {
  b := make([]byte, 64)
  n, err := rand.Read(b)
  if n != 64 || err != nil {
    log.Fatal("Error generating auth key", err)
    return []byte(nil)
  }

  return b
}

// encryption key does not exist, create one
func newEncryptionKey() []byte {
  skey := generateEncryptionKey()
  dkey := make([]byte, len(skey))

  c, err := aes.NewCipher([]byte(EKEY));
  if err != nil {
    log.Fatal("Error creating new encryption key", err)
    return []byte(nil)
  }

  for i := 0; i < 2; i++ {
    start := i * c.BlockSize()
    end := start + c.BlockSize()
    c.Encrypt(dkey[start:end], skey[start:end])
  }

  r.Do("SET", "0k", dkey[:16])
  r.Do("SET", "1k", dkey[16:])

  return skey
}

// auth key does not exist, create one
func newAuthKey() []byte {
  skey := generateAuthKey()
  dkey := make([]byte, len(skey))

  c, err := aes.NewCipher([]byte(AKEY));
  if err != nil {
    log.Fatal("Error creating new auth key", err)
    return []byte(nil)
  }

  for i := 0; i < 4; i++ {
    start := i * c.BlockSize()
    end := start + c.BlockSize()
    c.Encrypt(dkey[start:end], skey[start:end])
  }

  r.Do("SET", "0a", dkey[:32])
  r.Do("SET", "1a", dkey[32:])
  return skey
}

// setup the encryption key
func setEncryptionKey() []byte {
  var k0, k1 interface{}
  var err error

  if k0, err = r.Do("GET", "0k"); err != nil || k0 == nil {
    log.Println("error getting ekey, random key returned", err)
    return newEncryptionKey()
  }

  if k1, err = r.Do("GET", "1k"); err != nil || k1 == nil {
    log.Println("error getting ekey, random key returned", err)
    return newEncryptionKey()
  }

  c, err := aes.NewCipher([]byte(EKEY));
  if err != nil {
    log.Fatal("Error getting new aes cipher", err)
  }

  encKey := append(k0.([]byte), k1.([]byte)...)
  key := make([]byte, len(encKey))

  for i := 0; i < 2; i++ {
    start := i * c.BlockSize()
    end := start + c.BlockSize()
    c.Decrypt(key[start:end], encKey[start:end])
  }

  return key
}

// setup the auth key
func setAuthKey() []byte {
  var a0, a1 interface{}
  var err error

  if a0, err = r.Do("GET", "0a"); err != nil || a0 == nil {
    log.Println("error getting akey, random key returned", err)
    return newAuthKey()
  }

  if a1, err = r.Do("GET", "1a"); err != nil || a1 == nil {
    log.Println("error getting akey, random key returned", err)
    return newAuthKey()
  }

  c, err := aes.NewCipher([]byte(AKEY));
  if err != nil {
    log.Fatal("Error getting new aes cipher", err)
  }

  auKey := append(a0.([]byte), a1.([]byte)...)
  key := make([]byte, len(auKey))

  for i := 0; i < 4; i++ {
    start := i * c.BlockSize()
    end := start + c.BlockSize()
    c.Decrypt(key[start:end], auKey[start:end])
  }

  return key
}

func sessionOptions(maxAge int) *sessions.Options {
  return &sessions.Options{
    Path: "/",
    MaxAge: maxAge,
    HttpOnly: true,
//    Secure: true,
  }
}

func newSessionStore(maxAge int) *sessions.CookieStore {
  s := sessions.NewCookieStore(authKey, encryptKey)
  s.Options = sessionOptions(maxAge)
  return s
}

func GetSession(w http.ResponseWriter, r *http.Request) (*Session, error) {
  session, err := sessionStore.Get(r, SSID)
  return &Session{session}, err
}

func GetUserSession(w http.ResponseWriter, r *http.Request) (*UserSession, error) {
  session, err := authStore.Get(r, AUTH)
  return &UserSession{session}, err
}

// User Auth Session Methods

func (s *UserSession) SetAuthToken(ctx *Context, username string, id int64, rememberme bool) error {
  if rememberme {
    s.Options = sessionOptions(MAXAGE)
  }

  t := time.Now()
  hash := sha1.New()
  hash.Write([]byte(t.String()))

  // create token
  token := &Token{
    Hash: fmt.Sprintf("%x", hash.Sum(nil))[0:20],
    Time: t.Unix(),
    Id: id,
    Username: username,
  }

  // marshal the data
  jtoken, err := json.Marshal(token)
  if err != nil {
    log.Println(err)
    return err
  }

  s.Values["."] = jtoken

  // save session
  if err := s.Save(ctx.R, ctx.W); err != nil {
    log.Println(err)
    return err
  }

  return nil
}

func (s *UserSession) GetAuthToken() *Token {
  jtoken := s.Values["."]
  if jtoken == nil {
    return nil
  }

  token := new(Token)
  err := json.Unmarshal(jtoken.([]byte), &token)
  if err != nil {
    log.Println(err)
    return nil
  }

  if !validAuthToken(token) {
    return nil
  }

  return token
}

func (s *UserSession) DeleteAuthToken(ctx *Context) error {
  s.Options = sessionOptions(-1)
  s.Values["."] = nil

  if err := s.Save(ctx.R, ctx.W); err != nil {
    log.Println(err)
    return err
  }

  return nil
}

func validAuthToken(token *Token) bool {
  if len(token.Hash) != 20 {
    return false
  }

  if token.Id == 0 {
    return false
  }

  if len(token.Username) == 0 {
    return false
  }

  if (time.Now().Unix() - token.Time) > MAXAGE {
    // auth token is older than max age.
    return false
  }

  return true
}

func (s *UserSession) LoggedIn() bool {
  if token := s.GetAuthToken(); token != nil {
    return true
  }

  return false
}


// Secure cookie management for sessions.
func (s *Session) SetValue(w http.ResponseWriter, r *http.Request, key string, value interface{}) error {
  s.Values[key] = value
  if err := s.Save(r, w); err != nil {
    log.Println(err)
    return err
  }

  return nil
}

func (s *Session) GetValue(key string) interface{}  {
  return s.Values[key]
}

// Flash session management
func (s *Session) SetFlash(w http.ResponseWriter, r *http.Request, msg string) {
  s.AddFlash(msg, "info")
  s.Save(r, w)
}

func (s *Session) GetFlash(w http.ResponseWriter, r *http.Request) []interface{} {
  if flashes := s.Flashes("info"); len(flashes) > 0 {
    s.Save(r, w)
    return flashes
  }

  return nil
}

func (s *Session) SetForm(w http.ResponseWriter, r *http.Request, name string, data interface{}) {
  v, _ := json.Marshal(data)
  _ = s.SetValue(w, r, name + "Form", v)
}

func (s *Session) GetForm(name string) map[string][]string {
  v := s.GetValue(name + "Form")
  var form map[string][]string
  if v != nil {
    _ = json.Unmarshal(v.([]byte), &form)
  }
  return form
}

func (s *Session) ClearForm(w http.ResponseWriter, r *http.Request, name string) {
  _ = s.SetValue(w, r, name + "Form", nil)
}

func (s *Session) SetCSRF(w http.ResponseWriter, r *http.Request) string {
  csrf := encode(generateCSRFKey())
  s.SetValue(w, r, "csrf", csrf)
  return csrf
}

func (s *Session) GetOrSetCSRF(w http.ResponseWriter, r *http.Request) string {
  csrf := s.GetValue("csrf")

  if csrf == nil || len(csrf.(string)) == 0 {
    csrf = s.SetCSRF(w, r)
  }

  return csrf.(string)
}

func (s *Session) IsValidCSRF(csrf string) bool {
  savedCsrf := s.GetValue("csrf")

  if savedCsrf != nil && csrf == savedCsrf.(string) {
    return true
  }

  return false
}
