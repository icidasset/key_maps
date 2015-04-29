package api

import (
	"encoding/base64"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/gedex/inflector"
	"github.com/icidasset/key-maps-api/db"
	"github.com/labstack/echo"
)

type PublicItem map[string]interface{}
type PublicCollection []PublicItem

//
//  {public/handler} Show
//
func Public__Show(c *echo.Context) {
	theMap, err := publicGetMap(c.Param("hash"))

	if err != nil {
		c.JSON(500, FormatError(err))
		return
	} else if theMap.Id == 0 {
		c.JSON(501, map[string]string{"error": "Provided map id and slug do not match"})
		return
	}

	theMapSettings, err := publicGetMapSettings(&theMap)

	if err != nil {
		c.JSON(500, FormatError(err))
		return
	}

	theMapItems, err := publicGetMapItems(&theMap, &theMapSettings)

	if err != nil {
		c.JSON(500, FormatError(err))
		return
	}

	collection, err := publicMakeCollection(&theMap, &theMapItems)

	if err != nil {
		c.JSON(500, FormatError(err))
	} else if theMapSettings.IncludeJSONRoot {
		c.JSON(200, map[string]PublicCollection{theMap.Slug: collection})
	} else {
		c.JSON(200, collection)
	}
}

//
//  {private} Get map
//
func publicGetMap(hash string) (Map, error) {
	decodedHash, err := base64.StdEncoding.DecodeString(hash)
	processedHash := strings.Split(string(decodedHash[:]), "/")

	// params
	theMapId, err := strconv.ParseInt(processedHash[0], 10, 0)
	theMapSlug := processedHash[1]

	// map
	theMap := Map{}

	// make database request if no error
	if err == nil {
		err = db.Inst().Get(
			&theMap,
			"SELECT * FROM maps WHERE id = $1 AND slug = $2",
			theMapId,
			theMapSlug,
		)
	}

	// return
	return theMap, err
}

//
//  {private} Get map settings
//
func publicGetMapSettings(theMap *Map) (MapSettings, error) {
	theSettings := MapSettings{}
	err := json.Unmarshal([]byte(theMap.Settings), &theSettings)
	return theSettings, err
}

//
//  {private} Get map items
//
func publicGetMapItems(theMap *Map, theMapSettings *MapSettings) ([]MapItem, error) {
	mapItems := []MapItem{}
	mapItemsQuery := "SELECT structure_data FROM map_items WHERE map_id = $1"

	if theMapSettings.SortBy != "" {
		mapItemsQuery = mapItemsQuery +
			" ORDER BY structure_data::json->>'" + theMapSettings.SortBy + "'"
	}

	err := db.Inst().Select(
		&mapItems,
		mapItemsQuery,
		theMap.Id,
	)

	return mapItems, err
}

//
//  {private} Make collection
//
func publicMakeCollection(theMap *Map, theMapItems *[]MapItem) (PublicCollection, error) {
	var err error

	collection := make(PublicCollection, 0)

	for _, mapItem := range *theMapItems {
		item := make(PublicItem)
		err = json.Unmarshal([]byte(mapItem.StructureData), &item)

		if err != nil {
			break
		}

		err = publicProcessItem(theMap, &item)

		if err != nil {
			break
		} else {
			collection = append(collection, item)
		}
	}

	return collection, err
}

//
//  {private} Process item
//
func publicProcessItem(theMap *Map, item *PublicItem) error {
	var err error

	for itemKey, itemValue := range *item {

		// associations
		if strings.Contains(itemKey, "->") {
			err = publicProcessItemAssociation(theMap, itemKey, itemValue, item)
		}

		// check for errors
		if err != nil {
			break
		}

	}

	return err
}

//
//  {private} Process item association
//
func publicProcessItemAssociation(theMap *Map, itemKey string, itemValue interface{}, item *PublicItem) error {
	var associationType string
	var err error

	// get map structure
	theMapStructure := make([]map[string]string, 0)
	err = json.Unmarshal([]byte(theMap.Structure), &theMapStructure)

	if err != nil {
		return err
	}

	// find association type
	for _, s := range theMapStructure {
		if s["key"] == itemKey {
			associationType = strings.Split(s["type"], ".")[1]
			break
		}
	}

	// last step
	if associationType == "one" {
		err = publicProcessItemAssociationOne(itemKey, itemValue, item)
	} else if associationType == "many" {
		// TODO: many
	}

	return err
}

func publicProcessItemAssociationOne(key string, associationValue interface{}, item *PublicItem) error {
	mapItemId, err := strconv.ParseInt(associationValue.(string), 10, 0)

	if err != nil {
		return err
	}

	mapItem := MapItem{}
	mapItemPublic := make(PublicItem)

	// database request
	err = db.Inst().Get(
		&mapItem,
		`SELECT * FROM map_items
		WHERE id = $1`,
		mapItemId,
	)

	if err != nil {
		return err
	}

	// parse structure data
	err = json.Unmarshal(
		[]byte(mapItem.StructureData),
		&mapItemPublic,
	)

	if err != nil {
		return err
	}

	// add item to structure
	associationMapSlug := strings.Split(key, "->")[0]
	associationMapSlugSingularized := inflector.Singularize(associationMapSlug)

	(*item)[associationMapSlugSingularized] = mapItemPublic

	// remove old item from structure
	delete(*item, key)

	// return
	return nil
}
