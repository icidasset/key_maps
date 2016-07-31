defmodule RouterTest do
  use ExUnit.Case, async: true
  use Plug.Test

  import KeyMaps.TestHelpers

  alias KeyMaps.{Models}

  @user_default %{ email: "default@email.com", password: "test-default", username: "default" }
  @user_auth %{ email: "auth@email.com", password: "test-auth", username: "auth" }


  setup_all do
    :ok = Ecto.Adapters.SQL.Sandbox.checkout(KeyMaps.Repo)

    # setup db
    Ecto.Adapters.SQL.Sandbox.mode(KeyMaps.Repo, { :shared, self() })

    # create test user
    { :ok, user } = Models.User.create(@user_default)
    { :ok, token, _ } = Guardian.encode_and_sign(user)
    params = %{ user_id: user.id }

    # pre-build map
    map_attributes = %{ name: "Quotes", attributes: ["quote", "author"] }
    map = Models.Map.create(params, map_attributes, nil)

    # pre-build map item – 1
    map_item_attributes = %{ map: map_attributes.name, quote: "1st Q", author: "1st A" }
    map_item_1 = Models.MapItem.create(map, map_item_attributes)

    # pre-build map item – 2 – created a second later
    t = Task.async fn ->
      receive do
        :ok ->
          map_item_attributes = %{ map: map_attributes.name, quote: "2nd Q", author: "2nd A" }
          Models.MapItem.create(map, map_item_attributes)
      end
    end

    Process.send_after(t.pid, :ok, 1500)
    Task.await(t)

    # --> share data with tests
    { :ok, %{
      map: map,
      map_item_1: map_item_1,
      token: token,
      user_id: user.id }}
  end



  #
  # USERS / AUTHENTICATION
  #

  @tag :users
  test "users -- sign up and in" do
    conn = request_with_json_body(:post, "/sign-up", @user_auth)
    token = data_response(conn)["token"]

    # assert
    assert conn.status == 201
    assert token
    assert String.length(token) > 0

    # --- sign in with email
    params = %{ login: @user_auth.email, password: @user_auth.password }
    conn = request_with_json_body(:post, "/sign-in", params)
    token = data_response(conn)["token"]

    # assert
    assert conn.status == 200
    assert token
    assert String.length(token) > 0

    # --- sign in with username
    params = %{ login: @user_auth.username, password: @user_auth.password }
    conn = request_with_json_body(:post, "/sign-in", params)
    token = data_response(conn)["token"]

    # assert
    assert conn.status == 200
    assert token
    assert String.length(token) > 0
  end


  @tag :users
  test "users -- should have a unique email" do
    attr = Map.put(@user_default, :username, "something")
    conn = request_with_json_body(:post, "/sign-up", attr)

    # assert
    assert conn.status == 400
  end


  @tag :users
  test "users -- should have a unique username" do
    attr = Map.put(@user_default, :email, "other-email@example.com")
    conn = request_with_json_body(:post, "/sign-up", attr)

    # assert
    assert conn.status == 400
  end


  @tag :users
  test "users -- should be authenticated for graphql queries (ie. /api)" do
    conn = graphql_request(:query, :maps, ~w(name))
    message = error_response(conn)["message"]

    # assert
    assert conn.status == 403
    assert message == "Forbidden"
  end



  #
  # MAPS
  #

  @tag :maps
  test "maps -- create", context do
    conn = do_graphql_request(
      :mutation,
      context.token,
      """
      mutation M {
        createMap(
          name: "Test",
          attributes: [ "example" ],
          types: { example: "string", test: 1 },
          settings: { test: 1 }
        ) {
          name,
          attributes,
          types,
          settings
        }
      }
      """
    )

    # data
    data = data_response(conn)["createMap"]

    # assert
    assert conn.status == 200
    assert List.first(data["attributes"]) == "example"
    assert data["types"]["example"] == "string"
    assert data["types"]["test"] == 1
    assert data["settings"]["test"] == 1
  end


  @tag :maps
  test "maps -- create -- name should be unique (case insensitive)", context do
    try do
      graphql_request(
        :mutation,
        :createMap,
        %{ name: "quotes", attributes: ["something"] },
        ~w(name),
        context.token
      )
    rescue
      err -> assert err.status == 422
    end
  end


  @tag :maps
  test "maps -- create -- name should only be unique per user" do
    user_attr = %{
      email: "maps-create-unique@email.com",
      password: "test-maps-create",
      username: "mcu"
    }

    { :ok, user } = Models.User.create(user_attr)
    { :ok, token, _ } = Guardian.encode_and_sign(user)

    # make map
    conn = graphql_request(
      :mutation,
      :createMap,
      %{ name: "Quotes", attributes: ["something"] },
      ~w(name),
      token
    )

    # assert
    assert conn.status == 200
  end


  @tag :maps
  test "maps -- create -- should have attributes", context do
    conn = graphql_request(
      :mutation,
      :createMap,
      %{ name: "Test - MHA", attributes: [] },
      ~w(name),
      context.token
    )

    # assert
    assert error_response(conn)["message"] =~ "at least 1 item"
  end


  @tag :maps
  test "maps -- create -- should have valid attributes", context do
    conn = graphql_request(
      :mutation,
      :createMap,
      %{ name: "Test - MHVA", attributes: [0, 1, 2] },
      ~w(attributes),
      context.token
    )

    # assert
    assert error_response(conn)["message"] =~ "is invalid"
  end


  @tag :maps
  test "maps -- create -- should sluggify attributes", context do
    conn = graphql_request(
      :mutation,
      :createMap,
      %{ name: "Test - MHVA", attributes: ["must be slugged"] },
      ~w(attributes),
      context.token
    )

    # assert
    assert conn.status == 200
    assert List.first(data_response(conn)["createMap"]["attributes"]) != "must be slugged"
    assert List.first(data_response(conn)["createMap"]["attributes"]) == "must-be-slugged"
  end


  @tag :maps
  test "maps -- get", context do
    conn = graphql_request(:query, :map, %{ name: "Quotes" }, ~w(name), context.token)

    # assert
    assert data_response(conn)["map"]["name"] == "Quotes"
  end


  @tag :maps
  test "maps -- all", context do
    conn = graphql_request(:query, :maps, ~w(name), context.token)

    # response
    map_item = data_response(conn)["maps"] |> List.first

    # assert
    assert map_item["name"] == "Quotes"
  end


  @tag :maps
  test "maps -- update", context do
    map_attributes = %{ name: "ZZZ - Update test", attributes: ["something"] }
    map = Models.Map.create(context, map_attributes, nil)
    new_name = String.replace(map.name, "test", "test success")

    conn = graphql_request(
      :mutation,
      :updateMap,
      %{ id: map.id, name: new_name },
      ~w(name),
      context.token
    )

    # assert
    map = Models.Map.get(context, %{ id: map.id }, nil)

    assert conn.status == 200
    assert map.name == "ZZZ - Update test success"
  end


  @tag :maps
  test "maps -- remove (by name)", context do
    map_attributes = %{ name: "ZZZ - Remove test", attributes: ["something"] }
    map = Models.Map.create(context, map_attributes, nil)

    conn = graphql_request(
      :mutation,
      :removeMap,
      %{ name: map.name },
      ~w(name),
      context.token
    )

    # assert
    assert conn.status == 200
    assert Models.Map.get(context, %{ name: map.name }, nil) == nil
  end


  @tag :maps
  test "maps -- remove (by id)", context do
    map_attributes = %{ name: "ZZZ - Remove test", attributes: ["something"] }
    map = Models.Map.create(context, map_attributes, nil)

    conn = graphql_request(
      :mutation,
      :removeMap,
      %{ id: map.id },
      ~w(name),
      context.token
    )

    # assert
    assert conn.status == 200
    assert Models.Map.get(context, %{ id: map.id }, nil) == nil
  end



  #
  # MAP ITEMS
  #

  @tag :map_items
  test "map items -- create", context do
    conn = graphql_request(
      :mutation,
      :createMapItem,
      %{ map: "Quotes", quote: "A", author: "B" },
      ~w(attributes),
      context.token
    )

    # response
    map_item = data_response(conn)["createMapItem"]

    # assert
    assert conn.status == 200
    assert map_item["attributes"]["quote"] == "A"
    assert map_item["attributes"]["author"] == "B"
  end


  @tag :map_items
  test "map items -- create -- should filter other attributes", context do
    conn = graphql_request(
      :mutation,
      :createMapItem,
      %{ map: "Quotes", quote: "A", shouldNotBeHere: true },
      ~w(attributes),
      context.token
    )

    # response
    map_item = data_response(conn)["createMapItem"]

    # assert
    assert conn.status == 200
    assert map_item["attributes"]["quote"] == "A"
    assert map_item["attributes"]["shouldNotBeHere"] == nil
  end


  @tag :map_items
  test "map items -- all", context do
    conn = graphql_request(
      :query,
      :mapItems,
      %{ map: "Quotes" },
      ~w(id attributes),
      context.token
    )

    # response
    map_item = data_response(conn)["mapItems"]
      |> List.first

    # assert
    assert conn.status == 200
    assert map_item["id"]
    assert map_item["attributes"]
  end


  @tag :map_items
  test "map items -- get", context do
    conn = graphql_request(
      :query,
      :mapItem,
      %{ id: context.map_item_1.id },
      ~w(id attributes),
      context.token
    )

    # response
    map_item = data_response(conn)["mapItem"]

    # assert
    assert conn.status == 200
    assert map_item["id"] == Integer.to_string(context.map_item_1.id)
    assert map_item["attributes"]
  end


  @tag :map_items
  test "map items -- update", context do
    map_item_attributes = %{ map: context.map.name, quote: "Z", author: "Z" }
    map_item = Models.MapItem.create(context.map, map_item_attributes)

    conn = graphql_request(
      :mutation,
      :updateMapItem,
      %{ id: map_item.id, quote: "Updated" },
      ~w(id),
      context.token
    )

    # assert
    assert conn.status == 200
    assert KeyMaps.Repo.get_by(Models.MapItem, id: map_item.id).attributes["quote"] == "Updated"
    assert KeyMaps.Repo.get_by(Models.MapItem, id: map_item.id).attributes["author"] == "Z"
  end


  @tag :map_items
  test "map items -- remove", context do
    map_item_attributes = %{ map: context.map.name, quote: "Z", author: "Z" }
    map_item = Models.MapItem.create(context.map, map_item_attributes)

    conn = graphql_request(
      :mutation,
      :removeMapItem,
      %{ id: map_item.id },
      ~w(id),
      context.token
    )

    # assert
    assert conn.status == 200
    assert KeyMaps.Repo.get_by(Models.MapItem, id: map_item.id) == nil
  end



  #
  # PUBLIC
  #

  @tag :public
  test "public -- map" do
    conn = request(:get, "/public/default/quotes", nil)

    data = data_response(conn)
    data = Enum.filter data, fn(d) ->
      if Map.has_key?(d, "author"),
        do: String.last(d["author"]) === "A"
    end

    first_item = Enum.at(data, 0)
    second_item = Enum.at(data, 1)

    # assert
    assert conn.status == 200
    assert is_list(data)

    assert first_item["quote"] == "2nd Q"
    assert first_item["author"] == "2nd A"

    assert second_item["quote"] == "1st Q"
    assert second_item["author"] == "1st A"
  end


  test "public -- map -- sorted" do
    conn = request(:get, "/public/default/quotes", %{ sort_by: "author" })

    data = data_response(conn)
    data = Enum.filter data, fn(d) ->
      if Map.has_key?(d, "author"),
        do: String.last(d["author"]) === "A"
    end

    first_item = Enum.at(data, 0)

    # assert
    assert conn.status == 200
    assert first_item["quote"] == "1st Q"
    assert first_item["author"] == "1st A"
  end


  test "public -- map -- timestamps" do
    conn = request(:get, "/public/default/quotes", %{ timestamps: "1" })

    data = data_response(conn)
    data = Enum.filter data, fn(d) ->
      if Map.has_key?(d, "author"),
        do: String.last(d["author"]) === "A"
    end

    first_item = Enum.at(data, 0)

    # assert
    assert conn.status == 200
    assert first_item["inserted_at"]
    assert first_item["updated_at"]
  end


  test "public -- map -- single item", context do
    id = Integer.to_string(context.map_item_1.id)
    conn = request(:get, "/public/default/quotes/" <> id, %{})
    data = data_response(conn)

    # assert
    assert conn.status == 200
    assert data["quote"] == "1st Q"
  end

end
