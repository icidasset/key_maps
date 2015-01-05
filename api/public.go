package api

import (
  "encoding/base64"
  "github.com/go-martini/martini"
  "github.com/icidasset/key-maps/db"
  "github.com/martini-contrib/render"
  "net/http"
  "strconv"
)


func Public__Show(params martini.Params, r render.Render, w http.ResponseWriter) {
  data, _ := base64.StdEncoding.DecodeString(params["hash"])
  str := string(data[:])

  // params
  map_id, _ := strconv.ParseInt(str, 10, 0)
  map_id = map_id / 25

  // query
  map_items := []MapItem{}

  db.Inst().Select(
    &map_items,
    "SELECT structure_data FROM map_items WHERE map_id = $1",
    map_id,
  )

  r.JSON(200, map_items)
}
