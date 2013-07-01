package app

import (
  "regexp"
  "github.com/eaigner/hood"
  "time"
)

const (
  invalidUrl = 0
  invalidText = 1
  invalidDate = 2
  invalidUsername = 3
  invalidEmail = 4
  invalidName = 5
  blankDateInt = 0
  invalidDateInt = 1
  blankDateString = "01/01/0001"
  invalidDateString = "01/01/0011"
  validDateFormat = "01/02/2006"
  validDateRegex = `^(0[1-9]|1[012])\/(0[1-9]|[12][0-9]|3[01])\/(19|20)\d\d$`
  validWritersRegex = `^[[:alpha:]\s-,.]+$`
  validUsernameLength = 15
  validEmailLength = 254
  validNameLength = 20
  validMinPasswordLength = 8
  validUsernameRegex = `^[a-z0-9_]{1,15}$`
  validEmailRegex = `^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,4}$`
  validNameRegex = `^[[:alpha:]\s-_&@,.]{1,20}$`
)

var (
  ValidUsernameRegex = regexp.MustCompile(validUsernameRegex)
  ValidEmailRegex = regexp.MustCompile(validEmailRegex)
  ValidNameRegex = regexp.MustCompile(validNameRegex)
  ValidDateRegex = regexp.MustCompile(validDateRegex)
  ValidWritersRegex = regexp.MustCompile(validWritersRegex)
  ValidImdbRegex = regexp.MustCompile(`^https?://(www\.)?imdb.com/(title/tt|name/nm)[0-9]+`)
  ValidWikiRegex = regexp.MustCompile(`^https?://(www\.|en\.)?wikipedia.org/wiki/.+`)
  ValidSourceRegex = regexp.MustCompile(`^(https?:\/\/)?([\da-z\.-]+)\.([a-z\.]{2,6})([\/\w \.-]*)*\/?`)
  NiceDateRegex = regexp.MustCompile(`^[A-Z][a-z]{2} [0-9]{1,2}, (19|20)\d\d$`)
  SafeUrlRegex  = regexp.MustCompile("[^A-Za-z0-9]+")
  ValidVersionRegex = regexp.MustCompile(`^[A-Za-z0-9\s]{1,20}$`)
)

// Validation functions
func formatIntDate(date uint32) string {
  d := IntToTime(date)
  return d.Format(validDateFormat)
}

func validatePassword(password string) error {
  if len(password) < validMinPasswordLength {
    return hood.NewValidationError(invalidName, "Password needs to be at least 8 characters long")
  }

  return nil
}

func validateLength(field string, length int) bool {
  if l := len(field); l > length || l == 0 {
    return false
  }

  return true
}

func validateTitle(title string) error {
  if len(title) <= 0 {
    return hood.NewValidationError(invalidText, "Invalid Title.")
  }
  return nil
}

func validateWriters(writers string) error {
  if len(writers) > 0 {
    if !ValidWritersRegex.MatchString(writers) {
      return hood.NewValidationError(invalidText, "Invalid characters in Writers field")
    }
  }

  return nil
}

func validateLogline(logline string) error {
  if len(logline) > 512 {
    return hood.NewValidationError(invalidText, "Logline must be less than 512 characters.")
  }
  return nil
}

func validateDate(date string) error {
  if !ValidDateRegex.MatchString(date) {
    return hood.NewValidationError(invalidDate, "Invalid date, require MM/DD/YYYY")
  }

  return nil
}

func validateDraftDate(date time.Time) error {
  d := date.Format(validDateFormat)
  if d == blankDateString {
    return nil
  }

  if d == invalidDateString {
    return hood.NewValidationError(invalidDate, "Invalid draft date, require MM/DD/YYYY")
  }

  return validateDate(d)
}

func validateImdbLink(link string) error {
  if len(link) > 0 {
    if !ValidImdbRegex.MatchString(link) {
      return hood.NewValidationError(invalidUrl, "Invalid IMDB Url")
    }
  }

  return validateUrl(link)
}


func validateWikiLink(link string) error {
  if len(link) > 0 {
    if !ValidWikiRegex.MatchString(link) {
      return hood.NewValidationError(invalidUrl, "Invalid Wiki Url")
    }
  }

  return validateUrl(link)
}

func validateSourceLink(link string) error {
    if !ValidSourceRegex.MatchString(link) {
      return hood.NewValidationError(invalidUrl, "Invalid Source Url")
    }
    return validateUrl(link)
}

func validateUrl(link string) error {
  if len(link) > 255 {
    return hood.NewValidationError(invalidUrl, "URL must be less than 255 characters")
  }
  return nil
}

func validateVersion(version string) error {
  if len(version) > 0 {
    if !ValidVersionRegex.MatchString(version) {
      return hood.NewValidationError(invalidText, "Invalid Version. 20 characters or less.")
    }
  }

  return nil
}

// Validation methods

// Feedback validations
func (f *Feedback) ValidateTitle() error {
  return validateTitle(f.Title)
}

func (f *Feedback) ValidateSourceLink() error {
  if len(f.Source) == 0 {
    return nil
  }

  return validateSourceLink(f.Source)
}

func (f *Feedback) ValidateText() error {
  if len(f.Text) > 10000 {
    return hood.NewValidationError(invalidText, "Text needs to be less than 10,000 characters")
  }

  return nil
}

// Comments validations
func (c *Comments) ValidateText() error {
  if len(c.Text) == 0 {
    return hood.NewValidationError(invalidText, "Comment cannot be blank.")
  }

  if len(c.Text) > 10000 {
    return hood.NewValidationError(invalidText, "Comment needs to be less than 10,000 characters")
  }

  return nil
}

// News validations
func (n *News) ValidateTitle() error {
  return validateTitle(n.Title)
}

func (n *News) ValidateSourceLink() error {
  if len(n.Source) == 0 {
    return nil
  }

  return validateSourceLink(n.Source)
}

func (n *News) ValidateText() error {
  if len(n.Text) > 10000 {
    return hood.NewValidationError(invalidText, "Text needs to be less than 10,000 characters")
  }

  return nil
}

// Script validations
func (s *Scripts) ValidateVersion() error {
  return validateVersion(s.Version)
}

func (s *Scripts) ValidateDraftDate() error {
  return validateDraftDate(s.Drafted)
}

func (s *Scripts) ValidateStored() error {
  d := formatIntDate(s.Stored)
  return validateDate(d)
}

func (s *Scripts) ValidateTitle() error {
  return validateTitle(s.Title)
}

func (s *Scripts) ValidateWriters() error {
  return validateWriters(s.Writers)
}

func (s *Scripts) ValidateLogline() error {
    return validateLogline(s.Logline)
}

func (s *Scripts) ValidateImdbLink() error {
  return validateImdbLink(s.Imdb)
}

func (s *Scripts) ValidateWikiLink() error {
  return validateWikiLink(s.Wiki)
}

func (s *Scripts) ValidateSourceLink() error {
    return validateSourceLink(s.Source)
}

// Request validations

func (r *Requests) ValidateVersion() error {
  return validateVersion(r.Version)
}

func (r *Requests) ValidateDraftDate() error {
  return validateDraftDate(r.Drafted)
}

func (r *Requests) ValidateStoredDate() error {
  d := formatIntDate(r.Stored)
  return validateDate(d)
}

func (r *Requests) ValidateTitle() error {
  return validateTitle(r.Title)
}

func (r *Requests) ValidateWriters() error {
  return validateWriters(r.Writers)
}

func (r *Requests) ValidateLogline() error {
    return validateLogline(r.Logline)
}

func (r *Requests) ValidateImdbLink() error {
  return validateImdbLink(r.Imdb)
}

func (r *Requests) ValidateWikiLink() error {
  return validateWikiLink(r.Wiki)
}

func (r *Requests) ValidateSourceLink() error {
  if r.Status == statusNew {
    return nil
  }

  return validateSourceLink(r.Source)
}

// User Validations
func (u *Users) ValidateLogline() error {
  if len(u.Logline) > 160 {
    return hood.NewValidationError(invalidText, "Your logline needs to be less than 160 characters.")
  }

  return nil
}

func (u *Users) ValidateUsername() error {
  if !validateLength(u.Username, validUsernameLength) {
    return hood.NewValidationError(invalidUsername, "Username needs to be 1-15 characters.")
  }

  if !ValidUsernameRegex.MatchString(u.Username) {
    return hood.NewValidationError(invalidUsername, "Invalid Username, must be alphanumeric.")
  }

  return nil
}

func (u *Users) ValidateEmail() error {
  if !validateLength(u.Email, validEmailLength) {
    return hood.NewValidationError(invalidEmail, "Require an email address, should be less than 255 characters.")
  }

  if !ValidEmailRegex.MatchString(u.Email) {
    return hood.NewValidationError(invalidEmail, "Invalid email address.")
  }

  return nil
}

func (u *Users) ValidateName() error {
  if !validateLength(u.Name, validNameLength) {
    return hood.NewValidationError(invalidName, "Name needs to be 1-20 characters.")
  }

  if !ValidNameRegex.MatchString(u.Name) {
    return hood.NewValidationError(invalidName, "Invalid characters in name.")
  }

  return nil
}
