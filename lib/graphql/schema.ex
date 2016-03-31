defmodule KeyMaps.GraphQL.Schema do
  alias KeyMaps.GraphQL.{Definitions}
  alias KeyMaps.{Models.Map, Models.MapItem}

  def schema do
    %GraphQL.Schema{

      query: %GraphQL.Type.ObjectType{
        name: "Queries",
        description: "Key Maps API Queries",
        fields: %{
          map: Definitions.build(Map, :get, ~w(id)a),
          maps: Definitions.build(Map, :all),
          mapItems: Definitions.build(MapItem, :all),
        }
      },

      mutation: %GraphQL.Type.ObjectType{
        name: "Mutations",
        description: "Key Maps API Mutations",
        fields: %{
          createMap: Definitions.build(Map, :create, ~w(name attributes)a),
          createMapItem: Definitions.build(MapItem, :create, ~w(map_id attributes)a),
        },
      }

    }
  end

end
