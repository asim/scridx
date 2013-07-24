package main

import (
  "github.com/eaigner/hood"
  _ "github.com/ziutek/mymysql/godrv"
)

func (m *M) AddNewsIndex_1374154258_Up(hd *hood.Hood) {
  hd.CreateIndex("news", "sdb_news_title_index", false, "title")
  hd.CreateIndex("news", "sdb_news_source_index", true, "source")
  hd.CreateIndex("news", "sdb_news_rank_index", false, "rank")
}

func (m *M) AddNewsIndex_1374154258_Down(hd *hood.Hood) {
  hd.DropIndex("news", "sdb_news_title_index")
  hd.DropIndex("news", "sdb_news_source_index")
  hd.DropIndex("news", "sdb_news_rank_index")
}
