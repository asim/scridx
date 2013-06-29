package app

import (
  "log"
  "net/url"
  "time"
  "fmt"
  "strings"
  "strconv"
  "path/filepath"
  "github.com/eaigner/hood"
  "github.com/gorilla/schema"
  "github.com/hoisie/mustache"
  "scridx/sessions"
)

type FeedbackForm struct {
  Title string
  Text string
  Source string
}

type NewsForm struct {
  Title string
  Text string
  Source string
}

type RequestsForm struct {
  Title string
  Logline string
  Writers string
  Source string
  Imdb string
  Wiki string
  Drafted string
  Version string
  Status int8
}

type ScriptsForm struct {
  Title string
  Logline string
  Writers string
  Source string
  Imdb string
  Wiki string
  Drafted string
  Version string
}

type SignupForm struct {
  Name          string
  Username      string
  Email         string
  Password      string
}

type LoginForm struct {
  Username string
  Password string
  Remember string
}

type SettingsForm struct {
  Name string
  Email string
  Logline string
}

var decoder = schema.NewDecoder()
var layoutPath = filepath.Join("app/views", "layout.m")
var templateCache = make(map[string]*mustache.Template)

func ParseReferer(ref string) *url.URL {
  u, err := url.Parse(ref)
  if err != nil {
     return nil
  }

  return u
}

// Parse Forms
func ParseFormDate(date string) time.Time {
  // Parse Drafted date string in form
  // Expecting "MM-DD-YYYY". 
  // If its blank return blankDateString, blank is valid
  if date == "" {
    t, _ := time.Parse(validDateFormat, blankDateString)
    return t
  }

  if !ValidDateRegex.MatchString(date) {
    t, _ := time.Parse(validDateFormat, invalidDateString)
    return t
  }

  t, _ := time.Parse(validDateFormat, date)
  return t
}

func ParseCommentsForm(form url.Values) *Comments {
  c := new(Comments)
  decoder.Decode(c, form)
  return c
}

func ParseFeedbackForm(form url.Values) *Feedback {
  f := new(Feedback)
  decoder.Decode(f, form)
  return f
}

func ParseNewsForm(form url.Values) *News {
  news := new(News)
  decoder.Decode(news, form)
  return news
}

func ParseScriptForm(form url.Values) *Scripts {
  script := new(Scripts)
  decoder.Decode(script, form)
  return script
}

func ParseRequestForm(form url.Values) *Requests {
  request := new(Requests)
  decoder.Decode(request, form)
  return request
}

func ParseSettingsForm(form url.Values) *SettingsForm {
  sf := new(SettingsForm)
  decoder.Decode(sf, form)
  return sf
}

func ParseSignupForm(form url.Values) *SignupForm {
  sf := new(SignupForm)
  decoder.Decode(sf, form)
  return sf
}

func ParseLoginForm(form url.Values) *LoginForm {
  lf := new(LoginForm)
  decoder.Decode(lf, form)
  return lf
}

// Refill Forms
func refillNewsForm(s *sessions.Session) *NewsForm {
  form := new(NewsForm)

  f := s.GetForm("news")
  if f != nil && len(f) > 0 {
    decoder.Decode(form, f)
  }

  return form
}

func refillFeedbackForm(s *sessions.Session) *FeedbackForm {
  form := new(FeedbackForm)

  f := s.GetForm("feedback")
  if f != nil && len(f) > 0 {
    decoder.Decode(form, f)
  }

  return form
}

func refillSubmit(s *sessions.Session) *ScriptsForm {
  form := new(ScriptsForm)

  f := s.GetForm("script")
  if f != nil && len(f) > 0 {
    decoder.Decode(form, f)
  }

  return form
}

func refillApproveRequest(s *sessions.Session) *RequestsForm {
  form := new(RequestsForm)

  f := s.GetForm("approve")
  if f != nil && len(f) > 0 {
    decoder.Decode(form, f)
  }

  return form
}

func refillFulfillRequest(s *sessions.Session) *RequestsForm {
  form := new(RequestsForm)

  f := s.GetForm("fulfill")
  if f != nil && len(f) > 0 {
    decoder.Decode(form, f)
  }

  return form
}


func refillRequest(s *sessions.Session) *RequestsForm {
  form := new(RequestsForm)

  f := s.GetForm("request")
  if f != nil && len(f) > 0 {
    decoder.Decode(form, f)
  }

  return form
}

func refillSignupForm(s *sessions.Session) *SignupForm {
  form := new(SignupForm)

  f := s.GetForm("signup")
  if f != nil && len(f) > 0 {
    decoder.Decode(form, f)
  }

  return form
}

// Printers

func printVersion(version string) string {
  if len(version) > 0 {
    return version
  }
  return "Unspecified"
}

func printWriters(writers string) string {
  if len(writers) > 0 {
    return writers
  }

  return "Unspecified"
}

func printTitle(title string) string {
  if len(title) > 0 {
    return title
  }

  return "Unspecified Title"
}

func (c Comments) PrintStored() string {
  return timeAgo(c.Stored)
}

func (f Feedback) PrintTitle() string {
  return printTitle(f.Title)
}

func (f Feedback) PrintStored() string {
  return timeAgo(f.Stored)
}

func (n News) PrintTitle() string {
  return printTitle(n.Title)
}

func (n News) PrintStored() string {
  return timeAgo(n.Stored)
}

func (r Requests) PrintTitle() string {
  return printTitle(r.Title)
}

func (r Requests) PrintWriters() string {
  return printWriters(r.Writers)
}

func (r Requests) PrintDraftdate() string {
  return niceDate(r.Drafted)
}

func (r Requests) PrintStored() string {
  return timeAgo(r.Stored)
}

func (r Requests) PrintVersion() string {
  return printVersion(r.Version)
}


func (r Requests) DraftedFormat() string {
  d := r.Drafted.Format(validDateFormat)

  if d == blankDateString {
    return ""
  }

  return d
}

func (r Requests) MetaStatus() string {
  var status string
  aHref := `<li><span class="req-status %s">%s</span></li>`

  switch {
  case r.Status == statusNew:
    info := fmt.Sprintf(`<a href="%s">New</a>`, r.FulfillUrl())
    status = fmt.Sprintf(aHref, "req-new", info)
  case r.Status == statusPending:
    info := fmt.Sprintf(`<a href="%s">Pending</a>`, r.ApproveUrl())
    status = fmt.Sprintf(aHref, "req-pending", info)
  case r.Status == statusCompleted:
    status = fmt.Sprintf(aHref, "req-fulfilled", "Fulfilled")
  }

  return status
}

func metaLinks(imdb, wiki string) string {
  metalinks := make([]string, 2)
  aHref := "<li><a href=\"%s\" title=\"%s\"><img class=\"link-icon\" src=\"%s\"/></a></li>"

  if imdb != "" {
    link := fmt.Sprintf(aHref, escapeUrl(imdb), "Imdb", "/static/img/imdb.png")
    metalinks = append(metalinks, link)
  }

  if wiki != "" {
    link := fmt.Sprintf(aHref, escapeUrl(wiki), "Wikipedia", "/static/img/w.png")
    metalinks = append(metalinks, link)
  }

  if imdb == "" && wiki == "" {
    return ""
  }

  return fmt.Sprintf("<ul class=\"inline meta-icons\">%s</ul>", strings.Join(metalinks, ""))
}

func (r Requests) MetaLinks() string {
  return metaLinks(r.Imdb, r.Wiki)
}

func (s Scripts) MetaLinks() string {
  return metaLinks(s.Imdb, s.Wiki)
}

func (s Scripts) PrintTitle() string {
  return printTitle(s.Title)
}

func (s Scripts) PrintWriters() string {
  return printWriters(s.Writers)
}

func (s Scripts) PrintDraftdate() string {
  return niceDate(s.Drafted)
}

func (s Scripts) PrintStored() string {
  return timeAgo(s.Stored)
}

func (s Scripts) PrintVersion() string {
  return printVersion(s.Version)
}

// Clean/Format Strings
func cleanString(str string) string {
  str = strings.Replace(str, "<", "", -1)
  str = strings.Replace(str, ">", "", -1)
  return strings.TrimSpace(str)
}

func cleanWriters(w string) string {
  w = strings.Replace(w, "&", ",", -1)
  aw := strings.Split(cleanString(w), ",")
  for i, writer := range aw {
    aw[i] = strings.Join(strings.Fields(writer), " ")
  }

  return strings.Join(aw, ", ")
}

func cleanTitle(t string) string {
  return strings.Join(strings.Fields(cleanString(t)), " ")
}

// Called before saving.
func (c *Comments) Clean() {
  c.Text = cleanString(c.Text)
}

func (f *Feedback) Clean() {
  f.Title = cleanTitle(f.Title)
  f.Text = cleanString(f.Text)
  f.Source = cleanString(f.Source)
}

func (n *News) Clean() {
  n.Title = cleanTitle(n.Title)
  n.Text = cleanString(n.Text)
  n.Source = cleanString(n.Source)
}

func (u *Users) Clean() {
  u.Username = cleanString(strings.ToLower(u.Username))
  u.Email = cleanString(strings.ToLower(u.Email))
  u.Name = cleanString(u.Name)
  u.Logline = cleanString(u.Logline)
}

// Called before saving.
func (s *Scripts) Clean() {
  s.Title = cleanTitle(s.Title)
  s.Source = cleanString(s.Source)
  s.Imdb = cleanString(s.Imdb)
  s.Wiki = cleanString(s.Wiki)
  s.Logline = cleanString(s.Logline)
  s.Writers = cleanWriters(s.Writers)
}

// Called before saving.
func (r *Requests) Clean() {
  r.Title = cleanTitle(r.Title)
  r.Source = cleanString(r.Source)
  r.Imdb = cleanString(r.Imdb)
  r.Wiki = cleanString(r.Wiki)
  r.Logline = cleanString(r.Logline)
  r.Writers = cleanWriters(r.Writers)
}

// Pager
func GetPageOffset(vars url.Values, limit int) (int, int) {
  page, err := strconv.Atoi(vars.Get("page"))
  if err != nil {
    page = 1
  }

  next := page - 1
  if page == 1 {
    next = 0
  }

  offset := next * limit
  return page, offset
}

func GetPager(u *url.URL, page, limit, items int) map[string]string {
  pager := make(map[string]string)

  if page == 0 || page == 1 {
    pager["previousPage"] = "#"
    pager["previousState"] = "disabled"
  } else {
    prev := u
    vars := prev.Query()
    vars.Set("page", strconv.Itoa(page - 1))
    prev.RawQuery = vars.Encode()
    pager["previousPage"] = prev.RequestURI()
  }

  if items < limit {
    pager["nextPage"] = "#"
    pager["nextState"] = "disabled"
  } else {
    next := u
    vars := next.Query()
    vars.Set("page", strconv.Itoa(page + 1))
    next.RawQuery = vars.Encode()
    pager["nextPage"] = next.RequestURI()
  }

  return pager
}

// Template methods
func Render(ctx *sessions.Context, view string) {
  viewPath := filepath.Join("app/views", view + ".m")
  r := genTmpl(viewPath, layoutPath, ctx.Data)
  fmt.Fprintf(ctx.W, "%s", r)
}

func genTmpl(view, layout string, data ...interface{}) string {
  tmpl, tprs := templateCache[view]
  lyot, lprs := templateCache[layout]

  if !tprs {
    var err error
    tmpl, err = mustache.ParseFile(view)
    if err != nil {
      log.Println(err)
      return ""
    }
    templateCache[view] = tmpl
  }

  if !lprs {
    var err error
    lyot, err = mustache.ParseFile(layout)
    if err != nil {
      log.Println(err)
      return ""
    }
    templateCache[layout] = lyot
  }

  return tmpl.RenderInLayout(lyot, data...)
}

// Time formatter
func IntToTime(t uint32) time.Time {
  return time.Unix(int64(t), 0)
}

func TimeToInt(t time.Time) uint32 {
  return uint32(t.Unix())
}

func TimeNowToInt() uint32 {
  return uint32(time.Now().Unix())
}

func timeAgo(t uint32) string {
  d := IntToTime(t)
  timeAgo := ""
  startDate := time.Now().Unix()
  deltaMinutes := float64(startDate - d.Unix()) / 60.0
  if deltaMinutes <= 523440  { // less than 363 days
    timeAgo = fmt.Sprintf("%s ago", distanceOfTime(deltaMinutes))
  } else {
    timeAgo = d.Format("2 Jan")
  }

  return timeAgo
}

func distanceOfTime(minutes float64) string {
  switch {
  case minutes < 1:
    return fmt.Sprintf("%d secs", int(minutes * 60))
  case minutes < 59:
    return fmt.Sprintf("%d minutes", int(minutes))
  case minutes < 90:
    return "about an hour"
  case minutes < 120:
    return "almost 2 hours"
  case minutes < 1080:
    return fmt.Sprintf("%d hours", int(minutes / 60))
  case minutes < 1680:
    return "about a day"
  case minutes < 2160:
    return "more than a day"
  case minutes < 2520:
    return "almost 2 days"
  case minutes < 2880:
    return "about 2 days"
  default:
    return fmt.Sprintf("%d days", int(minutes / 1440))
  }

  return ""
}

// Date formatters
func niceIntDate(date uint32) string {
  if date == blankDateInt {
    return "Unspecified"
  }

  return niceDate(IntToTime(date))
}

func niceDate(date time.Time) string {
  d := date.Format("Jan 1, 2006")
  if !NiceDateRegex.MatchString(d) {
    d = "Unspecified"
  }

  return d
}

// URL Methods
func escapeUrl(u string) string {
  ul,_ := url.Parse(u)
  return ul.String()
}

func createUrl(id hood.Id, title, u string) string {
  title = SafeUrlRegex.ReplaceAllString(title, "_")
  title = strings.ToLower(strings.Trim(title, "_"))
  if len(title) > 64 {
    title = title[:64]
  }
  return fmt.Sprintf(u, id, title)
}

func sourceOrUrl(s, u string) string {
  if s != "" {
    return s
  }

  return u
}
func (s Scripts) SourceOrUrl() string {
  return sourceOrUrl(s.Source, s.Url())
}

func (r Requests) SourceOrUrl() string {
  return sourceOrUrl(r.Source, r.Url())
}

func (f Feedback) SourceOrUrl() string {
  return sourceOrUrl(f.Source, f.Url())
}

func (n News) SourceOrUrl() string {
  return sourceOrUrl(n.Source, n.Url())
}

func (f Feedback) Url() string {
  return createUrl(f.Id, f.Title, "/feedback/%d/%s")
}

func (n News) Url() string {
  return createUrl(n.Id, n.Title, "/news/%d/%s")
}

func (s Scripts) Url() string {
  return createUrl(s.Id, s.Title, "/scripts/%d/%s")
}

func (r Requests) Url() string {
  return createUrl(r.Id, r.Title, "/requests/%d/%s")
}

func (u *Users) Url() string {
  return fmt.Sprintf("/u/%s", u.Username)
}

func (c Comments) Url() string {
  return fmt.Sprintf("/comments/%d", c.Id)
}

func (s Scripts) VoteUrl() string {
  return fmt.Sprintf("/vote/script/%d", s.Id)
}

func (r Requests) VoteUrl() string {
  return fmt.Sprintf("/vote/request/%d", r.Id)
}

func (n News) VoteUrl() string {
  return fmt.Sprintf("/vote/news/%d", n.Id)
}

func (f Feedback) VoteUrl() string {
  return fmt.Sprintf("/vote/feedback/%d", f.Id)
}

func (c Comments) VoteUrl() string {
  return fmt.Sprintf("/vote/comment/%d", c.Id)
}

func (r Requests) FulfillUrl() string {
  title := SafeUrlRegex.ReplaceAllString(r.Title, "_")
  title = strings.ToLower(strings.Trim(title, "_"))
  return fmt.Sprintf("/submit/fulfill/r/%d/%s", r.Id, title)
}

func (r Requests) ApproveUrl() string {
  title := SafeUrlRegex.ReplaceAllString(r.Title, "_")
  title = strings.ToLower(strings.Trim(title, "_"))
  return fmt.Sprintf("/submit/approve/r/%d/%s", r.Id, title)
}


