defmodule KeyMaps.Public.Plug do
  import KeyMaps.Utils

  alias Plug.Conn
  alias KeyMaps.{Models, Public.Processor}

  @behaviour Plug

  def init(_) do
    []
  end


  def call(%Conn{method: m} = conn, _) when m in ["GET"] do
    path = conn.request_path
      |> String.replace_prefix("/public/", "")
      |> String.replace_trailing("/", "")
      |> String.split("/")

    # opts
    opts = %{
      id: Enum.at(path, 0),
      map_name: Enum.at(path, 1),
      map_item_id: Enum.at(path, 2),
    }

    # run
    if opts.id == nil || opts.map_name == nil,
      do: render_error(conn, 422, "Insufficient parameters"),
    else: check_user(conn, opts)
  end


  def call(%Conn{method: _} = conn, _) do
    render_error(conn, 405, "GraphQL only supports GET and POST requests.")
  end


  #
  # Private
  #
  defp check_user(conn, opts) do
    user = Models.User.get_by_id(opts.id)

    if user,
      do: check_map(conn, user.id, opts),
    else: render_error(conn, 422, "Invalid user id")
  end


  defp check_map(conn, user_id, opts) do
    opts = Map.put(opts, :user_id, user_id)
    map = Models.Map.get(opts, %{ name: opts.map_name }, nil)

    if map,
      do: collect_and_render(conn, map, opts),
    else: render_error(conn, 422, "Invalid map name")
  end


  defp collect_and_render(conn, map, opts) do
    processor_options = conn.query_params

    if opts.map_item_id != nil do
      processor_options = Map.put(processor_options, "map_item_id", opts.map_item_id)
    end

    map_items = Processor.run(map, processor_options)

    # render
    render_data(conn, 200, map_items)
  end

end
