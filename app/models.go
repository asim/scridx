package app

import (
  "log"
    "github.com/garyburd/redigo/redis"
  "github.com/eaigner/hood"
  _ "github.com/ziutek/mymysql/godrv"
  "scridx/db"
  "strings"
  "fmt"
)

const (
  // Request status
  statusNew = 0
  statusPending = 1
  statusCompleted = 2

  // Entity types for comments
  ScriptType = 1
  RequestType = 2
  NewsType = 3
  FeedbackType = 4
  CommentType = 5
)

// Entity Tables
type Scripts db.Scripts
type Requests db.Requests
type Users db.Users
type News db.News
type Feedback db.Feedback
type Comments db.Comments
// Vote Tables
type ScriptVotes db.ScriptVotes
type RequestVotes db.RequestVotes
type NewsVotes db.NewsVotes
type FeedbackVotes db.FeedbackVotes
type CommentVotes db.CommentVotes

// User sensitive data table
type UserSenseData db.UserSenseData

var DB *hood.Hood
var Redis redis.Conn

var rankingAlgorithm = "(karma + %d) - 1 / ((stored - %d / 60) / 60 + 2)^1.5"

func initRedis() error {
  var err error
  Redis, err = redis.Dial("tcp", ":6379")
  if err != nil {
    return err
  }

  return nil
}

func InitDB(config, env string) (*hood.Hood, error) {
  // hood connection
  var err error
  DB, err = hood.Load(config, env)
  if err != nil {
    return nil, err
  }

  // redis connection
//  err = initRedis()
//  if err != nil {
//    return nil, err
//  }

  // Init worker
 // Worker()
  Indexer()

  return DB, nil
}

var Script *Scripts
var Request *Requests
var User *Users
var NewsItem *News
var FB *Feedback
var Comment *Comments

var ScriptVote *ScriptVotes
var RequestVote *RequestVotes
var NewsVote *NewsVotes
var FBVote *FeedbackVotes
var CommentVote *CommentVotes

func rollback(hd *hood.Hood, err error) error {
  hd.Rollback()
  log.Println(err)
  return err
}

func cacheSave(i *LRUItem) {
  go func() {
    LRUChan <- i
  }()
}

func save(i interface{}) (hood.Id, error) {
  err := DB.Validate(i)
  if err != nil {
    log.Println(err)
    return -1, err
  }

  var id hood.Id
  id, err = DB.Save(i)
  if err != nil {
    log.Println(err)
    log.Println(err.Error())
    return -1, err
  }
  return id, nil
}

func vote(table string, incr int, id hood.Id, user int64, v interface{}) error {
  // New Transaction
  hd := DB.Begin()

  // Save Vote
  if _, err := hd.Save(v); err != nil {
    return rollback(hd, err)
  }

  // Set Item Karma
  if err := setKarmaAndRank(hd, table, incr, int64(id)); err != nil {
    return rollback(hd, err)
  }

  // Set User Karma
  if err := setKarma(hd, "users", incr, user); err != nil {
    return rollback(hd, err)
  }

  // Commit
  if err := hd.Commit(); err != nil {
    return rollback(hd, err)
  }

  return nil
}

// Must be executed in transaction along with all the other shit
func setKarmaAndRank(hd *hood.Hood, table string, incr int, id int64) error {
  ra := fmt.Sprintf(rankingAlgorithm, incr, TimeNowToInt())
  query := fmt.Sprintf("UPDATE %s SET rank = %s, karma = karma + %d where id = %d LIMIT 1;", table, ra, incr, id)

  if _, err := hd.Exec(query); err != nil {
    return err
  }

  return nil
}

func setKarma(hd *hood.Hood, table string, incr int, id int64) error {
  query := fmt.Sprintf("UPDATE %s SET karma = karma + %d where id = %d LIMIT 1;", table, incr, id)

  if _, err := hd.Exec(query); err != nil {
    return err
  }

  return nil
}

func ListQuery(q *hood.Hood, order, username string, limit, offset int) *hood.Hood {
  // Filter by user
  // TODO: needs some error handling
  if ValidUsernameRegex.MatchString(username) {
    qq := User.Select("id").Where("username", "=", username).Limit(1)
    var u []Users

    if User.Get(qq, &u); len(u) > 0 {
      q = q.Where("user_id", "=", u[0].Id)
    }
  }

  switch order {
  case "latest":
    q = q.OrderBy("id").Desc()
  case "top":
    q = q.OrderBy("rank` desc,`id").Desc()
  case "a-z":
    q = q.OrderBy("title")
  default:
    q = q.OrderBy("rank` desc,`id").Desc()
  }

  return q.Limit(limit).Offset(offset)
}

// Database methods
func (c *Comments) Select(columns ...hood.Path) *hood.Hood {
  return DB.Select("comments", columns...)
}

func (c *Comments) Get(q *hood.Hood, comments *[]Comments) bool {
  err := q.Find(comments)
  if err != nil {
    log.Println(err)
    return false
  }

  return true
}

func (n *News) Select(columns ...hood.Path) *hood.Hood {
  return DB.Select("news", columns...)
}

func (n *News) Get(q *hood.Hood, news *[]News) bool {
  err := q.Find(news)
  if err != nil {
    log.Println(err)
    return false
  }

  return true
}

func (f *Feedback) Select(columns ...hood.Path) *hood.Hood {
  return DB.Select("feedback", columns...)
}

func (f *Feedback) Get(q *hood.Hood, feedback *[]Feedback) bool {
  err := q.Find(feedback)
  if err != nil {
    log.Println(err)
    return false
  }

  return true
}

func (r *Requests) Select(columns ...hood.Path) *hood.Hood {
  return DB.Select("requests", columns...)
}

func (r *Requests) Get(q *hood.Hood, requests *[]Requests) bool {
  err := q.Find(requests)
  if err != nil {
    log.Println(err)
    return false
  }

  return true
}

func (s *Scripts) Select(columns ...hood.Path) *hood.Hood {
  return DB.Select("scripts", columns...)
}

func (s *Scripts) Get(q *hood.Hood, scripts *[]Scripts) bool {
  err := q.Find(scripts)
  if err != nil {
    log.Println(err)
    return false
  }

  return true
}

func (u *Users) Select(columns ...hood.Path) *hood.Hood {
  return DB.Select("users", columns...)
}

func (u *Users) Get(q *hood.Hood, users *[]Users) bool {
  err := q.Find(users)
  if err != nil {
    log.Println(err)
    return false
  }

  return true
}

func (v *ScriptVotes) Select(columns ...hood.Path) *hood.Hood {
  return DB.Select("script_votes", columns...)
}

func (v *ScriptVotes) Get(q *hood.Hood, votes *[]ScriptVotes) bool {
  err := q.Find(votes)
  if err != nil {
    log.Println(err)
    return false
  }

  return true
}

func (v *RequestVotes) Select(columns ...hood.Path) *hood.Hood {
  return DB.Select("request_votes", columns...)
}

func (v *RequestVotes) Get(q *hood.Hood, votes *[]RequestVotes) bool {
  err := q.Find(votes)
  if err != nil {
    log.Println(err)
    return false
  }

  return true
}

func (v *NewsVotes) Select(columns ...hood.Path) *hood.Hood {
  return DB.Select("news_votes", columns...)
}

func (v *NewsVotes) Get(q *hood.Hood, votes *[]NewsVotes) bool {
  err := q.Find(votes)
  if err != nil {
    log.Println(err)
    return false
  }

  return true
}

func (v *FeedbackVotes) Select(columns ...hood.Path) *hood.Hood {
  return DB.Select("feedback_votes", columns...)
}

func (v *FeedbackVotes) Get(q *hood.Hood, votes *[]FeedbackVotes) bool {
  err := q.Find(votes)
  if err != nil {
    log.Println(err)
    return false
  }

  return true
}

func (v *CommentVotes) Select(columns ...hood.Path) *hood.Hood {
  return DB.Select("comment_votes", columns...)
}

func (v *CommentVotes) Get(q *hood.Hood, votes *[]CommentVotes) bool {
  err := q.Find(votes)
  if err != nil {
    log.Println(err)
    return false
  }

  return true
}

func GetScriptVotes(user_id int64, ids []string) (*[]ScriptVotes, bool) {
  if len(ids) == 0 {
    return nil, false
  }

  var results []ScriptVotes
  var query = fmt.Sprintf("SELECT * from script_votes WHERE script_id IN (%s) AND user_id = %d", strings.Join(ids, ","), user_id)

  err := DB.FindSql(&results, query)
  if err != nil {
    log.Println(err)
    log.Println(err.Error())
    return nil, false
  }

  return &results, true
}

func GetRequestVotes(user_id int64, ids []string) (*[]RequestVotes, bool) {
  if len(ids) == 0 {
    return nil, false
  }

  var results []RequestVotes
  var query = fmt.Sprintf("SELECT * from request_votes WHERE request_id IN (%s) AND user_id = %d", strings.Join(ids, ","), user_id)

  err := DB.FindSql(&results, query)
  if err != nil {
    log.Println(err)
    log.Println(err.Error())
    return nil, false
  }

  return &results, true
}

func GetNewsVotes(user_id int64, ids []string) (*[]NewsVotes, bool) {
  if len(ids) == 0 {
    return nil, false
  }

  var results []NewsVotes
  var query = fmt.Sprintf("SELECT * from news_votes WHERE news_id IN (%s) AND user_id = %d", strings.Join(ids, ","), user_id)

  err := DB.FindSql(&results, query)
  if err != nil {
    log.Println(err)
    log.Println(err.Error())
    return nil, false
  }

  return &results, true
}

func GetFeedbackVotes(user_id int64, ids []string) (*[]FeedbackVotes, bool) {
  if len(ids) == 0 {
    return nil, false
  }

  var results []FeedbackVotes
  var query = fmt.Sprintf("SELECT * from feedback_votes WHERE feedback_id IN (%s) AND user_id = %d", strings.Join(ids, ","), user_id)

  err := DB.FindSql(&results, query)
  if err != nil {
    log.Println(err)
    log.Println(err.Error())
    return nil, false
  }

  return &results, true
}

func GetCommentVotes(user_id int64, ids []string) (*[]CommentVotes, bool) {
  if len(ids) == 0 {
    return nil, false
  }

  var results []CommentVotes
  var query = fmt.Sprintf("SELECT * from comment_votes WHERE comment_id IN (%s) AND user_id = %d", strings.Join(ids, ","), user_id)

  err := DB.FindSql(&results, query)
  if err != nil {
    log.Println(err)
    log.Println(err.Error())
    return nil, false
  }

  return &results, true
}

func GetLatestRequests(limit int, offset int) ([]Requests, bool){
  var requests []Requests
  err := DB.Select("requests").OrderBy("id").Desc().Limit(limit).Offset(offset).Find(&requests)
  if err != nil {
    log.Println(err)
    return requests, false
  }

  return requests, true
}

func GetLatestScripts(limit int, offset int) ([]Scripts, bool) {
  var scripts []Scripts
  err := DB.Select("scripts").OrderBy("id").Desc().Limit(limit).Offset(offset).Find(&scripts)
  if err != nil {
    log.Println(err)
    return scripts, false
  }

  return scripts, true
}

func GetLRUItems(limit int, offset int) ([]*LRUItem, bool) {
  c := make(chan *LRU, 1)
  LRURefChan <- c // pass sender chan for lru ref
  ref := <-c // receive pointer to lru
  data := *ref
  size := len(data)

  // Offset too far
  length := size - offset
  if size == 0 || length <= 0 || offset < 0 {
    return nil, false
  }

  if length < limit {
    limit = length
  }

  lru := make([]*LRUItem, limit)
  for i := 0; i < limit; i++ {
    lru[i] = data[length - (1+i)]
  }
  return lru, true
}

func GetRequest(id string) (*Requests, bool) {
  var results []Requests
  err := DB.Where("id", "=", id).Limit(1).Find(&results)
  if err != nil {
    log.Println(err)
    log.Println(err.Error())
    return nil, false
  }

  if len(results) == 0 {
    return nil, false
  }

  s := results[0]
  return &s, true
}

func GetComment(id string) (*Comments, bool) {
  var results []Comments
  err := DB.Where("id", "=", id).Limit(1).Find(&results)
  if err != nil {
    log.Println(err)
    log.Println(err.Error())
    return nil, false
  }

  if len(results) == 0 {
    return nil, false
  }

  s := results[0]
  return &s, true
}

func GetFeedback(id string) (*Feedback, bool) {
  var results []Feedback
  err := DB.Where("id", "=", id).Limit(1).Find(&results)
  if err != nil {
    log.Println(err)
    log.Println(err.Error())
    return nil, false
  }

  if len(results) == 0 {
    return nil, false
  }

  s := results[0]
  return &s, true
}

func GetNews(id string) (*News, bool) {
  var results []News
  err := DB.Where("id", "=", id).Limit(1).Find(&results)
  if err != nil {
    log.Println(err)
    log.Println(err.Error())
    return nil, false
  }

  if len(results) == 0 {
    return nil, false
  }

  s := results[0]
  return &s, true
}

func GetScript(id string) (*Scripts, bool) {
  var results []Scripts
  err := DB.Where("id", "=", id).Limit(1).Find(&results)
  if err != nil {
    log.Println(err)
    log.Println(err.Error())
    return nil, false
  }

  if len(results) == 0 {
    return nil, false
  }

  s := results[0]
  return &s, true
}

func (s *Scripts) CheckIfSourceExists() (*Scripts, bool) {
  var results []Scripts
  err := DB.Where("source", "=", s.Source).Limit(1).Find(&results)
  if err != nil {
    log.Println(err)
    log.Println(err.Error())
    return nil, false
  }

  if len(results) == 0 {
    return nil, false
  }

  r := results[0]
  return &r, true
}

func (u *Users) Exists() bool {
  var results []Users
  err := DB.Where("username", "=", u.Username).Limit(1).Find(&results)
  if err != nil {
    log.Println(err)
    log.Println(err.Error())
    return false
  }

  if len(results) == 0 {
    return false
  }

  return true
}

func GetUserByUsername(username string) (*Users, bool) {
  var results []Users
  err := DB.Where("username", "=", username).Limit(1).Find(&results)
  if err != nil {
    log.Println(err)
    log.Println(err.Error())
    return nil, false
  }

  if len(results) == 0 {
    return nil, false
  }

  u := results[0]
  return &u, true
}

func GetUsersById(ids []string) (*map[int64]Users, bool) {
  if len(ids) == 0 {
    return nil, false
  }

  var results []Users
  var query = fmt.Sprintf("SELECT id, username, name from users WHERE id IN (%s)", strings.Join(ids, ","))

  err := DB.FindSql(&results, query)
  if err != nil {
    log.Println(err)
    log.Println(err.Error())
    return nil, false
  }

  users := make(map[int64]Users)

  for _, user := range results {
    users[int64(user.Id)] = user
  }

  return &users, true
}

func GetUserById(id int64) (*Users, bool) {
  var results []Users
  err := DB.Where("id", "=", id).Limit(1).Find(&results)
  if err != nil {
    log.Println(err)
    log.Println(err.Error())
    return nil, false
  }

  if len(results) == 0 {
    return nil, false
  }

  u := results[0]
  return &u, true
}

func GetUserSenseData(id int64) (*UserSenseData, bool) {
  var results []UserSenseData
  err := DB.Where("id", "=", id).Limit(1).Find(&results)
  if err != nil {
    log.Println(err)
    log.Println(err.Error())
    return nil, false
  }

  if len(results) == 0 {
    return nil, false
  }

  us := results[0]
  return &us, true
}

func (c *Comments) Save() error {
  _, err := save(c)
  return err
}

func (f *Feedback) Save() error {
//  votes := make(map[int64]bool)
//  i := &LRUItem{0, true, "feedback", f, votes}
//  cacheSave(i)
  _, err := save(f)
  return err
}

func (n *News) Save() error {
//  votes := make(map[int64]bool)
//  i := &LRUItem{0, true, "news", n, votes}
//  cacheSave(i)
  _, err := save(n)
  return err
}

func (r *Requests) Save() error {
//  votes := make(map[int64]bool)
//  i := &LRUItem{0, true, "request", r, votes}
//  cacheSave(i)
  _, err := save(r)
  return err
}

func (s *Scripts) Save() error {
//  votes := make(map[int64]bool)
//  i := &LRUItem{0, true, "script", s, votes}
//  cacheSave(i)
  id, err := save(s)
  if err == nil && id > 0 {
    s.Id = id
    Index(s)
  }

  return err
}

func (u *Users) Save() error {
  _, err := save(u)
  return err
}

func (us *UserSenseData) Save() error {
  _, err := save(us)
  return err
}

// Voting - occurs in transaction
func (s *Scripts) Vote(user_id int64, incr int) error {
  v := &ScriptVotes{
    ScriptId: int64(s.Id),
    UserId: user_id,
  }

  return vote("scripts", 1, s.Id, s.UserId, v)
}

func (r *Requests) Vote(user_id int64, incr int) error {
  v := &RequestVotes{
    RequestId: int64(r.Id),
    UserId: user_id,
  }

  return vote("requests", 1, r.Id, r.UserId, v)
}

func (n *News) Vote(user_id int64, incr int) error {
  v := &NewsVotes{
    NewsId: int64(n.Id),
    UserId: user_id,
  }

  return vote("news", 1, n.Id, n.UserId, v)
}

func (f *Feedback) Vote(user_id int64, incr int) error {
  v := &FeedbackVotes{
    FeedbackId: int64(f.Id),
    UserId: user_id,
  }

  return vote("feedback", 1, f.Id, f.UserId, v)
}

func (c *Comments) Vote(user_id int64, incr int) error {
  v := &CommentVotes{
    CommentId: int64(c.Id),
    UserId: user_id,
  }

  return vote("comments", incr, c.Id, c.UserId, v)
}

