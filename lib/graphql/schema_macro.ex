defmodule KeyMaps.GraphQL.SchemaMacro do
  alias Ectograph.Definitions, as: D
  alias KeyMaps.Models.{Map, MapItem}

  defmacro build_schema do
    Macro.escape %GraphQL.Schema{

      query: %GraphQL.Type.ObjectType{
        name: "Queries",
        description: "Key Maps API Queries",
        fields: %{
          map:        D.build(Map,      :get, ~w(id name)a    ),
          maps:       D.build(Map,      :all, ~w()a           ),
          mapItem:    D.build(MapItem,  :get, ~w(id)a         ),
          mapItems:   D.build(MapItem,  :all, ~w(map map_id)a ) |> add_map_attr,
        }
      },

      mutation: %GraphQL.Type.ObjectType{
        name: "Mutations",
        description: "Key Maps API Mutations",
        fields: %{
          createMap:      D.build(Map,      :create, ~w(name attributes)a   ),
          createMapItem:  D.build(MapItem,  :create, ~w(map map_id)a        ) |> add_map_attr,
          updateMap:      D.build(Map,      :update, ~w(id name attributes)a),
          updateMapItem:  D.build(MapItem,  :update, ~w(id)a                ),
          removeMap:      D.build(Map,      :delete, ~w(id name)a           ),
          removeMapItem:  D.build(MapItem,  :delete, ~w(id)a                ),
        },
      }

    }
  end


  defp add_map_attr(d) do
    D.extend_arguments(d, %{ map: %{ type: %GraphQL.Type.String{} }})
  end

end
