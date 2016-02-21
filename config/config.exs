use Mix.Config

config :key_maps, KeyMaps.Repo,
  adapter: Ecto.Adapters.Postgres,
  database: "key_maps_development",
  username: "icidasset",
  password: "",
  hostname: "localhost"
