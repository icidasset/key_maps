# Key Maps

Key maps provides an interface to fill in custom data structures and a JSON api with CORS support to fetch the data.

__Work in progress.__

_You can find old code in the legacy branches._



## How it works

### Authentication

```markdown
POST  /sign-up
GET   /sign-in

__Request body:__

{
  "email": "...",
  "password": "..."
}

__Response body:__

{
  "token": "..."
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

```
mutation M createMap(
  name: "Quotes",
  attributes: [ "quote", "author" ]
)

mutation M createMapItem(
  map: "Quotes",

  quote: "Specialization tends to shut off the wide-band tuning searches and thus to preclude further discovery.",
  author: "Buckminster Fuller"
) {
  attributes
}

query Q mapItems(map: "Quotes") {
  attributes
}

query Q maps() {
  name,
  attributes
}

mutation M removeMapItem(map: "Quotes", id: ITEM_ID)
mutation M removeMap(name: "Quotes")
```

__Notes__  
The map name must be unique, it will be casted to lowercase for validation.


### Public API

__TODO__



## Development

```
source .env

mix deps.get

mix ecto.create MIX_ENV=test
mix test

mix ecto.create
mix run --no-halt
```
