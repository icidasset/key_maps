package db

import (
  _ "database/sql"
  "github.com/jmoiron/sqlx"
  _ "github.com/lib/pq"
)


var db *sqlx.DB


func Inst() *sqlx.DB {
  return db;
}


func Open() error {
  var err error
  db, err = sqlx.Open("postgres", "user=icidasset dbname=keymaps_development sslmode=disable")
  return err
}


func Close() {
  db.Close()
}
