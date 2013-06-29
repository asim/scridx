package app

import (
  "github.com/gorilla/mux"
  "net/http"
)

func RegisterRoutes() *mux.Router {
  r := mux.NewRouter()
  r.Handle("/", Handler(IndexHandler))

  // Requests
  r.StrictSlash(true).Handle("/requests", Handler(ListHandler))
  r.StrictSlash(true).Handle("/requests/{id:[0-9]+}", Handler(RequestHandler))
  r.StrictSlash(true).Handle("/requests/{id:[0-9]+}/{title}", Handler(RequestHandler))

  // Scripts
  r.StrictSlash(true).Handle("/scripts", Handler(ListHandler))
  r.StrictSlash(true).Handle("/scripts/{id:[0-9]+}", Handler(ScriptHandler))
  r.StrictSlash(true).Handle("/scripts/{id:[0-9]+}/{title}", Handler(ScriptHandler))

  // News
  r.StrictSlash(true).Handle("/news", Handler(ListHandler))
  r.StrictSlash(true).Handle("/news/{id:[0-9]+}", Handler(NewsHandler))
  r.StrictSlash(true).Handle("/news/{id:[0-9]+}/{title}", Handler(NewsHandler))

  // News
  r.StrictSlash(true).Handle("/feedback", Handler(ListHandler))
  r.StrictSlash(true).Handle("/feedback/{id:[0-9]+}", Handler(FeedbackHandler))
  r.StrictSlash(true).Handle("/feedback/{id:[0-9]+}/{title}", Handler(FeedbackHandler))

  // Comments
  r.StrictSlash(true).Handle("/comments", Handler(ListHandler))
  r.StrictSlash(true).Handle("/comments/{id:[0-9]+}", Handler(CommentHandler))

  // Forms (must be under /submit to authenticate login)
  r.Handle("/submit/{form}", Handler(FormHandler))
  r.Handle("/submit/{form}/r/{id:[0-9]+}/{title}", Handler(FormHandler))

  // Votes handler
  r.Handle("/vote/{entity}/{id:[0-9]+}", Handler(VoteHandler))
  r.Handle("/vote/{entity}/{id:[0-9]+}/{dir}", Handler(VoteHandler))

  // User
  r.StrictSlash(true).Handle("/u/{user}", Handler(UserHandler))
  r.Handle("/u/{user}/{content}", Handler(UserHandler))
  r.Handle("/settings", Handler(UserSettingsHandler))

  // Authentication
  r.Handle("/signup", Handler(SignupFormHandler))
  r.Handle("/login", Handler(LoginHandler))
  r.Handle("/logout", Handler(LogoutHandler))

  // Static
  r.Handle("/favicon.ico", http.HandlerFunc(FaviconHandler))
  r.Handle("/robots.txt", http.HandlerFunc(RobotHandler))
  r.PathPrefix("/static/").Handler(http.FileServer(http.Dir("app")))

  // Not found
  r.NotFoundHandler = Handler(NotFoundHandler)

  return r
}
