defmodule KeyMaps.Public.Processor do
  alias KeyMaps.Models

  def run(map, opts \\ %{}) do
    type = if Map.has_key?(opts, "map_item_id"),
      do: :single,
    else: :multiple

    # load map items
    map = case type do
      :single -> Models.Map.load_map_item_into(map, opts["map_item_id"])
      :multiple -> Models.Map.load_map_items_into(map)
    end

    # data
    include_timestamps = Map.has_key?(opts, "timestamps")

    map_items = Enum.map map.map_items, fn(m) ->
      if include_timestamps do
        m.attributes
        |> Map.put(:inserted_at, m.inserted_at)
        |> Map.put(:updated_at, m.updated_at)
      else
        m.attributes
      end
    end

    # to be continued
    do_run(type, map_items, opts)
  end


  defp do_run(:single, map_items, _) do
    Enum.at(map_items, 0)
  end


  defp do_run(:multiple, map_items, opts) do
    if Map.has_key?(opts, "sort_by") do
      first_map_item = List.first(map_items)
      the_sort_key = opts["sort_by"]

      sort_direction = Map.get(opts, "sort_direction", "asc")
      sort_direction = String.to_atom(sort_direction)

      sort_method = if sort_direction == :desc,
          do: &>=/2,
        else: &<=/2

      Enum.sort_by(
        map_items,
        ( fn(m) -> Map.get(m, the_sort_key) || "" end ),
        ( sort_method )
      )

    else
      map_items

    end
  end

end
