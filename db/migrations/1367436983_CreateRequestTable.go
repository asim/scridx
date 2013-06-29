package main

import (
  "github.com/eaigner/hood"
  _ "github.com/ziutek/mymysql/godrv"
  "time"
)

type Requests struct {
  Id hood.Id
  Title string `sql:"size(255),notnull"`
  Logline string `sql:"size(512)"`
  Writers string `sql:"size(512)"`
  Source string `sql:"size(255)"`
  Imdb string  `sql:"size(255)"`
  Wiki string `sql:"size(255)"`
  Rank float64 `sql:"default(0.0)"`
  Karma int32 `sql:"default(1)"`
  Stored uint32
  Drafted time.Time
  Version string `sql:"size(20)"`
  Status int8 `sql:"default(0)"`
  UserId int64 `sql:"notnull"`
  FulfillUserId int64
}

func (m *M) CreateRequestTable_1367436983_Up(hd *hood.Hood) {
  hd.CreateTable(&Requests{})
}

func (m *M) CreateRequestTable_1367436983_Down(hd *hood.Hood) {
  hd.DropTable(&Requests{})
}
