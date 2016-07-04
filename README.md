# Key Maps

A simple data-structure API.

__Work in progress.__

_You can find old code in the legacy branches._



## How it works

1. You make a map (like a database table), e.g. "Quotes"
2. You add data to the map, e.g. a quote and its author
3. You can fetch this data through a public JSON API



## The API

### Authentication

```markdown
POST  /sign-up

__Request body:__

{
  "email": "...",
  "password": "...",
  "username": "..."
}

__Response body:__

{
  "data": {
    "token": "..."
  }
}

POST  /sign-in

__Request body:__

{
  "login": "...",
  "password": "..."
}

__Response body:__

{
  "data": {
    "token": "..."
  }
}
```

Use the `token` to authenticate requests.  
For example:

```http
GET /api?query=query
Authorization: PLACE_TOKEN_HERE
```


### Private API

__Endpoint__

```http
GET /api?query=PLACE_QUERY_HERE
```

__GraphQL queries__

```bash
mutation M { createMap(
  name: "Quotes",
  attributes: [ "quote", "author" ]
) {}}

mutation M { createMapItem(
  map: "Quotes",

  quote: "Specialization tends to shut off the wide-band tuning searches and thus to preclude further discovery.",
  author: "Buckminster Fuller"
) { attributes }}

query Q { mapItems(map: "Quotes") {
  id,
  attributes
}}

query Q { maps() {
  id,
  name,
  attributes
}}

query Q { mapItem(id: ITEM_ID) { attributes }}
query Q { map(name: "Quotes") { attributes }}
# -- uses the name argument to select the map,
#    but you can also use the map id.

mutation M { updateMapItem(id: ITEM_ID, quote: "Updated quote") { quote }}
mutation M { updateMap(id: MAP_ID, name: "Updated name") { name }}

mutation M { removeMapItem(id: ITEM_ID) { id }}
mutation M { removeMap(id: MAP_ID) { id }}
```

__Notes__  
The map name must be unique, it will be casted to lowercase for validation.

#### Defining types for your attributes (optional)

You can define types, but it's totally optional, and it doesn't actually do anything either.
So why define these, well, you could do this to show a certain input field in your UI.
For example, if you want to show a date selector.

```
mutation M { createMap(
  name: "Author",
  attributes: [ "date_of_birth" ]
  types: { date_of_birth: "date" }
) {}}
```


### Public API

__All__ items for a single map:

```
GET /public/:username/:map_name

:username, your case-insensitive username
:map_name, your case-insensitive map name
```

__One__ item for a single map:

```
GET /public/:username/:map_name/:map_item_id
```

__Options__

```markdown
GET /public/:username/:map_name?sort_by=author

**sort_by**, e.g. 'author', when not specified, it is sorted by insertion date.  
**sort_direction**, 'asc' or 'desc', default is 'asc'.  
**timestamps**, include this to add the timestamps of the item.  
```


### Responses

The API will always return data in one of the following formats.

__Data__

```json
{
  "data": {
    "attribute": "example"
  }
}
```

In the case of a GraphQL query or mutation.

```json
{
  "data": {
    "query_or_mutation_name": {
      "attribute": "example"
    }
  }
}
```

__Errors__

```json
{
  "errors": [
    { "message": "Error message" }
  ]
}
```



## Development

```
source .env

mix deps.get
mix test

mix ecto.create
mix run --no-halt
```
