defmodule KeyMaps.GraphQL.Schema do
  alias Ectograph.Definitions, as: D
  alias KeyMaps.{Models.Map, Models.MapItem}

  def schema do
    %GraphQL.Schema{

      query: %GraphQL.Type.ObjectType{
        name: "Queries",
        description: "Key Maps API Queries",
        fields: %{
          map:        D.build(Map,      :get, ~w(name)a   ),
          maps:       D.build(Map,      :all, ~w()a       ),
          mapItem:    D.build(MapItem,  :get, ~w(id)a     ),
          mapItems:   D.build(MapItem,  :all, ~w(map)a    ) |> add_map_attr,
        }
      },

      mutation: %GraphQL.Type.ObjectType{
        name: "Mutations",
        description: "Key Maps API Mutations",
        fields: %{
          createMap:      D.build(Map,      :create, ~w(name attributes)a   ),
          createMapItem:  D.build(MapItem,  :create, ~w(map)a               ) |> add_map_attr,
          removeMap:      D.build(Map,      :delete, ~w(name)a              ),
          removeMapItem:  D.build(MapItem,  :delete, ~w(id)a                ),
        },
      }

    }
  end


  defp add_map_attr(d) do
    D.extend_arguments(d, %{ map: %{ type: %GraphQL.Type.String{} }})
  end

end
