use Mix.Config


config :key_maps, KeyMaps.Repo,
  adapter: Ecto.Adapters.Postgres,
  url: System.get_env("DATABASE_URL"),
  pool_size: System.get_env("DATABASE_POOL_SIZE"),
  ssl: true
