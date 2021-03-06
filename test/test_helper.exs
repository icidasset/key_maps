ExUnit.start

Mix.Task.run "ecto.create", ["--quiet"]
Mix.Task.run "ecto.migrate", ["--quiet"]


defmodule KeyMaps.TestHelpers do
  use Plug.Test

  alias KeyMaps.{Router}


  def request(method, path, attr, token \\ nil) do
    params = if is_binary(attr), do: attr
    path = if is_binary(attr), do: path, else: path <> "?" <> Plug.Conn.Query.encode(attr)

    conn(method, path, params)
      |> put_req_header( "accept", "application/json" )
      |> put_req_header( "content-type", "application/json" )
      |> put_req_header( "authorization", token || "" )
      |> Router.call( Router.init([]) )
  end


  def request_with_json_body(method, path, map, token \\ nil) do
    request(method, path, Poison.encode!(map), token)
  end


  def response(conn, key) do
    Poison.decode!(conn.resp_body)[key]
  end


  def data_response(conn), do: response(conn, "data")
  def error_response(conn), do: response(conn, "errors") |> List.first


  def graphql_request(type, name, attr), do: graphql_request(type, name, %{}, attr, nil)
  def graphql_request(type, name, attr, token), do: graphql_request(type, name, %{}, attr, token)


  def graphql_request(type, name, args, attr, token) do
    query = format_graphql_query(type, name, args, attr)
    do_graphql_request(type, token, query)
  end


  def do_graphql_request(type, token, query) do
    method = if type === :mutation, do: :post, else: :get
    request(method, "/api", %{ query: query }, token)
  end


  #
  # Private
  #
  defp format_graphql_query(type, name, args, attr) do
    type = Atom.to_string(type)
    name = Atom.to_string(name)
    args = Map.to_list(args)
    attr = Enum.join(attr, ",")
    id = type |> String.at(0) |> String.upcase

    args = Enum.map args, fn(arg) ->
      k = elem(arg, 0)
      v = Poison.encode!(elem(arg, 1))

      "#{k}: #{v}"
    end

    args = if length(args) > 0,
      do: "(" <> Enum.join(args, ", ") <> ")",
    else: ""

    "#{type} #{id} { #{name} #{args} { #{attr} } }"
  end

end
