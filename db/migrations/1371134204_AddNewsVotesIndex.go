package main

import (
    "github.com/eaigner/hood"
  _ "github.com/ziutek/mymysql/godrv"
)

func (m *M) AddNewsVotesIndex_1371134204_Up(hd *hood.Hood) {
  hd.CreateIndex("news_votes", "sdb_news_votes_news_id_user_id_index", true, "news_id", "user_id")
}

func (m *M) AddNewsVotesIndex_1371134204_Down(hd *hood.Hood) {
  hd.DropIndex("news_votes", "sdb_news_votes_news_id_user_id_index")
}
