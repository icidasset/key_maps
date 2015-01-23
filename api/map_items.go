package api

import (
  "encoding/json"
  "github.com/gocraft/web"
  "github.com/icidasset/key-maps/db"
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
  MapItem MapItem         `json:"map_item"`
}



//
//  {get} SHOW
//
func (c *Context) MapItems__Show(rw web.ResponseWriter, req *web.Request) {
  mi := MapItem{}

  // execute query
  err := db.Inst().Get(
    &mi,
    "SELECT * FROM map_items WHERE id = $1",
    req.PathParams["id"],
  )

  // render
  if err != nil {
    RenderJSON(rw, 500, FormatError(err))
  } else if mi.Id == 0 {
    RenderJSON(rw, 404, nil)
  } else {
    RenderJSON(rw, 200, map[string]MapItem{ "map_item": mi })
  }
}



//
//  {post} CREATE
//
func (c *Context) MapItems__Create(rw web.ResponseWriter, req *web.Request) {
  query := "INSERT INTO map_items (structure_data, created_at, updated_at, map_id) VALUES (:structure_data, :created_at, :updated_at, :map_id) RETURNING *"

  // parse json from request body
  mifd := MapItemFormData{}
  json_decoder := json.NewDecoder(req.Body)
  json_decoder.Decode(&mifd)

  // make new map item
  now := time.Now()

  new_map_item := MapItem{StructureData: mifd.MapItem.StructureData, CreatedAt: now, UpdatedAt: now, MapId: mifd.MapItem.MapId}

  // execute query
  rows, err := db.Inst().NamedQuery(query, new_map_item)

  // return if error
  if err != nil {
    RenderJSON(rw, 500, FormatError(err))
    return
  }

  // scan rows
  for rows.Next() {
    err = rows.StructScan(&new_map_item)
  }

  // render
  if err != nil {
    RenderJSON(rw, 500, FormatError(err))
  } else {
    RenderJSON(rw, 201, map[string]MapItem{ "map_item": new_map_item })
  }
}



//
//  {put} UPDATE
//
func (c *Context) MapItems__Update(rw web.ResponseWriter, req *web.Request) {
  mifd := MapItemFormData{}
  json_decoder := json.NewDecoder(req.Body)
  json_decoder.Decode(&mifd)

  // update map item
  _, err := db.Inst().Exec(
    "UPDATE map_items SET structure_data = $1, updated_at = $2 WHERE id = $3",
    mifd.MapItem.StructureData,
    time.Now(),
    req.PathParams["id"],
  )

  // return if error
  if err != nil {
    RenderJSON(rw, 500, FormatError(err))
    return
  }

  // fetch
  mi := MapItem{}

  err = db.Inst().Get(
    &mi,
    "SELECT * FROM map_items WHERE id = $1",
    req.PathParams["id"],
  )

  // render
  if err != nil {
    RenderJSON(rw, 500, FormatError(err))
  } else {
    RenderJSON(rw, 200, map[string]MapItem{ "map_item": mi })
  }
}



//
//  {delete} DESTROY
//
func (c *Context) MapItems__Destroy(rw web.ResponseWriter, req *web.Request) {
  _, err := db.Inst().Exec(
    "DELETE FROM map_items WHERE id = $1",
    req.PathParams["id"],
  )

  // render
  if err != nil {
    RenderJSON(rw, 500, FormatError(err))
  } else {
    RenderJSON(rw, 204, nil)
  }
}
