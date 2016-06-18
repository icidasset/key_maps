use Mix.Config


check_env_variables = fn(vars) ->
  has_missing_vars = Enum.any? vars, fn(var) ->
    var = System.get_env(var)
    var == nil or String.length(var) === 0
  end

  if has_missing_vars,
    do: throw("One or more ENV variables is missing")
end


# pre-flight check
check_env_variables.([
  "SECRET_KEY"
])


# config
config :guardian, Guardian,
  issuer: "KeyMaps",
  ttl: { 365, :days },
  verify_issuer: true,
  secret_key: System.get_env("SECRET_KEY"),
  serializer: KeyMaps.Guardian.Serializer


config :key_maps, ecto_repos: [KeyMaps.Repo]


import_config "#{Mix.env}.exs"
