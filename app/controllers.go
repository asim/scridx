package app

import (
  "log"
  "fmt"
  "net/http"
  "strings"
  "github.com/gorilla/mux"
  "scridx/crypt"
  "scridx/sessions"
)

type AppError struct {
    Error   error
    Message string
    Code    int
}

type Handler func(*sessions.Context) *AppError

func errorHandler(ctx *sessions.Context, err *AppError) {
  // Display error page
  ctx.W.WriteHeader(err.Code)
  d := *ctx.Data
  d["message"] = err.Message
  Render(ctx, "error")
}

func NotFoundHandler(ctx *sessions.Context) *AppError {
  log.Printf("%s %s %s %s", ctx.R.RemoteAddr, ctx.R.Method, ctx.R.URL.Path, ctx.R.Proto)
  errorHandler(ctx, &AppError{nil, "Page not found", 404})
  return nil
}

// Global HTTP Handler
func (fn Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  if r == nil {
    log.Println("oh shit")
  }

  log.Printf("%s %s %s %s", r.RemoteAddr, r.Method, r.URL.RequestURI(), r.Proto)

  ctx, err := sessions.NewContext(w, r)
  if err != nil {
    log.Printf(err.Error())
    w.WriteHeader(502)
    fmt.Fprintf(w, "An error has occcurred")
    return
  }

  if err := fn(ctx); err != nil { // e is *Error, not os.Error.
    if err.Error != nil {
      log.Printf("%v", err.Error)
    } else {
      log.Printf("%s", err.Message)
    }

   errorHandler(ctx, err)
  }
}

func parseForm(ctx *sessions.Context, form string) bool {
  ctx.R.ParseForm()
  if !ctx.Session.IsValidCSRF(ctx.R.FormValue("_csrf")) {
    if form != "" {
      formError(ctx, form, "Refresh the page and try again.")
    }

    return false
  }

  return true
}

func formError(ctx *sessions.Context, form, errmsg string) {
  ctx.Session.SetFlash(ctx.W, ctx.R, errmsg)
  ctx.Session.SetForm(ctx.W, ctx.R, form, ctx.R.Form)
  http.Redirect(ctx.W, ctx.R, ctx.R.Referer(), 302)
}

// Request Handler methods
func FaviconHandler(w http.ResponseWriter, r *http.Request) {
  http.ServeFile(w, r, "app/static/img/favicon.ico")
}

func RobotHandler(w http.ResponseWriter, r *http.Request) {
  http.ServeFile(w, r, "app/static/html/robots.txt")
}

// GET /
func IndexHandler(ctx *sessions.Context) *AppError {
  ListScripts(ctx)
  Render(ctx, "index")
  return nil
}
/*func IndexHandler(ctx *sessions.Context) *AppError {
  vars := ctx.R.URL.Query()
  limit := 20
  d := *ctx.Data

  page, offset := GetPageOffset(vars, limit)

  var entities []*LRUItem
  switch vars.Get("order") {
  case "top":
    entities, _ = GetLRUItems(limit, offset)
  case "latest":
    entities, _ = GetLRUItems(limit, offset)
  default:
    entities, _ = GetLRUItems(limit, offset)
  }

  if ctx.UserSession.LoggedIn() {
    id := d["_user"].(*sessions.Token).Id
    for i, entity := range entities {
      var user_id int64
      switch entity.Thing.(type) {
      case Scripts:
        user_id = entity.Thing.(Scripts).UserId
      case Requests:
        user_id = entity.Thing.(Requests).UserId
      case Feedback:
        user_id = entity.Thing.(Feedback).UserId
      case News:
        user_id = entity.Thing.(News).UserId
      }

      if id == user_id || entity.Votes[id] {
        entities[i].Votable = false
      } else {
        entities[i].Votable = true
      }
    }
  }

  pager := GetPager(ctx.R.URL, page, limit, len(entities))

  d["entities"] = entities
  d["pager"] = &pager
  Render(ctx, "index")
  return nil
}
*/

func ListHandler(ctx *sessions.Context) *AppError {
  switch ctx.R.URL.Path {
  case "/comments":
    ListComments(ctx)
    Render(ctx, "comments")
  case "/feedback":
    ListFeedback(ctx)
    Render(ctx, "feedback")
  case "/scripts":
    ListScripts(ctx)
    Render(ctx, "scripts")
  case "/requests":
    ListRequests(ctx)
    Render(ctx, "requests")
  case "/news":
    ListNews(ctx)
    Render(ctx, "news")
  default:
    return &AppError{nil, "Page not found", 404}
  }

  return nil
}

// GET /vote/script/id
func VoteHandler(ctx *sessions.Context) *AppError {
  // Is the user logged in?
  if !ctx.UserSession.LoggedIn() {
    ctx.Session.SetFlash(ctx.W, ctx.R, "You need to login first.") 
    http.Redirect(ctx.W, ctx.R, ctx.R.Referer(), 302)
    return nil
  }

  d := *ctx.Data
  user := d["_user"].(*sessions.Token)
  vars := mux.Vars(ctx.R)

  type item interface {
    Vote(int64, int) error
  }

  var ok bool
  var vote = "up"

  var voteFn = func(i item, user_id int64, exists bool) bool {
    if !exists || user.Id == user_id {
     return false
    }

    var err error
    switch vote {
    case "up":
      err = i.Vote(user.Id, 1)
    case "down":
      err = i.Vote(user.Id, -1)
    }

    if err != nil {
      return false
    }

    return true
  }

  switch vars["entity"] {
  case "script":
    i, exists := GetScript(vars["id"])
    ok = voteFn(i, i.UserId, exists)
  case "request":
    i, exists := GetRequest(vars["id"])
    ok = voteFn(i, i.UserId, exists)
  case "feedback":
    i, exists := GetFeedback(vars["id"])
    ok = voteFn(i, i.UserId, exists)
  case "news":
    i, exists := GetNews(vars["id"])
    ok = voteFn(i, i.UserId, exists)
  case "comment":
    // Only comments can be downvoted
    if vars["dir"] == "down" {
      vote = "down"
    }

    i, exists := GetComment(vars["id"])
    ok = voteFn(i, i.UserId, exists)
  default:
    return &AppError{nil, "Not found", 404}
  }

  if !ok {
    return &AppError{nil, "Error voting", 404}
  }

  // Get the item and update its vote
  return nil
}


func CommentHandler(ctx *sessions.Context) *AppError {
  if ctx.R.Method != "GET" {
    return nil
  }

  vars := mux.Vars(ctx.R)
  if comment, ok := GetComment(vars["id"]); ok {
    d := *ctx.Data
    votable := true

    // If logged in get votes
    if ctx.UserSession.LoggedIn() {
      user_id := d["_user"].(*sessions.Token).Id
      // Get votes and create map
      votes, ok := GetCommentVotes(user_id, []string{fmt.Sprintf("%d",comment.Id)})
      if ok && len(*votes) > 0 || user_id == comment.UserId {
        votable = false
      }
    }

    d["comment"] = createEntity(0, votable, comment)
    d["comments"] = threadedComments(comment.EntityType, comment.EntityId, int64(comment.Id), ctx)
    if user, ok := GetUserById(comment.UserId); ok {
      d["submitter"] = user
    }

    et := fmt.Sprintf("%d", comment.EntityId)
    switch comment.EntityType {
    case ScriptType:
      p, _ := GetScript(et)
      d["parent_text"] = p.Title
      d["parent_url"] = p.Url()
    case RequestType:
      p, _ := GetRequest(et)
      d["parent_text"] = p.Title
      d["parent_url"] = p.Url()
    case FeedbackType:
      p, _ := GetFeedback(et)
      d["parent_text"] = p.Title
      d["parent_url"] = p.Url()
    case NewsType:
      p, _ := GetNews(et)
      d["parent_text"] = p.Title
      d["parent_url"] = p.Url()
    }


    Render(ctx, "comment")
  } else {
    return &AppError{nil, "Page not found", 404}
  }

  return nil
}

func FeedbackHandler(ctx *sessions.Context) *AppError {
  if ctx.R.Method != "GET" {
    return nil
  }

  vars := mux.Vars(ctx.R)
  if feedback, ok := GetFeedback(vars["id"]); ok {
    d := *ctx.Data

    votable := true

    // If logged in get votes
    if ctx.UserSession.LoggedIn() {
      user_id := d["_user"].(*sessions.Token).Id
      // Get votes and create map
      votes, ok := GetFeedbackVotes(user_id, []string{fmt.Sprintf("%d",feedback.Id)})
      if ok && len(*votes) > 0 || user_id == feedback.UserId {
        votable = false
      }
    }

    d["feedback"] = createEntity(0, votable, feedback)
    d["comment"] = &Comments{
      ParentId: 0,
      EntityId: int64(feedback.Id),
      EntityType: FeedbackType,
    }

    d["comments"] = threadedComments(int16(FeedbackType), int64(feedback.Id), 0, ctx)

    if user, ok := GetUserById(feedback.UserId); ok {
      d["submitter"] = user
    }
    Render(ctx, "feedback_item")
  } else {
    return &AppError{nil, "Page not found", 404}
  }

  return nil
}

func NewsHandler(ctx *sessions.Context) *AppError {
  if ctx.R.Method != "GET" {
    return nil
  }

  vars := mux.Vars(ctx.R)
  if news, ok := GetNews(vars["id"]); ok {
    d := *ctx.Data
    votable := true

    // If logged in get votes
    if ctx.UserSession.LoggedIn() {
      user_id := d["_user"].(*sessions.Token).Id
      // Get votes and create map
      votes, ok := GetNewsVotes(user_id, []string{fmt.Sprintf("%d",news.Id)})
      if ok && len(*votes) > 0 || user_id == news.UserId {
        votable = false
      }
    }

    d["news"] = createEntity(0, votable, news)
    d["comment"] = &Comments{
      ParentId: 0,
      EntityId: int64(news.Id),
      EntityType: NewsType,
    }

    d["comments"] = threadedComments(int16(NewsType), int64(news.Id), 0, ctx)

    if user, ok := GetUserById(news.UserId); ok {
      d["submitter"] = user
    }
    Render(ctx, "news_item")
  } else {
    return &AppError{nil, "Page not found", 404}
  }

  return nil
}

// GET /script/id/title
func ScriptHandler(ctx *sessions.Context) *AppError {
  if ctx.R.Method != "GET" {
    return nil
  }

  vars := mux.Vars(ctx.R)
  if script, ok := GetScript(vars["id"]); ok {
    d := *ctx.Data
    votable := true

    // If logged in get votes
    if ctx.UserSession.LoggedIn() {
      user_id := d["_user"].(*sessions.Token).Id
      // Get votes and create map
      votes, ok := GetScriptVotes(user_id, []string{fmt.Sprintf("%d",script.Id)})
      if ok && len(*votes) > 0 || user_id == script.UserId {
        votable = false
      }
    }

    d["script"] = createEntity(0, votable, script)
    d["comment"] = &Comments{
      ParentId: 0,
      EntityId: int64(script.Id),
      EntityType: ScriptType,
    }

    d["comments"] = threadedComments(int16(ScriptType), int64(script.Id), 0, ctx)

    if user, ok := GetUserById(script.UserId); ok {
      d["submitter"] = user
    }
    Render(ctx, "script")
  } else {
    return &AppError{nil, "Page not found", 404}
  }

  return nil
}

// GET /request/id/title
func RequestHandler(ctx *sessions.Context) *AppError {
  if ctx.R.Method != "GET" {
    return nil
  }

  vars := mux.Vars(ctx.R)
  if request, ok := GetRequest(vars["id"]); ok {
    d := *ctx.Data
    votable := true

    // If logged in get votes
    if ctx.UserSession.LoggedIn() {
      user_id := d["_user"].(*sessions.Token).Id
      // Get votes and create map
      votes, ok := GetRequestVotes(user_id, []string{fmt.Sprintf("%d",request.Id)})
      if ok && len(*votes) > 0 || user_id == request.UserId {
        votable = false
      }
    }

    d["request"] = createEntity(0, votable, request)
    d["comment"] = &Comments{
      ParentId: 0,
      EntityId: int64(request.Id),
      EntityType: RequestType,
    }

    d["comments"] = threadedComments(int16(RequestType), int64(request.Id), 0, ctx)

    if user, ok := GetUserById(request.UserId); ok {
      d["submitter"] = user
    }

    if user, ok := GetUserById(request.FulfillUserId); ok {
      d["fulfiller"] = user
    }

    Render(ctx, "request")
  } else {
    return &AppError{nil, "Page not found", 404}
  }

  return nil
}

// GET|POST /submit/{form}
// GET|POST /submit/{form}/r/id/title
func FormHandler(ctx *sessions.Context) *AppError {
  if !ctx.UserSession.LoggedIn() {
    // avoid redirect loop
    url := ctx.R.Referer()
    if strings.HasSuffix(url, ctx.R.URL.Path) {
      url = "/"
    } else {
      ctx.Session.SetFlash(ctx.W, ctx.R, "You need to login first.")
    }

    http.Redirect(ctx.W, ctx.R, url, 302)
    return nil
  }

  // Make sure anything changing data is coming from a form on this site.
  if ctx.R.Method == "POST" {
    ref := ParseReferer(ctx.R.Referer())
    if ref == nil || ref.Host != ctx.R.Host {
      return &AppError{nil, "Not authorized", 401}
    }
  }

  vars := mux.Vars(ctx.R)

  switch vars["form"] {
  case "comment":
    return submitCommentHandler(ctx)
  case "feedback":
    return feedbackFormHandler(ctx)
  case "news":
    return newsFormHandler(ctx)
  case "script":
    return scriptFormHandler(ctx)
  case "request":
    return requestFormHandler(ctx)
  case "fulfill":
    return fulfillFormHandler(ctx)
  case "approve":
    return approveFormHandler(ctx)
  default:
    return &AppError{nil, "Page not found", 404}
  }
}

func submitCommentHandler(ctx *sessions.Context) *AppError {
  if ctx.R.Method == "POST" {
    if !parseForm(ctx, "") {
      log.Println("error submitting comment")
      ctx.Session.SetFlash(ctx.W, ctx.R, "Refresh the page and try again")
      http.Redirect(ctx.W, ctx.R, ctx.R.Referer(), 302)
      return nil
    }

    form := ParseCommentsForm(ctx.R.Form)
    d := *ctx.Data

    c := &Comments{
      Text: form.Text,
      UserId: d["_user"].(*sessions.Token).Id,
      ParentId: form.ParentId,
      EntityId: form.EntityId,
      EntityType: form.EntityType,
      Stored: TimeNowToInt(),
    }

    c.Clean()

    if err := c.Save(); err != nil {
      formError(ctx, "comment", err.Error())
      return nil
    }

    ctx.Session.SetFlash(ctx.W, ctx.R, "Comment submitted!")
    ctx.Session.ClearForm(ctx.W, ctx.R, "comment")
    http.Redirect(ctx.W, ctx.R, ctx.R.Referer(), 302)
    return nil
  }

  return nil
}

// GET|POST /submit/approve/r/id/title
func approveFormHandler(ctx *sessions.Context) *AppError {
  vars := mux.Vars(ctx.R)
  request, ok := GetRequest(vars["id"])
  if !ok {
    return &AppError{nil, "Page not found", 404}
  }

  if request.Status != statusPending {
    ctx.Session.SetFlash(ctx.W, ctx.R, "This request might need to be fulfilled first or its already been approved.")
    http.Redirect(ctx.W, ctx.R, request.Url(), 302)
    return nil
  }

  d := *ctx.Data
  if d["_user"].(*sessions.Token).Id != request.UserId {
    ctx.Session.SetFlash(ctx.W, ctx.R, "You can't approve this request.")
    http.Redirect(ctx.W, ctx.R, request.Url(), 302)
    return nil
  }

  if ctx.R.Method == "GET" {
    if len(d["info"].([]interface{})) > 0 {
      d["form"] = refillApproveRequest(ctx.Session)
    } else {
      d["form"] = request
    }

    Render(ctx, "approveRequestForm")
    return nil
  }

  if ctx.R.Method == "POST" {
    if !parseForm(ctx, "approve") {
      return nil
    }

    approve := ParseRequestForm(ctx.R.Form)
    if approve.Title != request.Title {
      return &AppError{nil, "Page not found", 404}
    }

    approve.Drafted = ParseFormDate(ctx.R.FormValue("Drafted"))
    approve.Stored = request.Stored
    approve.Status = statusCompleted
    approve.Title = request.Title
    approve.Source = request.Source
    approve.Id = request.Id
    approve.UserId = request.UserId
    approve.FulfillUserId = request.FulfillUserId
    approve.Clean()

    if err := approve.Save(); err != nil {
      formError(ctx, "approve", err.Error())
      return nil
    }

    // The request has been approved. Now save as a new script.
    script := &Scripts{}
    script.Title = approve.Title
    script.Logline = approve.Logline
    script.Writers = approve.Writers
    script.Source = approve.Source
    script.Imdb = approve.Imdb
    script.Stored = TimeNowToInt()
    script.Drafted = approve.Drafted

    // Check if the source already exists.
    if _, exists := script.CheckIfSourceExists(); !exists {
      script.Save()
    }

    ctx.Session.SetFlash(ctx.W, ctx.R, "Thanks! The request has now been approved and everyone can read the script.")
    ctx.Session.ClearForm(ctx.W, ctx.R, "approve")
    url := request.Url()
    http.Redirect(ctx.W, ctx.R, url, 302)
  }

  return nil
}

// GET|POST /submit/fulfill/r/id/title
func fulfillFormHandler(ctx *sessions.Context) *AppError {
  vars := mux.Vars(ctx.R)
  request, ok := GetRequest(vars["id"])
  if !ok {
    return &AppError{nil, "Page not found", 404}
  }


  if request.Status != statusNew {
    ctx.Session.SetFlash(ctx.W, ctx.R, "This request might have already been fulfilled or its waiting for approval.")
    http.Redirect(ctx.W, ctx.R, request.Url(), 302)
  }

  if ctx.R.Method == "GET" {
    d := *ctx.Data
    if len(d["info"].([]interface{})) > 0 {
      d["form"] = refillFulfillRequest(ctx.Session)
    } else {
      d["form"] = request
    }

    Render(ctx, "fulfillRequestForm")
    return nil
  }

  if ctx.R.Method == "POST" {
    if !parseForm(ctx, "fulfill") {
      return nil
    }

    fulfill := ParseRequestForm(ctx.R.Form)
    if fulfill.Title != request.Title {
      return &AppError{nil, "Page not found", 404}
    }

    d := *ctx.Data

    fulfill.Drafted = ParseFormDate(ctx.R.FormValue("Drafted"))
    fulfill.Stored = request.Stored
    fulfill.Status = statusPending
    fulfill.Title = request.Title
    fulfill.Id = request.Id
    fulfill.UserId = request.UserId
    fulfill.FulfillUserId = d["_user"].(*sessions.Token).Id
    fulfill.Clean()
    log.Println(fulfill)
    if err := fulfill.Save(); err != nil {
      formError(ctx, "fulfill", err.Error())
      return nil
    }

    ctx.Session.SetFlash(ctx.W, ctx.R, "Thanks! The script is now pending approval.")
    ctx.Session.ClearForm(ctx.W, ctx.R, "fulfill")
    url := request.Url()
    http.Redirect(ctx.W, ctx.R, url, 302)
  }

  return nil
}

// GET|POST /submit/feedback
func feedbackFormHandler(ctx *sessions.Context) *AppError {
  d := *ctx.Data
  if ctx.R.Method == "GET" {
    d["form"] = new(FeedbackForm)
    if len(d["info"].([]interface{})) > 0 {
      d["form"] = refillFeedbackForm(ctx.Session)
    }

    Render(ctx, "feedbackForm")
    return nil
  }

  if ctx.R.Method == "POST" {
    if !parseForm(ctx, "feedback") {
      return nil
    }

    f := ParseFeedbackForm(ctx.R.Form)
    f.Stored = TimeNowToInt()
    f.UserId = d["_user"].(*sessions.Token).Id
    f.Clean()

    if err := f.Save(); err != nil {
      formError(ctx, "feedback", err.Error())
      return nil
    }

    ctx.Session.SetFlash(ctx.W, ctx.R, "Feedback request created!")
    ctx.Session.ClearForm(ctx.W, ctx.R, "feedback")
    url := f.Url()
    http.Redirect(ctx.W, ctx.R, url, 302)
  }

  return nil
}

// GET|POST /submit/news
func newsFormHandler(ctx *sessions.Context) *AppError {
  d := *ctx.Data
  if ctx.R.Method == "GET" {
    d["form"] = new(NewsForm)
    if len(d["info"].([]interface{})) > 0 {
      d["form"] = refillNewsForm(ctx.Session)
    }

    Render(ctx, "newsForm")
    return nil
  }

  if ctx.R.Method == "POST" {
    if !parseForm(ctx, "news") {
      return nil
    }

    news := ParseNewsForm(ctx.R.Form)
    news.Stored = TimeNowToInt()
    news.UserId = d["_user"].(*sessions.Token).Id
    news.Clean()

    if err := news.Save(); err != nil {
      formError(ctx, "news", err.Error())
      return nil
    }

    ctx.Session.SetFlash(ctx.W, ctx.R, "News item created!")
    ctx.Session.ClearForm(ctx.W, ctx.R, "news")
    url := news.Url()
    http.Redirect(ctx.W, ctx.R, url, 302)
  }

  return nil
}

// GET|POST /submit/request
func requestFormHandler(ctx *sessions.Context) *AppError {
  d := *ctx.Data
  if ctx.R.Method == "GET" {
    d["form"] = new(RequestsForm)
    if len(d["info"].([]interface{})) > 0 {
      d["form"] = refillRequest(ctx.Session)
    }

    Render(ctx, "requestForm")
    return nil
  }

  if ctx.R.Method == "POST" {
    if !parseForm(ctx, "request") {
      return nil
    }

    request := ParseRequestForm(ctx.R.Form)
    request.Drafted = ParseFormDate(ctx.R.FormValue("Drafted"))
    request.Stored = TimeNowToInt()
    request.UserId = d["_user"].(*sessions.Token).Id
    request.Clean()

    if err := request.Save(); err != nil {
      formError(ctx, "request", err.Error())
      return nil
    }

    ctx.Session.SetFlash(ctx.W, ctx.R, "Thanks! You're script request is now waiting to be fulfilled.")
    ctx.Session.ClearForm(ctx.W, ctx.R, "request")
    url := request.Url()
    http.Redirect(ctx.W, ctx.R, url, 302)
  }

  return nil
}

// GET|POST /submit/script
func scriptFormHandler(ctx *sessions.Context) *AppError {
  d := *ctx.Data
  if ctx.R.Method == "GET" {
    d["form"] = new(ScriptsForm)
    if len(d["info"].([]interface{})) > 0 || len(d["info"].([]interface{})) > 0 {
      d["form"] = refillSubmit(ctx.Session)
    }

    Render(ctx, "scriptForm")
    return nil
  }

  if ctx.R.Method == "POST" {
    if !parseForm(ctx, "script") {
      return nil
    }

    script := ParseScriptForm(ctx.R.Form)
    script.Drafted = ParseFormDate(ctx.R.FormValue("Drafted"))
    script.Stored = TimeNowToInt()
    script.UserId = d["_user"].(*sessions.Token).Id
    script.Clean()

    // Check if the source already exists.
    if s, exists := script.CheckIfSourceExists(); exists {
     msg := fmt.Sprintf("Looks like we've got that script. <a href=\"%s\"><strong>Check Here</strong></a>", s.Url())
     ctx.Session.SetFlash(ctx.W, ctx.R, msg)
     ctx.Session.SetForm(ctx.W, ctx.R, "script", ctx.R.Form)
     http.Redirect(ctx.W, ctx.R, ctx.R.Referer(), 302)
     return nil
    }

    if err := script.Save(); err != nil {
      formError(ctx, "script", err.Error())
      return nil
    }

    ctx.Session.SetFlash(ctx.W, ctx.R, "Thanks! The script is now available for everyone to read.")
    ctx.Session.ClearForm(ctx.W, ctx.R, "script")
    url := script.Url()
    http.Redirect(ctx.W, ctx.R, url, 302)
  }

  return nil
}

func loginError(ctx *sessions.Context, errmsg string) {
  log.Println(errmsg)
  ctx.Session.SetFlash(ctx.W, ctx.R, "Invalid username or password")
  http.Redirect(ctx.W, ctx.R, ctx.R.Referer(), 302)
}

func signupError(ctx *sessions.Context, errmsg string) {
  form := ctx.R.Form
  form["Password"] = []string{""}
  form["VerifyPassword"] = []string{""}

  log.Println(errmsg)
  ctx.Session.SetFlash(ctx.W, ctx.R, errmsg)
  ctx.Session.SetForm(ctx.W, ctx.R, "signup", form)
  http.Redirect(ctx.W, ctx.R, ctx.R.Referer(), 302)
}

func UserHandler(ctx *sessions.Context) *AppError {
  vars := mux.Vars(ctx.R)
  content := vars["content"]
  username := vars["user"]
  if username == "" {
    return &AppError{nil, "Page not found", 404}
  }

  if ctx.R.Method == "GET" {
    user, exists := GetUserByUsername(username)

    if !exists {
      return &AppError{nil, "User not found", 404}
    }

    d := *ctx.Data
    d["user"] = user

    if ctx.UserSession.LoggedIn() {
      if username == d["_user"].(*sessions.Token).Username {
        d["edit"] = true
      }
    }

    v := ctx.R.URL.Query()
    v.Set("user", user.Username)
    ctx.R.URL.RawQuery = v.Encode()

    switch content {
    case "submitted":
      ListScripts(ctx)
      Render(ctx, "userSubmitted")
    case "requested":
      ListRequests(ctx)
      Render(ctx, "userRequested")
    case "news":
      ListNews(ctx)
      Render(ctx, "userNews")
    case "feedback":
      ListFeedback(ctx)
      Render(ctx, "userFeedback")
    case "comments":
      ListComments(ctx)
      Render(ctx, "userComments")
    default:
      ListScripts(ctx)
      Render(ctx, "user")
    }

  }

  return nil
}

func UserSettingsHandler(ctx *sessions.Context) *AppError {
  if !ctx.UserSession.LoggedIn() {
    ctx.Session.SetFlash(ctx.W, ctx.R, "You are not logged in.")
    http.Redirect(ctx.W, ctx.R, "/", 302)
    return nil
  }

  d := *ctx.Data
  username := d["_user"].(*sessions.Token).Username
  user, exists := GetUserByUsername(username)

  if !exists {
    return &AppError{nil, "User not found", 404}
  }

  d["user"] = user

  if ctx.R.Method == "GET" {
    Render(ctx, "settings")
    return nil
  }

  if ctx.R.Method == "POST" {
    if !parseForm(ctx, "") {
      log.Println("erroring parsing settings form")
      ctx.Session.SetFlash(ctx.W, ctx.R, "Refresh the page and try again")
      http.Redirect(ctx.W, ctx.R, ctx.R.Referer(), 302)
      return nil
    }

    form := ParseSettingsForm(ctx.R.Form)

    user.Name = form.Name
    user.Email = form.Email
    user.Logline = form.Logline
    user.Updated = TimeNowToInt()
    user.Clean()

    // Save user account
    if err := user.Save(); err != nil {
      log.Println("erroring saving user settings")
      ctx.Session.SetFlash(ctx.W, ctx.R, err.Error())
      http.Redirect(ctx.W, ctx.R, ctx.R.Referer(), 302)
      return nil
    }

    ctx.Session.SetFlash(ctx.W, ctx.R, "Settings Updated!")
    http.Redirect(ctx.W, ctx.R, ctx.R.Referer(), 302)
  }

  return nil
}

// Handlers

func LogoutHandler(ctx *sessions.Context) *AppError {
  ctx.UserSession.DeleteAuthToken(ctx)
  ctx.Session.SetFlash(ctx.W, ctx.R, "You have been logged out successfully.")
  http.Redirect(ctx.W, ctx.R, ctx.R.Referer(), 302)
  return nil
}

func LoginHandler(ctx *sessions.Context) *AppError {
  if ctx.UserSession.LoggedIn() {
    ctx.Session.SetFlash(ctx.W, ctx.R, "You're already logged in.")
    http.Redirect(ctx.W, ctx.R, "/", 302)
    return nil
  }

  if ctx.R.Method == "GET" {
    Render(ctx, "loginForm")
    return nil
  }

  if ctx.R.Method == "POST" {
    if !parseForm(ctx, "") {
      loginError(ctx, "Refresh the page and try again.")
      return nil
    }

    form := ParseLoginForm(ctx.R.Form)

    // input is blank or username does not match valid regex.
    if form.Username == "" || form.Password == "" || !ValidUsernameRegex.MatchString(form.Username) {
      loginError(ctx, "username or password was blank")
      return nil
    }

    user, exists := GetUserByUsername(form.Username)
    if !exists {
      loginError(ctx, "Username not in database")
      return nil
    }

    usd, e := GetUserSenseData(int64(user.Id))
    if !e {
      loginError(ctx, "Cant find user sense data...shit")
      return nil
    }

    b := crypt.ValidatePassword(usd.Password, form.Password, usd.Algo)
    if !b {
      loginError(ctx, "invalid username")
      return nil
    }

   remember := false
   if form.Remember == "True" {
     remember = true
   }

    // Successful login
    ctx.Session.SetCSRF(ctx.W, ctx.R)
    ctx.UserSession.SetAuthToken(ctx, user.Username, int64(user.Id), remember)
    http.Redirect(ctx.W, ctx.R, user.Url(), 302)
  }

  return nil
}

func SignupFormHandler(ctx *sessions.Context) *AppError {
  if ctx.UserSession.LoggedIn() {
    ctx.Session.SetFlash(ctx.W, ctx.R, "You're already logged in.")
    http.Redirect(ctx.W, ctx.R, "/", 302)
    return nil
  }

  if ctx.R.Method == "GET" {
    d := *ctx.Data
    d["form"] = new(SignupForm)

    if len(d["info"].([]interface{})) > 0 {
      d["form"] = refillSignupForm(ctx.Session)
    }

    Render(ctx, "signupForm")
    return nil
  }

  if ctx.R.Method == "POST" {
    if !parseForm(ctx, "") {
      signupError(ctx, "Refresh the page and try again.")
      return nil
    }

    form := ParseSignupForm(ctx.R.Form)

    hashedPass, algoSalt := crypt.HashPassword(form.Password)
    if hashedPass == "" || algoSalt == "" {
      signupError(ctx, "There was a problem signing you up, please try again.")
      return nil
    }

    t := TimeNowToInt()

    // Save user
    u := &Users{
      Username: form.Username,
      Email: form.Email,
      Name: form.Name,
      Created: t,
      Updated: t,
    }

    u.Clean()

    // Check if the user exists.
    if exists := u.Exists(); exists {
     signupError(ctx, "The username already exists.")
     return nil
    }

    if err := validatePassword(form.Password); err != nil {
      signupError(ctx, err.Error())
      return nil
    }

    // Save user account
    if err := u.Save(); err != nil {
      signupError(ctx, err.Error())
      return nil
    }

    usd := &UserSenseData {
      UserId: int64(u.Id),
      Password: hashedPass,
      Algo: algoSalt,
      Updated: t,
    }

    // Save user password (if it errors out you are FUCKED)
    if err := usd.Save(); err != nil {
      signupError(ctx, err.Error())
      return nil
    }

    // successful signup
    ctx.Session.SetCSRF(ctx.W, ctx.R)
    ctx.Session.SetFlash(ctx.W, ctx.R, "Thanks for signing up!")
    ctx.Session.ClearForm(ctx.W, ctx.R, "signup")
    ctx.UserSession.SetAuthToken(ctx, u.Username, int64(u.Id), false)
    http.Redirect(ctx.W, ctx.R, u.Url(), 302)
  }

  return nil
}

