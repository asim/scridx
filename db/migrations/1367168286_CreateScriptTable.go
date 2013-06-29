package main

import (
  "github.com/eaigner/hood"
  _ "github.com/ziutek/mymysql/godrv"
  "time"
)

type Scripts struct {
  Id hood.Id
  Title string `sql:"size(255),notnull"`
  Logline string `sql:"size(512)"`
  Writers string `sql:"size(512)"`
  Source string `sql:"size(255),notnull"`
  Imdb string  `sql:"size(255)"`
  Wiki string `sql:"size(255)"`
  Rank float64 `sql:"default(0.0),notnull"`
  Karma int32 `sql:"default(1),notnull"`
  Stored uint32
  Drafted time.Time
  Version string `sql:"size(20)"`
  UserId int64 `sql:"notnull"`
}

func (m *M) CreateScriptTable_1367168286_Up(hd *hood.Hood) {
  hd.CreateTable(&Scripts{})
}

func (m *M) CreateScriptTable_1367168286_Down(hd *hood.Hood) {
  hd.DropTable(&Scripts{})
}
