use Mix.Config


config :guardian, Guardian,
  issuer: "KeyMaps",
  ttl: { 365, :days },
  verify_issuer: true,
  secret_key: System.get_env("SECRET_KEY"),
  serializer: KeyMaps.Guardian.Serializer


import_config "#{Mix.env}.exs"
