package main

import (
  _ "github.com/ziutek/mymysql/godrv"
  "github.com/eaigner/hood"
)

type FeedbackVotes struct {
  Id       hood.Id
  FeedbackId int64 `sql:"notnull"`
  UserId   int64 `sql:"notnull"`
}

func (m *M) CreateFeedbackVotesTable_1371132200_Up(hd *hood.Hood) {
  hd.CreateTable(&FeedbackVotes{})
}

func (m *M) CreateFeedbackVotesTable_1371132200_Down(hd *hood.Hood) {
  hd.DropTable(&FeedbackVotes{})
}
