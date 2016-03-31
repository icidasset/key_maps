defmodule KeyMaps.GraphQL.Definitions do
  alias GraphQL.Type.{ID, List}

  def build(model, :all) do
    %{
      type: %List{ ofType: build_type(model) },
      resolve: &model.all/3,
    }
  end


  def build(model, :get, attributes) do
    type_def = build_type(model)

    %{
      type: type_def,
      args: pick_types(type_def, attributes),
      resolve: &model.get/3,
    }
  end


  def build(model, :create, attributes) do
    type_def = build_type(model)

    %{
      type: type_def,
      args: pick_types(type_def, attributes),
      resolve: &model.create/3,
    }
  end


  #
  # Private
  #

  defp build_type(model) do
    cast = Ectograph.Schema.cast_schema(model, :ecto_to_graphql)

    if elem(cast, 0) == :ok do
      elem(cast, 1)
    else
      raise "Could not cast Ecto schema `" <> model.__schema__(:source) <> "`"
    end
  end


  @doc """
    Pick specific fields from a type definition.

    # example:

    pick_types(
      %{ name: "Whatever", fields: %{ a: %{ type: ... }, b: %{ type: ... } }},
      [:a]
    )

    -> %{ a: %{ type: ... } }
  """
  defp pick_types(type_def, keys) do
    mapper = fn(f) ->
      k = elem(f, 0)
      v = elem(f, 1)

      {k, %{ type: v }}
    end

    m = type_def
      |> Map.fetch!(:fields)
      |> Map.take(keys)
      |> Enum.map(mapper)
      |> Map.new()

    m
  end

end
