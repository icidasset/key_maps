package api

import (
  "github.com/extemporalgenome/slug"
  "github.com/go-martini/martini"
  "github.com/lib/pq"
  "github.com/icidasset/key-maps/db"
  "github.com/martini-contrib/render"
  _ "net/http"
  "time"
)

//
//  [API]
//  -> https://gobyexample.com/json
//  -> http://www.alexedwards.net/blog/golang-response-snippets
//  -> http://godoc.org/github.com/lib/pq
//

type Map struct {
  Id int                  `json:"id"`
  Name string             `json:"name"`
  Slug string             `json:"slug"`
  Structure string        `json:"structure"`
  CreatedAt pq.NullTime   `json:"created_at" db:"created_at"`
  UpdatedAt pq.NullTime   `json:"updated_at" db:"updated_at"`
}


type MapFormData struct {
  Name string         `form:"maps[name]" binding:"required"`
  Structure string    `form:"maps[structure]" binding:"required"`
}


//
//  [MAPS]
//
func Maps__Index(r render.Render) {
  m := []Map{}

  // execute query
  err := db.Inst().Select(&m, "SELECT * FROM maps")

  // render
  if err != nil {
    panic(err)
  } else {
    r.JSON(200, m)
  }
}


func Maps__Show(params martini.Params, r render.Render) {
  m := Map{}

  // execute query
  db.Inst().Get(&m, "SELECT * FROM maps WHERE id = $1", params["id"])

  // if none found
  if m.Id == 0 {
    r.JSON(404, nil)

  // render map as json
  } else {
    r.JSON(200, m)

  }
}


func Maps__Create(mfd MapFormData, r render.Render) {
  query := "INSERT INTO maps (name, slug, structure, created_at, updated_at)" +
           " VALUES (:name, :slug, :structure, :created_at, :updated_at)"

  // make new map
  slug := slug.Slug(mfd.Name)
  now := pq.NullTime{}
  now.Scan(time.Now())

  new_map := Map{Name: mfd.Name, Slug: slug, Structure: mfd.Structure, CreatedAt: now, UpdatedAt: now}

  // execute query
  _, err := db.Inst().NamedExec(query, new_map)

  // if error
  if err != nil {
    r.JSON(500, err.Error())

  // render map as json
  } else {
    m := Map{}

    db.Inst().Get(&m, "SELECT * FROM maps WHERE slug = $1", slug)

    r.JSON(200, m)

  }
}
