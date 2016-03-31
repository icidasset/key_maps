# Key Maps

Key maps provides an interface to fill in custom data structures and a JSON api with CORS support to fetch the data.

__Work in progress.__

_You can find old code in the legacy branches._



## How it works

To be able to create maps, you must first authenticate yourself:

```
POST /sign-up { email: "...", password: "..." }
```

This will return a token you can use to authenticate your requests.
__Same goes for `sign-in`.__

You create map, which has a name/key and a set of attributes.  
For example:

```graphql
mutation M createMap(
  name: "Quotes",
  attributes: { quote: "string", author: "string" }
)
```

__Note:__ The map name must be unique, it will be casted to
lowercase for validation.

We now have a "repository" for our quotes,
and we can add one by, for example, executing this query:

```graphql
mutation M createMapItem(
  quote: "Specialization tends to shut off the wide-band tuning searches and thus to preclude further discovery.",
  author: "Buckminster Fuller"
)
```

Thus far we have one quote in our repository.
Let's fetch it from the api.

```graphql
query Q mapItem(map: "Quotes", id: ITEM_HASH) { quote, author }
```

Great! But what if we would like to see all quotes?

```graphql
query Q mapItems(map: "Quotes") { quote, author }
```

__That's it!__

PS. This is how you remove map items and maps:

```
mutation M removeMapItem(map: "Quotes", id: ITEM_HASH)
mutation M removeMap(name: "Quotes")
```
