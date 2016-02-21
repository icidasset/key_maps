defmodule KeyMaps.GraphQL.Schema do
  alias KeyMaps.GraphQL.{Definitions}
  alias KeyMaps.{Models}

  def schema do
    %GraphQL.Schema{

      query: %GraphQL.Type.ObjectType{
        name: "Queries",
        description: "Key Maps API Queries",
        fields: %{
          map: Definitions.build(Models.Map, :get, [:id]),
          maps: Definitions.build(Models.Map, :all),
        }
      },

      mutation: %GraphQL.Type.ObjectType{
        name: "Mutations",
        description: "Key Maps API Mutations",
        fields: %{
          createMap: Definitions.build(Models.Map, :create, [:name]),
        },
      }

    }
  end

end
