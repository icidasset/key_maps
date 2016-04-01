ExUnit.start

Mix.Task.run "ecto.create", ["--quiet"]
Mix.Task.run "ecto.migrate", ["--quiet"]
Ecto.Adapters.SQL.begin_test_transaction(KeyMaps.Repo)


defmodule KeyMaps.TestHelpers do
  use Plug.Test

  alias KeyMaps.{Router}


  def request(method, path, attributes, token \\ nil) do
    conn(method, path, attributes)
      |> put_req_header( "accept", "application/json" )
      |> put_req_header( "content-type", "application/json" )
      |> put_req_header( "authorization", token || "" )
      |> Router.call( Router.init([]) )
  end


  def response(conn, key) do
    Poison.decode!(conn.resp_body)[key]
  end


  def data_response(conn), do: response(conn, "data")
  def error_response(conn), do: response(conn, "errors")


  def graph_query_request(query, token \\ nil) do
    request(:get, "/api", %{ query: "query Q { #{query} }" }, token)
  end


  def graph_mutation_request(query, token \\ nil) do
    request(:post, "/api", %{ query: "mutation M { #{query} }" }, token)
  end

end
