package db

import (
  "database/sql"
  _ "github.com/lib/pq"
)


var db *sql.DB


func Open() {
  db, _ = sql.Open("postgres", "user=icidasset dbname=keymaps_development")
}


func Close() {
  db.Close()
}


func Query(query string) *sql.Rows {
  rows, _ := db.Query(query)
  return rows
}
