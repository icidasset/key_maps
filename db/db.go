package db

import (
  _ "database/sql"
  "github.com/jmoiron/sqlx"
  "github.com/lib/pq"
  "gopkg.in/yaml.v2"
  "io/ioutil"
  "os"
)


var db *sqlx.DB
var env string


func Inst() *sqlx.DB {
  return db;
}


func Open() error {
  config_file_data, err := ioutil.ReadFile("db/config.yml")

  // panic if readFile error
  if err != nil {
    panic(err)
  }

  // parse yaml
  config := make(map[string]interface{})
  err = yaml.Unmarshal([]byte(config_file_data), &config)

  // panic if yaml error
  if err != nil {
    panic(err)
  }

  // env
  env = os.Getenv("ENV")
  if env == "" { env = "development" }

  // set datasource
  env_config := config[env].(map[interface{}]interface{})
  env_datasource := env_config["datasource"].(string)

  switch env_datasource {
    case "HEROKU": env_datasource = getHerokuDataSource()
  }

  // open db
  db, err = sqlx.Open("postgres", env_datasource)

  // return db error
  return err
}


func Close() {
  db.Close()
}



//
//  Environment specific functions
//
func getHerokuDataSource() string {
  url := os.Getenv("DATABASE_URL")
  connection, _ := pq.ParseURL(url)
  connection += " sslmode=require"

  return connection
}
