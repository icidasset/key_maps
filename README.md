# Key Maps

A simple data-structure API.

_You can find old code in the legacy branches._



## How it works

1. You make a map (resembles a database table), e.g. "Quotes"
2. You add data to the map, e.g. a quote and its author
3. You can fetch this data through a public JSON API



## The API

### Authentication

Uses Auth0's passwordless authentication.

```markdown
POST /auth/start

__Request body:__

{
  "email": "..."
}

__Response body:__

Status 200 with empty body if successful
Status 422 with { errors: ... } if unsuccessful


POST /auth/exchange

__Request body:__

{
  auth0_id_token: "..."
}

__Response body:__

{
  "data": {
    "token": "..."
  }
}


POST /auth/validate

__Request body:__

{
  "token": "..."
}

__Response body:__

Status 200 with empty body if valid
Status 403 with empty body if invalid
```

Use the `token` to authenticate requests.
User info is located inside the token's claims.  
For example:

```http
GET /api?query=query
Authorization: PLACE_TOKEN_HERE
```

__Notes__

The ability to sign-up is disabled by default,
you have to define the ENV variable 'ENABLE_SIGN_UP=1' to enable it.


### Private API

__Endpoint__

```http
GET /api?query=PLACE_QUERY_HERE
```

__GraphQL queries__

```bash
# 1. Create a `map`
mutation M { createMap(
  name: "Quotes",
  attributes: [ "quote", "author" ]
) {
  id,
  attributes,
  types
}}

# 2. Create an item for a `map`
mutation M { createMapItem(
  map: "Quotes",

  quote: "Specialization tends to shut off the wide-band tuning searches and thus to preclude further discovery.",
  author: "Buckminster Fuller"
) {
  id,
  map_id,
  attributes
}}

# 3. Get all map items for a specific map
query Q { mapItems(map: "Quotes") {
  id,
  map_id,
  attributes
}}

# 4. Get all maps
query Q { maps() {
  id,
  name,
  attributes,
  types
}}

# 5. Get a specific item and a specific map
query Q { mapItem(id: ITEM_ID) { ... }}
query Q { map(name: "Quotes") { ... }}
# -- uses the name argument to select the map,
#    but you can also use the map id.

# 6. Update a map item and a map
mutation M { updateMapItem(id: ITEM_ID, quote: "Updated quote") { ... }}
mutation M { updateMap(id: MAP_ID, name: "Updated name") { ... }}

# 7. Remove a map item and a map
mutation M { removeMapItem(id: ITEM_ID) { ... }}
mutation M { removeMap(id: MAP_ID) { ... }}
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
) { ... }}
```

#### Storing map settings

This attribute is there in case you need to store some extra data for a map, e.g. if you want to sort your data in a particular way in your UI.

```
mutation M { updateMap(
  id: MAP_ID,

  settings: { ui_sort_by: "author", ui_sort_dir: "asc" }
) { ... }}
```


### Public API

__All__ items for a single map:

```
GET /public/:user_id/:map_name

:user_id, your case-insensitive user id
:map_name, your case-insensitive map name
```

__One__ item for a single map:

```
GET /public/:user_id/:map_name/:map_item_id
```

__Options__

```markdown
GET /public/:user_id/:map_name?sort_by=author

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
echo 'export SECRET_KEY=...' >> .env
echo 'export AUTH0_DOMAIN=...' >> .env
echo 'export AUTH0_CLIENT_ID=...' >> .env

source .env

mix deps.get
mix test

mix ecto.create
mix run --no-halt
```
