package main

import (
  "github.com/eaigner/hood"
  _ "github.com/ziutek/mymysql/godrv"
)

type Feedback struct {
  Id hood.Id
  Title string `sql:"size(255),notnull"`
  Source string `sql:"size(255)"`
  Rank float64 `sql:"default(0.0),notnull"`
  Karma int32 `sql:"default(1),notnull"`
  Stored uint32 `sql:"notnull"`
  UserId int64 `sql:"notnull"`
  Text string `sql:"size(65535)"`
}

func (m *M) CreateFeedbackTable_1370371087_Up(hd *hood.Hood) {
  hd.CreateTable(&Feedback{})
	// TODO: implement
}

func (m *M) CreateFeedbackTable_1370371087_Down(hd *hood.Hood) {
  hd.DropTable(&Feedback{})
	// TODO: implement
}
