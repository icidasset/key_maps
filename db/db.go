package db

import (
  _ "database/sql"
  "github.com/jmoiron/sqlx"
  _ "github.com/lib/pq"
)


var db *sqlx.DB


func Open() error {
  var err error
  db, err = sqlx.Open("postgres", "user=icidasset dbname=keymaps_development")
  return err
}


func Close() {
  db.Close()
}


func Select(hash interface{}, query string, args ...interface{}) error {
  return db.Select(hash, query)
}


func Get(hash interface{}, query string, args ...interface{}) error {
  return db.Get(hash, query)
}
