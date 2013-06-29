package app
import (
//  "encoding/json"
  "log"
)

type LRUItem struct {
  Index int
  Votable bool
  Type string
  Thing interface{}
  Votes map[int64]bool
}

/*type VoteItem struct {
  Item Item
  Points float64
  Karma  int32
  Vote int
}
*/

// List of latest entities
type LRU []*LRUItem

// Channel to add entities to the LRU
var LRUChan chan *LRUItem
var LRURefChan chan chan<- *LRU

func populateLRU(lru *LRU) {
  r, err := Redis.Do("GET", "lru")
  if err == nil && r != nil {
    log.Println("populating lru from redis")
    *lru = r.(LRU)
    return
  }

  log.Println("populating lru from databse")
  var scripts []Scripts
  q := Script.Select().OrderBy("id").Desc().Limit(500)
  Script.Get(q, &scripts)
  if size := len(scripts); size > 0 {
    for i := 0; i < size; i++ {
      s := size - i
      var svotes []ScriptVotes
      q := ScriptVote.Select("user_id").Where("script_id", "=", int64(scripts[s-1].Id))
      ScriptVote.Get(q, &svotes)
      votes := make(map[int64]bool)
      for _, vote := range svotes {
        votes[vote.UserId] = true
      }

      *lru = append(*lru, &LRUItem{s, true, "script", scripts[s-1], votes})
    }
  }
}

func saveLRU(lru *LRU) {
  _, err := Redis.Do("SET", "lru", *lru)
  if err != nil {
    log.Println("error saving lru to redis", err)
    return
  }

  log.Println("successfully saved lru to redis")
}

func resizeLRU(lru *LRU, ch chan<- *LRU) {
  if len(*lru) < 1000 {
    // not big enough
    return
  }

  start := len(*lru) - 1000
  alru := *lru
  nlru := alru[start:]
  ch <- &nlru
  log.Println("created new lru")
}

func Worker() {
  LRUChan = make(chan *LRUItem, 100)
  LRURefChan = make(chan chan<- *LRU, 100)
  newRefChan := make(chan *LRU, 1)

  var a LRU
  var lru = &a

  populateLRU(lru)

  go func() {
    for {
      select {
      case i := <-LRUChan:
        // add entity to LRU
        *lru = append(*lru, i)
        log.Println("added entity to lru")

        if len(*lru) > 1500 {
          go resizeLRU(lru, newRefChan)
        }
      case c := <-LRURefChan:
        // return LRU ref
        c <- lru
        log.Println("returned lru reference")
      case nlru := <-newRefChan:
        // change the lru reference
        lru = nlru
        log.Println("changed lru variable")
        saveLRU(lru)
      }
    }
  }()
}

//var VoteChan chan VoteItem

/*
func Worker() {
  VoteChan = make(chan Item, 100)

  go func() {
    for {
      entity := <-VoteChan
      item := entity.Item

      b, err := json.Marshal(item)
      if err != nil {
        log.Println(err)
        continue
      }

      var i interface{}
      key := fmt.Sprintf("%s:%s", item.Type, item.Title)
      i, err = Redis.Do("GET", key)
      if err != nil || i == nil {
        // TODO
      }
    }
  }
}

*/
