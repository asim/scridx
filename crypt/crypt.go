package crypt

import (
  "code.google.com/p/go.crypto/bcrypt"
  "code.google.com/p/go.crypto/scrypt"
  "crypto/rand"
  "crypto/subtle"
  "encoding/base64"
  "log"
  "io"
  "time"
)

const (
  BCRYPT = "b"
  SCRYPT = "s"
)

func encode(b []byte) string {
  return "eb" + base64.StdEncoding.EncodeToString(b)
}

func decode(s string) []byte {
  b, err := base64.StdEncoding.DecodeString(s[2:])
  if err != nil {
    log.Println(err)
    return []byte(nil)
  }

  return b
}

func genSalt() []byte {
  b := make([]byte, 10)
  n, err := io.ReadFull(rand.Reader, b)
  if n != len(b) || err != nil {
    log.Println(err)
    return []byte(nil)
  }

  return b
}

func sCrypt(salt []byte, password []byte) []byte {
  s, err := scrypt.Key(password, salt, 16384, 8, 1, 64)
  if err != nil {
    log.Println(err)
    return []byte(nil)
  }

  return s
}

func sCryptEqual(hashedPass, salt, password []byte) bool {
  enteredPass := sCrypt(salt, password)
  if subtle.ConstantTimeCompare(hashedPass, enteredPass) == 1 {
    return true
  }

  return false
}

func bCrypt(salt []byte, password []byte) []byte {
  password = append(salt, password...)
  b, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
  if err != nil {
    log.Println(err)
    return []byte(nil)
  }

  return b
}

func bCryptEqual(hashedPass, salt, password []byte) bool {
  password = append(salt, password...)
  err := bcrypt.CompareHashAndPassword(hashedPass, password)
  if err == nil {
    return true
  }

  log.Println(err)
  return false
}

func hashPassword(password string) ([]byte, []byte) {
  var algo string
  var hashedPass []byte

  salt := genSalt()
  pass := []byte(password)

  if hfunc := time.Now().Unix() % 2; hfunc == 0 {
    algo = BCRYPT
    hashedPass = bCrypt(salt, pass)
  } else if hfunc == 1 {
    algo = SCRYPT
    hashedPass = sCrypt(salt, pass)
  }

  algosalt := append([]byte(algo), salt...)
  return hashedPass, algosalt
}

func checkPassword(hashedPass []byte, password string, algosalt []byte) bool {
  var eq bool

  algo := string(algosalt[0])
  salt := []byte(algosalt[1:])
  pass := []byte(password)

  if algo == BCRYPT {
    eq = bCryptEqual(hashedPass, salt, pass)
  } else if algo == SCRYPT {
    eq = sCryptEqual(hashedPass, salt, pass)
  }

  return eq
}

// input: plain text password
// output:
func HashPassword(password string) (string, string) {
  hash, algosalt := hashPassword(password)
  if hash == nil || algosalt == nil {
    return "", ""
  }

  encodedPass := encode(hash)
  encodedSalt := encode(algosalt)
  if encodedPass == "" || encodedSalt == "" {
    return "", ""
  }

  return encodedPass, encodedSalt
}

// input: base64 encoded hashed password, plain text password,
//        and base64 encoded algosalt
// output: true or false
func ValidatePassword(hashedPassword string, plainPassword string, algoSalt string) bool {
  decodedPass := decode(hashedPassword)
  decodedSalt := decode(algoSalt)
  if decodedPass == nil || decodedSalt == nil {
    return false
  }

  return checkPassword(decodedPass, plainPassword, decodedSalt)
}

