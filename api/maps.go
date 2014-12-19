package api

import (
  "github.com/extemporalgenome/slug"
  "github.com/go-martini/martini"
  "github.com/icidasset/key-maps/db"
  _ "github.com/lib/pq"
  "github.com/martini-contrib/render"
  "time"
)


type Map struct {
  Id int                  `json:"id"`
  Slug string             `json:"slug"`
  Name string             `json:"name" form:"name" binding:"required"`
  Structure string        `json:"structure" form:"structure" binding:"required"`
  CreatedAt time.Time     `json:"created_at" db:"created_at"`
  UpdatedAt time.Time     `json:"updated_at" db:"updated_at"`
  UserId int              `json:"user_id" db:"user_id"`
}


type MapFormData struct {
  Map Map                 `form:"map" binding:"required"`
}



//
//  Routes
//
func Maps__Index(r render.Render, u User) {
  m := []Map{}

  // execute query
  err := db.Inst().Select(&m, "SELECT * FROM maps WHERE user_id = $1", u.Id)

  // render
  if err != nil {
    panic(err)
  } else {
    r.JSON(200, map[string][]Map{ "maps": m })
  }
}


func Maps__Show(params martini.Params, r render.Render, u User) {
  m := Map{}

  // execute query
  err := db.Inst().Get(
    &m,
    "SELECT * FROM maps WHERE id = $1 AND user_id = $2",
    params["id"],
    u.Id,
  )

  // render
  if err != nil {
    panic(err)
  } else if m.Id == 0 {
    r.JSON(404, nil)
  } else {
    r.JSON(200, map[string]Map{ "map": m })
  }
}


func Maps__Create(mfd MapFormData, r render.Render, u User) {
  query := "INSERT INTO maps (name, slug, structure, created_at, updated_at, user_id) VALUES (:name, :slug, :structure, :created_at, :updated_at, :user_id)"

  // make new map
  slug := slug.Slug(mfd.Map.Name)
  now := time.Now()

  new_map := Map{Name: mfd.Map.Name, Slug: slug, Structure: mfd.Map.Structure, CreatedAt: now, UpdatedAt: now, UserId: u.Id}

  // execute query
  _, err := db.Inst().NamedExec(query, new_map)

  // if error
  if err != nil {
    r.JSON(500, err.Error())

  // render map as json
  } else {
    m := Map{}
    db.Inst().Get(&m, "SELECT * FROM maps WHERE slug = $1", slug)
    r.JSON(200, map[string]Map{ "map": m })

  }
}
