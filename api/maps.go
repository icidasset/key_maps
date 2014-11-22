package api

import (
  // "fmt"
  "github.com/extemporalgenome/slug"
  "github.com/go-martini/martini"
  "github.com/icidasset/key-maps/db"
  "github.com/martini-contrib/render"
  // "net/http"
  "time"
)

//
//  [API]
//  -> https://gobyexample.com/json
//  -> http://www.alexedwards.net/blog/golang-response-snippets
//  -> http://godoc.org/github.com/lib/pq
//

type Map struct {
  Id int              `json:"id" db:"id"`
  Name string         `json:"name" db:"name"`
  Slug string         `json:"slug"`
  Structure string    `json:"structure"`
  CreatedAt time.Time `json:"created_at" db:"created_at"`
  UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}


type MapFormData struct {
  Name string `form:"maps[name]" binding:"required"`
  Structure string `form:"maps[structure]" binding:"required"`
}


//
//  [MAPS]
//
func Maps__Index() {
  maps := []Map{}
  db.Inst().Select(&maps, "SELECT * FROM maps")
}


func Maps__Show(params martini.Params, r render.Render) {
  m := Map{}

  // execute query
  db.Inst().Get(&m, "SELECT * FROM maps WHERE id=$1", params["id"])

  // if none found
  if m.Id == 0 {
    r.JSON(404, nil)

  // render map as json
  } else {
    r.JSON(200, m)
  }
}


func Maps__Create(mfd MapFormData, r render.Render) {
  query := "INSERT INTO maps (name, slug, structure)" +
           " VALUES (:name, :slug, :structure)"

  // make new map
  slug := slug.Slug(mfd.Name)
  now := time.Now()

  new_map := Map{Name: mfd.Name, Slug: slug, Structure: mfd.Structure, CreatedAt: now, UpdatedAt: now}

  // execute query
  _, err := db.Inst().NamedExec(query, new_map)

  // if error
  if err != nil {
    r.JSON(500, nil)

  // render map as json
  } else {
    m := Map{}

    db.Inst().Get(&m, "SELECT * FROM maps WHERE slug = ?", slug)

    r.JSON(200, m)

  }
}
