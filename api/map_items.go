package api

import (
  "fmt"
  "github.com/go-martini/martini"
  "github.com/icidasset/key-maps/db"
  "github.com/martini-contrib/render"
  "net/http"
  "time"
)


type MapItem struct {
  Id int                  `json:"id"`
  StructureData string    `json:"structure_data" db:"structure_data" form:"structure_data"`
  CreatedAt time.Time     `json:"created_at" db:"created_at"`
  UpdatedAt time.Time     `json:"updated_at" db:"updated_at"`
  MapId int               `json:"map_id" db:"map_id"`
}


type MapItemFormData struct {
  MapItem MapItem         `form:"map_item" binding:"required"`
}



//
//  {get} SHOW
//
func MapItems__Show(params martini.Params, r render.Render, u User) {
  mi := MapItem{}

  // execute query
  err := db.Inst().Get(
    &mi,
    "SELECT * FROM map_items WHERE id = $1",
    params["id"],
  )

  // render
  if err != nil {
    panic(err)
  } else if mi.Id == 0 {
    r.JSON(404, nil)
  } else {
    r.JSON(200, map[string]MapItem{ "map_item": mi })
  }
}



//
//  {post} CREATE
//
func MapItems__Create(mifd MapItemFormData, w http.ResponseWriter, r render.Render, u User) {
  query := "INSERT INTO map_items (structure_data, created_at, updated_at, map_id) VALUES (:structure_data, :created_at, :updated_at, :map_id)"

  // make new map item
  now := time.Now()

  new_map_item := MapItem{StructureData: mifd.MapItem.StructureData, CreatedAt: now, UpdatedAt: now, MapId: mifd.MapItem.MapId}

  // execute query
  result, _ := db.Inst().NamedExec(query, new_map_item)

  fmt.Fprintf(w, "%#v", result)
}
