defmodule KeyMaps.GraphQL.Schema do
  alias Ectograph.Definitions, as: D
  alias KeyMaps.Models.{Map, MapItem}

  defmacro build_schema do
    { :ok, map_items_type_def } = Ectograph.Schema.cast(MapItem)

    # schema
    Macro.escape %GraphQL.Schema{

      query: %GraphQL.Type.ObjectType{
        name: "Queries",
        description: "Key Maps API Queries",
        fields: %{
          map:        D.build(Map,      :get, ~w(id name)a),
          maps:       D.build(Map,      :all, ~w()a),
          mapItem:    D.build(MapItem,  :get, ~w(id)a),
          mapItems:   D.build(MapItem,  :all, ~w(map_id)a) |> add_map_attr,
        }
      },

      mutation: %GraphQL.Type.ObjectType{
        name: "Mutations",
        description: "Key Maps API Mutations",
        fields: %{
          createMap:      D.build(Map,      :create, ~w(name attributes types settings)a),
          createMapItem:  D.build(MapItem,  :create, ~w(map_id)a) |> add_map_attr,
          updateMap:      D.build(Map,      :update, ~w(id name attributes types settings)a),
          updateMapItem:  D.build(MapItem,  :update, ~w(id)a),
          removeMap:      D.build(Map,      :delete, ~w(id name)a),
          removeMapItem:  D.build(MapItem,  :delete, ~w(id)a),

          createMapItems: add_items(%{
            type: %GraphQL.Type.List{ ofType: map_items_type_def },
            args: D.pick_types(map_items_type_def, ~w(map_id)a),
            resolve: { MapItem, :create_multiple }
          })
        },
      }

    }
  end


  defp add_map_attr(d) do
    D.extend_arguments(d, %{ map: %{ type: %GraphQL.Type.String{} }})
  end


  defp add_items(d) do
    items_type = %{ type: %GraphQL.Type.String{} }

    d
    |> add_map_attr
    |> D.extend_arguments(%{ items: items_type })
  end


  def schema do
    build_schema
  end

end
