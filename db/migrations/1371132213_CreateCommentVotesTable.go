package main

import (
  _ "github.com/ziutek/mymysql/godrv"
  "github.com/eaigner/hood"
)

type CommentVotes struct {
  Id       hood.Id
  CommentId int64 `sql:"notnull"`
  UserId   int64 `sql:"notnull"`
}

func (m *M) CreateCommentVotesTable_1371132213_Up(hd *hood.Hood) {
  hd.CreateTable(&CommentVotes{})
}

func (m *M) CreateCommentVotesTable_1371132213_Down(hd *hood.Hood) {
  hd.DropTable(&CommentVotes{})
}
