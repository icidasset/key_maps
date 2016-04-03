defmodule KeyMaps.Public.Processor do
  alias KeyMaps.Models

  def run(map, opts \\ %{}) do
    map = Models.Map.load_map_items_into(map)

    # data
    include_timestamps = Map.has_key?(opts, "timestamps")

    map_items = Enum.map map.map_items, fn(m) ->
      attr = m.attributes

      if include_timestamps do
        attr = Map.put(attr, :inserted_at, m.inserted_at)
        attr = Map.put(attr, :updated_at, m.updated_at)
      end

      attr
    end

    # sort
    if Map.has_key?(opts, "sort_by") do
      first_map_item = List.first(map_items)

      the_sort_key = opts["sort_by"]
      has_sort_key = if first_map_item, do: Map.has_key?(first_map_item, the_sort_key)

      sort_direction = Map.get(opts, "sort_direction", "asc")
      sort_direction = String.to_atom(sort_direction)

      if has_sort_key do
        map_items = Enum.sort_by(
          map_items,
          ( fn(m) -> Map.get(m, the_sort_key) end ),
          ( if sort_direction == :desc do &>=/2 else &<=/2 end )
        )
      end
    end

    # return
    map_items
  end

end
