package app

import (
  "github.com/mattbaird/elastigo/core"
  "log"
  "time"
)

// Script
type IndexItem struct {
  Id int64
  Type string `json:"type"`
  Title string `json:"title"`
  Writers string `json:"writers"`
  Logline string `json:logline"`
  Drafted time.Time `json:"drafted"`
  Username string `json:"user"`
  Source string `json:"source"`
  Url string `json:"url"`
  Wiki string `json:"wiki"`
  Imdb string `json:"imdb"`
}

var IndexChan chan *IndexItem

func Index(script *Scripts) {
  username := "unknown"
  if u, exists := GetUserById(script.UserId); exists {
    username = u.Username
  }

  i := &IndexItem{
    Id: int64(script.Id),
    Type: "script",
    Title: script.Title,
    Writers: script.Writers,
    Logline: script.Logline,
    Drafted: script.Drafted,
    Username: username,
    Source: script.Source,
    Url: script.Url(),
    Wiki: script.Wiki,
    Imdb: script.Imdb,
  }

  IndexChan <- i
}

func Indexer() {
  log.Println("start indexer")

  IndexChan = make(chan *IndexItem, 100)

  go func() {
    for {
      // Get item to be indexed, blocking receive
      item := <-IndexChan

      // Define item type
      t := item.Type
      if t == "" {
        t = "misc"
      }

      // Index Item
      r, err := core.Index(true, "scridx", t, "", item)
      if err != nil {
        log.Println("error indexing:", err)
      } else {
        log.Println("indexed item (id, type):", r.Id, r.Type)
      }
    }
  }()
}
