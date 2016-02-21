defmodule KeyMaps.Api do
  use Plug.Router

  plug Plug.Parsers,
    parsers: [:urlencoded, :json],
    json_decoder: Poison

  plug :match
  plug :dispatch


  def start_link do
    { :ok, _ } = Plug.Adapters.Cowboy.http KeyMaps.Api, [], port: 8080
  end


  forward "/api",
    to: GraphQL.Plug.Endpoint,
    schema: { KeyMaps.GraphQL.Schema, :schema }

end
