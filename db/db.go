package db

import (
  "database/sql"
  "github.com/jmoiron/sqlx"
  "github.com/lib/pq"
  "github.com/rubenv/sql-migrate"
  "gopkg.in/yaml.v2"
  "io/ioutil"
  "os"
)

var db *sql.DB
var dbSqlx *sqlx.DB
var env string


func Inst() *sqlx.DB {
  return dbSqlx;
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
  env_ds := env_config["datasource"].(string)
  env_datasource := env_ds

  switch env_ds {
    case "HEROKU": env_datasource = getHerokuDataSource()
  }

  // open db
  db, err = sql.Open("postgres", env_datasource)
  dbSqlx = sqlx.NewDb(db, "postgres")

  // run migrations if on heroku
  if env_ds == "HEROKU" {
    runMigrations()
  }

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



//
//  Migrations
//
func runMigrations() {
  migrations := &migrate.FileMigrationSource{
    Dir: "db/migrations",
  }

  _, err := migrate.Exec(db, "postgres", migrations, migrate.Up)

  if err != nil {
    panic(err)
  }
}
