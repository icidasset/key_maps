defmodule RouterTest do
  use ExUnit.Case
  use Plug.Test

  import KeyMaps.TestHelpers

  alias KeyMaps.{Models}

  @user_default %{ email: "default@email.com", password: "test-default" }
  @user_auth %{ email: "auth@email.com", password: "test-auth" }


  setup_all do
    { :ok, user } = Models.User.create(@user_default)
    { :ok, token, _ } = Guardian.encode_and_sign(user)
    params = %{ user_id: user.id }

    # prebuild map
    map_attributes = %{ name: "Quotes", attributes: ["quote", "author"] }
    map = Models.Map.create(params, map_attributes, nil)

    # --> share data with tests
    { :ok, %{ token: token, map: map } }
  end


  test "sign up and in" do
    conn = request(:post, "/sign-up", Poison.encode!(@user_auth))
    token = data_response(conn)["token"]

    # assert
    assert conn.status == 201
    assert token
    assert String.length(token) > 0

    # --- sign in
    conn = request(:post, "/sign-in", Poison.encode!(@user_auth))
    token = data_response(conn)["token"]

    # assert
    assert conn.status == 200
    assert token
    assert String.length(token) > 0
  end


  test "must be authenticated for graphql queries (ie. /api)" do
    conn = graph_query_request("maps { name }")
    message = List.first(error_response(conn))["message"]

    # assert
    assert conn.status == 403
    assert message == "Forbidden"
  end


  test "maps - create", context do
    map_attr = "name: \"Test\", attributes: [\"example\"]"
    conn = graph_mutation_request("createMap (#{map_attr}) { name }", context.token)

    # assert
    assert conn.status == 200
  end


  test "maps - create - name must be unique (case insensitive)", context do
    try do
      map_attr = "name: \"quotes\", attributes: [\"something\"]"
      graph_mutation_request("createMap (#{map_attr}) { name }", context.token)
    rescue
      err -> assert err.status == 422
    end
  end


  test "maps - create - name must only be unique per user" do
    attr = %{ email: "maps-create-unique@email.com", password: "test-maps-create" }

    { :ok, user } = Models.User.create(attr)
    { :ok, token, _ } = Guardian.encode_and_sign(user)

    # make map
    map_attr = "name: \"Quotes\", attributes: [\"something\"]"
    conn = graph_mutation_request("createMap (#{map_attr}) { name }", token)

    # assert
    assert conn.status == 200
  end


  test "maps - get", context do
     conn = graph_query_request("map (name: \"Quotes\") { name }", context.token)

     # assert
     assert data_response(conn)["map"]["name"] == "Quotes"
  end


  test "maps - all", context do
     conn = graph_query_request("maps { name }", context.token)
     map_item = data_response(conn)["maps"] |> List.first

     # assert
     assert map_item["name"] == "Quotes"
  end

end
