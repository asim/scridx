package main

import (
  "flag"
  "log"
  "net/http"
  "os"
  "path/filepath"
  _ "github.com/ziutek/mymysql/godrv"
  "scridx/app"
)

var addr = flag.String("addr", ":8080", "http service address")

var (
  BasePath, _  = os.Getwd()
  AppPath      = filepath.Join(BasePath, "app")
  ViewsPath    = filepath.Join(BasePath, "views")
  DBConfigPath = filepath.Join(BasePath, "db/config.json")
  UserDBConfigPath = filepath.Join(BasePath, "users/db/config.json")
  Addr         = addr
  Environment  = "production"
)

// Main
func main() {
  hd, err := app.InitDB(DBConfigPath, Environment)
  defer hd.Db.Close()
  if err != nil {
    log.Fatal(err)
  }

  // Serve
  err = http.ListenAndServe(*Addr, app.RegisterRoutes())
  if err != nil {
    log.Fatal(err)
  }
}
