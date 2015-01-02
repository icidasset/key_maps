package api

import (
  "net/http"
  "github.com/extemporalgenome/slug"
  "github.com/go-martini/martini"
  "github.com/icidasset/key-maps/db"
  "github.com/martini-contrib/render"
  "strconv"
  "strings"
  "time"
)


type Map struct {
  Id int                  `json:"id"`
  Slug string             `json:"slug"`
  Name string             `json:"name"`
  Structure string        `json:"structure"`
  CreatedAt time.Time     `json:"created_at" db:"created_at"`
  UpdatedAt time.Time     `json:"updated_at" db:"updated_at"`
  MapItems IntSlice       `json:"map_items" db:"map_items"`
  UserId int              `json:"-" db:"user_id"`
}


type MapFormData struct {
  Map Map                 `json:"map" binding:"required"`
}



//
//  IntSlice
//
type IntSlice []int


func (i *IntSlice) Scan(src interface{}) error {
  as_bytes, _ := src.([]byte)
  as_string := string(as_bytes)

  a := strings.Split(as_string, ", ")
  b := make([]int, 0)

  for _, x := range a {
    i, _ := strconv.Atoi(x)
    b = append(b, i)
  }

  int_slice := IntSlice(b)

  if int_slice[0] == 0 {
    (*i) = nil
  } else {
    (*i) = int_slice
  }

  return nil
}



//
//  {get} INDEX
//
func Maps__Index(w http.ResponseWriter, r render.Render, u User) {
  maps := []Map{}

  // execute query
  rows, err := db.Inst().Queryx(
    `SELECT maps.* AS map_id,
            array_to_string(array(
              SELECT id FROM map_items
              WHERE maps.id = map_items.map_id
            ), ', ') AS map_items
     FROM maps
     WHERE maps.user_id = $1
     ORDER BY maps.id;`,
     u.Id,
  )

  for rows.Next() {
    m := Map{}
    err = rows.StructScan(&m)
    maps = append(maps, m)
  }

  // render
  if err != nil {
    r.JSON(500, err.Error())
  } else {
    r.JSON(200, map[string][]Map{ "maps": maps })
  }
}



//
//  {get} SHOW
//
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
    r.JSON(500, err.Error())
  } else if m.Id == 0 {
    r.JSON(404, nil)
  } else {
    r.JSON(200, map[string]Map{ "map": m })
  }
}



//
//  {post} CREATE
//
func Maps__Create(mfd MapFormData, r render.Render, u User) {
  query := "INSERT INTO maps (name, slug, structure, created_at, updated_at, user_id) VALUES (:name, :slug, :structure, :created_at, :updated_at, :user_id) RETURNING id"

  // make new map
  slug := slug.Slug(mfd.Map.Name)
  now := time.Now()

  new_map := Map{Name: mfd.Map.Name, Slug: slug, Structure: mfd.Map.Structure, CreatedAt: now, UpdatedAt: now, UserId: u.Id}

  // execute query
  rows, err := db.Inst().NamedQuery(query, new_map)

  for rows.Next() {
    rows.StructScan(&new_map)
  }

  // render
  if err != nil {
    r.JSON(500, err.Error())
  } else {
    r.JSON(200, map[string]Map{ "map": new_map })
  }
}



//
//  {put} UPDATE
//
func Maps__Update(mfd MapFormData, params martini.Params, r render.Render, u User) {
  _, err := db.Inst().Exec(
    "UPDATE maps SET structure = $1, updated_at = $2 WHERE id = $3 AND user_id = $4",
    mfd.Map.Structure,
    time.Now(),
    params["id"],
    u.Id,
  )

  // fetch
  m := Map{}

  db.Inst().Get(
    &m,
    "SELECT * FROM maps WHERE id = $1 AND user_id = $2",
    params["id"],
    u.Id,
  )

  // render
  if err != nil {
    r.JSON(500, err.Error())
  } else {
    r.JSON(200, map[string]Map{ "map": m })
  }
}



//
//  {delete} DESTROY
//
func Maps__Destroy(params martini.Params, r render.Render, u User) {
  _, err := db.Inst().Exec(
    "DELETE FROM maps WHERE id = $1 AND user_id = $2",
    params["id"],
    u.Id,
  )

  // render
  if err != nil {
    r.JSON(500, err.Error())
  } else {
    r.JSON(200, nil)
  }
}
