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


func Query(query string) {
  // rows, _ := db.Query("SELECT * FROM maps")

  // defer rows.Close()
  // for rows.Next() {}
  // return rows
}
