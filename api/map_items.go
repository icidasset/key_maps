package api

import (
	"encoding/json"
	"time"

	"github.com/icidasset/key-maps-api/db"
	"github.com/labstack/echo"
)

type MapItem struct {
	Id            int       `json:"id"`
	StructureData string    `json:"structure_data" db:"structure_data"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
	MapId         int       `json:"map_id,string" db:"map_id"`
}

type MapItemFormData struct {
	MapItem MapItem `json:"map_item"`
}

//
//  {get} SHOW
//
func MapItems__Show(c *echo.Context) error {
	mi := MapItem{}

	// execute query
	err := db.Inst().Get(
		&mi,
		`SELECT * FROM map_items
     WHERE id = $1 AND map_id IN (SELECT id FROM maps WHERE user_id = $2)`,
		c.Param("id"),
		c.Get("user").(User).Id,
	)

	// render
	if err != nil {
		if IsNoResultsError(err.Error()) {
			return c.JSON(404, nil)
		} else {
			return c.JSON(500, FormatError(err))
		}
	} else if mi.Id == 0 {
		return c.JSON(404, nil)
	} else {
		return c.JSON(200, map[string]MapItem{"map_item": mi})
	}
}

//
//  {post} CREATE
//
func MapItems__Create(c *echo.Context) error {
	query := `INSERT INTO map_items (structure_data, created_at, updated_at, map_id) VALUES (:structure_data, :created_at, :updated_at, :map_id) RETURNING *`

	// parse json from request body
	mifd := MapItemFormData{}
	json_decoder := json.NewDecoder(c.Request().Body)
	json_decoder.Decode(&mifd)

	// make new map item
	now := time.Now()

	new_map_item := MapItem{StructureData: mifd.MapItem.StructureData, CreatedAt: now, UpdatedAt: now, MapId: mifd.MapItem.MapId}

	// execute query
	rows, err := db.Inst().NamedQuery(query, new_map_item)

	// return if error
	if err != nil {
		return c.JSON(500, FormatError(err))
	}

	// scan rows
	for rows.Next() {
		err = rows.StructScan(&new_map_item)
	}

	// render
	if err != nil {
		return c.JSON(500, FormatError(err))
	} else {
		return c.JSON(201, map[string]MapItem{"map_item": new_map_item})
	}
}

//
//  {put} UPDATE
//
func MapItems__Update(c *echo.Context) error {
	mifd := MapItemFormData{}
	json_decoder := json.NewDecoder(c.Request().Body)
	json_decoder.Decode(&mifd)

	// update map item
	_, err := db.Inst().Exec(
		`UPDATE map_items
     SET structure_data = $1, updated_at = $2
     WHERE id = $3 AND map_id IN (SELECT id FROM maps WHERE user_id = $4)`,
		mifd.MapItem.StructureData,
		time.Now(),
		c.Param("id"),
		c.Get("user").(User).Id,
	)

	// return if error
	if err != nil {
		return c.JSON(500, FormatError(err))
	}

	// fetch
	mi := MapItem{}

	err = db.Inst().Get(
		&mi,
		`SELECT * FROM map_items
     WHERE id = $1 AND map_id IN (SELECT id FROM maps WHERE user_id = $2)`,
		c.Param("id"),
		c.Get("user").(User).Id,
	)

	// render
	if err != nil {
		return c.JSON(500, FormatError(err))
	} else {
		return c.JSON(200, map[string]MapItem{"map_item": mi})
	}
}

//
//  {delete} DESTROY
//
func MapItems__Destroy(c *echo.Context) error {
	_, err := db.Inst().Exec(
		`DELETE FROM map_items
     WHERE id = $1 AND map_id IN (SELECT id FROM maps WHERE user_id = $2)`,
		c.Param("id"),
		c.Get("user").(User).Id,
	)

	// render
	if err != nil {
		return c.JSON(500, FormatError(err))
	} else {
		return c.JSON(204, nil)
	}
}
