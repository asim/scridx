package app

import "scridx/sessions"
import "fmt"

type ThreadedComment struct {
  User *Users
  Depth int
  Votable bool
  Comments
}

type Entity struct {
  Index int
  Votable bool
  Thing *interface{}
}

type SetTC struct {
  Comment *ThreadedComment
  Map *map[int64][]*ThreadedComment
  Comments *[]*ThreadedComment
  Users *map[int64]Users
  Votes *map[int64]bool
}

func createEntity(index int, votable bool, i interface{}) *Entity {
  return &Entity{index, votable, &i}
}

//func setTC(t *ThreadedComment, depth int, m *map[int64][]*ThreadedComment, tc *[]*ThreadedComment, u *map[int64]Users) {
func setTC(st *SetTC, depth int, user_id int64) {
  id := int64(st.Comment.Id)
  uid := int64(st.Comment.UserId)
  threadMap := *st.Map
  threadUsers := *st.Users
  st.Comment.Depth = depth
  votes := *st.Votes

  if user, prs := threadUsers[uid]; prs {
    st.Comment.User = &user
  } else {
    st.Comment.User = &Users{
      Username: "-",
      Name: "Anonymous",
    }
  }

  // Votable?
  if user_id == st.Comment.UserId || votes[id] {
    st.Comment.Votable = false
  } else {
    st.Comment.Votable = true
  }

  *st.Comments = append(*st.Comments, st.Comment)
  for _, ct := range threadMap[id] {
    if ct != nil {
//      setTC(ct, depth + 1, m, tc, u)
       std := &SetTC{ct, st.Map, st.Comments, st.Users, st.Votes}
       setTC(std, depth + 1, user_id)
    }
  }
}

func threadedComments(entity_type int16, entity_id int64, parent_id int64, ctx *sessions.Context) []*ThreadedComment {
  q := Comment.Select().Where("entity_type", "=", entity_type).And("entity_id", "=", entity_id).OrderBy("id").Desc()

  var comments []Comments
  Comment.Get(q, &comments)

  var user_ids []string
  var votes *[]CommentVotes
  var user_id int64 = int64(-1)
  var ids []string
  var voted map[int64]bool = make(map[int64]bool)


  // m is map containing key value pair of parent id and child comment pointers tc.
  m := make(map[int64][]*ThreadedComment)

  // Loop the flat comments array, create threaded comment and append to parent id slice.
  for _, comment := range comments {
    ids = append(ids, fmt.Sprintf("%d", comment.Id))

    pid := comment.ParentId
    t := &ThreadedComment{nil, 0, true, comment}
    m[pid] = append(m[pid], t)
    user_ids = append(user_ids, fmt.Sprintf("%d", comment.UserId))
  }

  // If logged in get votes
  if ctx.UserSession.LoggedIn() {
    d := *ctx.Data
    user_id = d["_user"].(*sessions.Token).Id
    // Get votes and create map
    var ok bool
    if votes, ok = GetCommentVotes(user_id, ids); ok {
      for _, vote := range *votes {
        voted[vote.CommentId] = true
      }
    }
  }

  users, _ := GetUsersById(user_ids)

  // Loop through the map starting parent_id of 0 since those are the root comments.
  var tcs []*ThreadedComment
  for _, t  := range m[parent_id] {
    if t != nil {
      st := &SetTC{t, &m, &tcs, users, &voted}
     // setTC(t, 0, &m, &tcs, users)
     setTC(st, 0, user_id)
    }
  }

  return tcs
}

func ListComments(ctx *sessions.Context) {
  vars := ctx.R.URL.Query()

  d := *ctx.Data
  id := int64(-1)
  voted := make(map[int64]bool)
  var ids []string
  var comments []Comments
  var votes *[]CommentVotes

  limit := 20
  page, offset := GetPageOffset(vars, limit)

  q := Comment.Select()

  // Filter by user
  // TODO: needs some error handling
  if username := vars.Get("user"); ValidUsernameRegex.MatchString(username) {
    qq := User.Select("id").Where("username", "=", username).Limit(1)
    var u []Users

    if User.Get(qq, &u); len(u) > 0 {
      q = q.Where("user_id", "=", u[0].Id)
    }
  }

  // Order By
  switch vars.Get("order") {
  case "latest":
    q = q.OrderBy("id").Desc().Limit(limit).Offset(offset)
  case "top":
    q = q.OrderBy("karma").Desc().Limit(limit).Offset(offset)
  default:
    q = q.OrderBy("id").Desc().Limit(limit).Offset(offset)
  }

  Comment.Get(q, &comments)
  pager := GetPager(ctx.R.URL, page, limit, len(comments))

  var user_ids []string
  var tcs []*ThreadedComment

  for _, comment := range comments {
    ids = append(ids, fmt.Sprintf("%d", comment.Id))

    tc := &ThreadedComment{nil, 0, true, comment}
    tcs = append(tcs, tc)
    user_ids = append(user_ids, fmt.Sprintf("%d", comment.UserId))
  }

  // If logged in get votes
  if ctx.UserSession.LoggedIn() {
    id = d["_user"].(*sessions.Token).Id
    // Get votes and create map
    var ok bool
    if votes, ok = GetCommentVotes(id, ids); ok {
      for _, vote := range *votes {
        voted[vote.CommentId] = true
      }
    }
  }

  users, ok := GetUsersById(user_ids)

  for _, t := range tcs {
    if ok {
      us := *users
      uid := int64(t.UserId)
      if user, prs := us[uid]; prs {
        t.User = &user
      } else {
        t.User = &Users{
          Username: "-",
          Name: "Anonymous",
        }
      }
    }

    // Votable?
    if id == t.UserId || voted[int64(t.Id)] {
      t.Votable = false
    } else {
      t.Votable = true
    }
  }

  d["comments"] = &tcs
  d["pager"] = &pager
}

func ListFeedback(ctx *sessions.Context) {
  vars := ctx.R.URL.Query()

  d := *ctx.Data
  id := int64(-1)
  votable := true
  voted := make(map[int64]bool)
  var ids []string
  var feedback []Feedback
  var votes *[]FeedbackVotes

  limit := 20
  page, offset := GetPageOffset(vars, limit)

  q := ListQuery(FB.Select(), vars.Get("order"), vars.Get("user"), limit, offset)
  FB.Get(q, &feedback)

  // Create list of request ids
  for _, item := range feedback {
    ids = append(ids, fmt.Sprintf("%d", item.Id))
  }

  // If logged in get votes
  if ctx.UserSession.LoggedIn() {
    id = d["_user"].(*sessions.Token).Id
    // Get votes and create map
    var ok bool
    if votes, ok = GetFeedbackVotes(id, ids); ok {
      for _, vote := range *votes {
        voted[vote.FeedbackId] = true
      }
    }
  }

  // Create Entities which will be displayed on page
  entities := make([]*Entity, len(feedback))
  for i := 0; i < len(feedback); i++ {
    // Votable?
    if id == feedback[i].UserId || voted[int64(feedback[i].Id)] {
      votable = false
    } else {
      votable = true
    }

    entities[i] = createEntity(i+1+offset, votable, feedback[i])
  }

  pager := GetPager(ctx.R.URL, page, limit, len(feedback))

  d["feedback"] = &entities
  d["pager"] = &pager
}

func ListNews(ctx *sessions.Context) {
  vars := ctx.R.URL.Query()

  d := *ctx.Data
  id := int64(-1)
  votable := true
  voted := make(map[int64]bool)
  var ids []string
  var news []News
  var votes *[]NewsVotes

  limit := 20
  page, offset := GetPageOffset(vars, limit)

  q := ListQuery(NewsItem.Select(), vars.Get("order"), vars.Get("user"), limit, offset)
  NewsItem.Get(q, &news)

  // Create list of request ids
  for _, item := range news {
    ids = append(ids, fmt.Sprintf("%d", item.Id))
  }

  // If logged in get votes
  if ctx.UserSession.LoggedIn() {
    id = d["_user"].(*sessions.Token).Id

    // Get votes and create map
    var ok bool
    if votes, ok = GetNewsVotes(id, ids); ok {
      for _, vote := range *votes {
        voted[vote.NewsId] = true
      }
    }
  }

  // Create Entities which will be displayed on page
  entities := make([]*Entity, len(news))
  for i := 0; i < len(news); i++ {
    // Votable?
    if id == news[i].UserId || voted[int64(news[i].Id)] {
      votable = false
    } else {
      votable = true
    }

    entities[i] = createEntity(i+1+offset, votable, news[i])
  }

  pager := GetPager(ctx.R.URL, page, limit, len(news))

  d["news"] = &entities
  d["pager"] = &pager
}

func ListRequests(ctx *sessions.Context) {
  vars := ctx.R.URL.Query()

  d := *ctx.Data
  id := int64(-1)
  votable := true
  voted := make(map[int64]bool)
  var ids []string
  var requests []Requests
  var votes *[]RequestVotes

  limit := 20
  page, offset := GetPageOffset(vars, limit)

  q := ListQuery(Request.Select(), vars.Get("order"), vars.Get("user"), limit, offset)
  Request.Get(q, &requests)

  // Create list of request ids
  for _, request := range requests {
    ids = append(ids, fmt.Sprintf("%d", request.Id))
  }

  // If logged in get votes
  if ctx.UserSession.LoggedIn() {
    id = d["_user"].(*sessions.Token).Id
    // Get votes and create map

    var ok bool
    if votes, ok = GetRequestVotes(id, ids); ok {
      for _, vote := range *votes {
        voted[vote.RequestId] = true
      }
    }
  }

  // Create Entities which will be displayed on page
  entities := make([]*Entity, len(requests))
  for i := 0; i < len(requests); i++ {
    // Votable?
    if id == requests[i].UserId || voted[int64(requests[i].Id)] {
      votable = false
    } else {
      votable = true
    }
    entities[i] = createEntity(i+1+offset, votable, requests[i])
  }

  pager := GetPager(ctx.R.URL, page, limit, len(requests))

  d["requests"] = &entities
  d["pager"] = &pager
}

func ListScripts(ctx *sessions.Context) {
  vars := ctx.R.URL.Query()

  d := *ctx.Data
  id := int64(-1)
  votable := true
  voted := make(map[int64]bool)
  var ids []string
  var scripts []Scripts
  var votes *[]ScriptVotes

  limit := 20
  page, offset := GetPageOffset(vars, limit)

  q := ListQuery(Script.Select(), vars.Get("order"), vars.Get("user"), limit, offset)
  Script.Get(q, &scripts)

  // Create list of script ids
  for _, script := range scripts {
    ids = append(ids, fmt.Sprintf("%d", script.Id))
  }

  // If logged in get votes
  if ctx.UserSession.LoggedIn() {
    id = d["_user"].(*sessions.Token).Id
    // Get votes and create map

    var ok bool
    if votes, ok = GetScriptVotes(id, ids); ok {
     for _, vote := range *votes {
       voted[vote.ScriptId] = true
     }
    }
  }

  // Create Entities which will be displayed on page
  entities := make([]*Entity, len(scripts))
  for i := 0; i < len(scripts); i++ {
    // Votable?
    if id == scripts[i].UserId || voted[int64(scripts[i].Id)] {
      votable = false
    } else {
      votable = true
    }

    entities[i] = createEntity(i+1+offset, votable, scripts[i])
  }

  // Pager
  pager := GetPager(ctx.R.URL, page, limit, len(scripts))

  d["scripts"] = &entities
  d["pager"] = &pager
}
