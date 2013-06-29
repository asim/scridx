package db

import (
	"github.com/eaigner/hood"
	"time"
)

type Scripts struct {
	Id      hood.Id
	Title   string  `sql:"size(255),notnull"`
	Logline string  `sql:"size(512)"`
	Writers string  `sql:"size(512)"`
	Source  string  `sql:"size(255),notnull"`
	Imdb    string  `sql:"size(255)"`
	Wiki    string  `sql:"size(255)"`
	Rank    float64 `sql:"default(0.0),notnull"`
	Karma   int32   `sql:"default(1),notnull"`
	Stored  uint32
	Drafted time.Time
	Version string `sql:"size(20)"`
	UserId  int64  `sql:"notnull"`
}

func (table *Scripts) Indexes(indexes *hood.Indexes) {
	indexes.Add("sdb_scripts_title_index", "title")
	indexes.Add("sdb_scripts_stored_index", "stored")
	indexes.Add("sdb_scripts_drafted_index", "drafted")
	indexes.AddUnique("sdb_scripts_source_index", "source")
}

type Requests struct {
	Id            hood.Id
	Title         string  `sql:"size(255),notnull"`
	Logline       string  `sql:"size(512)"`
	Writers       string  `sql:"size(512)"`
	Source        string  `sql:"size(255)"`
	Imdb          string  `sql:"size(255)"`
	Wiki          string  `sql:"size(255)"`
	Rank          float64 `sql:"default(0.0)"`
	Karma         int32   `sql:"default(1)"`
	Stored        uint32
	Drafted       time.Time
	Version       string `sql:"size(20)"`
	Status        int8   `sql:"default(0)"`
	UserId        int64  `sql:"notnull"`
	FulfillUserId int64
}

type Users struct {
	Id            hood.Id
	Username      string `sql:"size(15),notnull"`
	Email         string `sql:"size(255),notnull"`
	Name          string `sql:"size(20),notnull"`
	Logline       string `sql:"size(160)"`
	Karma         int32  `sql:"default(0)"`
	Updated       uint32
	Created       uint32
	LoginAttempts int  `sql:"default(0),notnull"`
	Locked        bool `sql:"default(0),notnull"`
}

func (table *Users) Indexes(indexes *hood.Indexes) {
	indexes.AddUnique("sdb_users_users_username_index", "username")
}

type UserSenseData struct {
	Id       hood.Id
	UserId   int64
	Password string `sql:"size(255),notnull"`
	Algo     string `sql:"size(128),notnull"`
	Updated  uint32
}

func (table *UserSenseData) Indexes(indexes *hood.Indexes) {
	indexes.AddUnique("sdb_users_user_sense_data_user_id_index", "user_id")
}

type News struct {
	Id     hood.Id
	Title  string  `sql:"size(255),notnull"`
	Source string  `sql:"size(255)"`
	Rank   float64 `sql:"default(0.0),notnull"`
	Karma  int32   `sql:"default(1),notnull"`
	Stored uint32
	UserId int64  `sql:"notnull"`
	Text   string `sql:"size(65535)"`
}

type Feedback struct {
	Id     hood.Id
	Title  string  `sql:"size(255),notnull"`
	Source string  `sql:"size(255)"`
	Rank   float64 `sql:"default(0.0),notnull"`
	Karma  int32   `sql:"default(1),notnull"`
	Stored uint32  `sql:"notnull"`
	UserId int64   `sql:"notnull"`
	Text   string  `sql:"size(65535)"`
}

type Comments struct {
	Id         hood.Id
	Text       string  `sql:"size(65535),notnull"`
	Rank       float64 `sql:"default(0.0),notnull"`
	Karma      int32   `sql:"default(1),notnull"`
	Stored     uint32  `sql:"notnull"`
	UserId     int64   `sql:"notnull"`
	ParentId   int64   `sql:"notnull"`
	EntityId   int64   `sql:"notnull"`
	EntityType int16   `sql:"notnull"`
}

type ScriptVotes struct {
	Id       hood.Id
	ScriptId int64 `sql:"notnull"`
	UserId   int64 `sql:"notnull"`
}

func (table *ScriptVotes) Indexes(indexes *hood.Indexes) {
	indexes.AddUnique("sdb_script_votes_script_id_user_id_index", "script_id", "user_id")
}

type RequestVotes struct {
	Id        hood.Id
	RequestId int64 `sql:"notnull"`
	UserId    int64 `sql:"notnull"`
}

func (table *RequestVotes) Indexes(indexes *hood.Indexes) {
	indexes.AddUnique("sdb_request_votes_request_id_user_id_index", "request_id", "user_id")
}

type FeedbackVotes struct {
	Id         hood.Id
	FeedbackId int64 `sql:"notnull"`
	UserId     int64 `sql:"notnull"`
}

func (table *FeedbackVotes) Indexes(indexes *hood.Indexes) {
	indexes.AddUnique("sdb_feedback_votes_feedback_id_user_id_index", "feedback_id", "user_id")
}

type NewsVotes struct {
	Id     hood.Id
	NewsId int64 `sql:"notnull"`
	UserId int64 `sql:"notnull"`
}

func (table *NewsVotes) Indexes(indexes *hood.Indexes) {
	indexes.AddUnique("sdb_news_votes_news_id_user_id_index", "news_id", "user_id")
}

type CommentVotes struct {
	Id        hood.Id
	CommentId int64 `sql:"notnull"`
	UserId    int64 `sql:"notnull"`
}

func (table *CommentVotes) Indexes(indexes *hood.Indexes) {
	indexes.AddUnique("sdb_comment_votes_comment_id_user_id_index", "comment_id", "user_id")
}
