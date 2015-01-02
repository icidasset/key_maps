package api

import (
  "github.com/go-martini/martini"
  "github.com/icidasset/key-maps/db"
  "github.com/martini-contrib/render"
  "net/http"
  "time"
)


type MapItem struct {
  Id int                  `json:"id"`
  StructureData string    `json:"structure_data" db:"structure_data"`
  CreatedAt time.Time     `json:"created_at" db:"created_at"`
  UpdatedAt time.Time     `json:"updated_at" db:"updated_at"`
  MapId int               `json:"map_id,string" db:"map_id"`
}


type MapItemFormData struct {
  MapItem MapItem         `json:"map_item" binding:"required"`
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
    r.JSON(500, err.Error())
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
  query := "INSERT INTO map_items (structure_data, created_at, updated_at, map_id) VALUES (:structure_data, :created_at, :updated_at, :map_id) RETURNING *"

  // make new map item
  now := time.Now()

  new_map_item := MapItem{StructureData: mifd.MapItem.StructureData, CreatedAt: now, UpdatedAt: now, MapId: mifd.MapItem.MapId}

  // execute query
  rows, err := db.Inst().NamedQuery(query, new_map_item)

  for rows.Next() {
    rows.StructScan(&new_map_item)
  }

  // render
  if err != nil {
    r.JSON(500, err.Error())
  } else {
    r.JSON(200, map[string]MapItem{ "map_item": new_map_item })
  }
}



//
//  {put} UPDATE
//
func MapItems__Update(mifd MapItemFormData, params martini.Params, r render.Render, u User) {
  _, err := db.Inst().Exec(
    "UPDATE map_items SET structure_data = $1, updated_at = $2 WHERE id = $3",
    mifd.MapItem.StructureData,
    time.Now(),
    params["id"],
  )

  // fetch
  mi := MapItem{}

  db.Inst().Get(
    &mi,
    "SELECT * FROM map_items WHERE id = $1",
    params["id"],
  )

  // render
  if err != nil {
    r.JSON(500, err.Error())
  } else {
    r.JSON(200, map[string]MapItem{ "map_item": mi })
  }
}



//
//  {delete} DESTROY
//
func MapItems__Destroy(params martini.Params, r render.Render, u User) {
  _, err := db.Inst().Exec(
    "DELETE FROM map_items WHERE id = $1",
    params["id"],
  )

  // render
  if err != nil {
    r.JSON(500, err.Error())
  } else {
    r.JSON(200, nil)
  }
}
