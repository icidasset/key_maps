defmodule KeyMaps.Models.MapItem do
  use Ecto.Schema

  alias KeyMaps.{Repo, Models}

  import Ecto.Changeset
  import Ecto.Query


  schema "map_items" do
    field :attributes, :map

    belongs_to :map, Models.Map

    timestamps
  end


  def graphql_attributes do
    %{
      map: %{ type: %GraphQL.Type.String{} }
    }
  end


  def changeset(user, params \\ :empty) do
    user
    |> cast(params, ~w(attributes map_id)a)
    |> validate_required(~w(attributes map_id)a)
    |> validate_attributes(:attributes)
  end


  #
  # {field} Attributes
  #
  def validate_attributes(changeset, field) do
    validate_change changeset, field, fn(_, attributes) ->
      cond do
        does_not_have_attributes(attributes) ->
          [{ field, "must not be empty" }]
        has_nil_attribute(attributes) ->
          [{ field, "must have atleast one attribute that's not nil" }]
        true ->
          []
      end
    end
  end


  def does_not_have_attributes(attributes) do
    (attributes |> Map.keys |> length) === 0
  end


  def has_nil_attribute(attributes) do
    Enum.find(Map.values(attributes), &(&1 != nil)) == nil
  end


  #
  # Queries
  #
  def all(params, attr, _) do
    map = Models.Map.get(params, %{ name: attr.map }, nil)

    if map do
      Repo.all(from m in Models.MapItem, where: m.map_id == ^map.id)
    else
      __raise_map_error()
    end
  end


  def create(params, attr, internal) do
    map = Models.Map.get(params, %{ name: attr.map }, nil)
    other = KeyMaps.Utils.extract_other_attributes(internal)

    if map,
      do: __create(other, map),
    else: __raise_map_error()
  end


  #
  # Private
  #
  defp __create(attr, map) do
    attr = Enum.filter attr, fn(a) ->
      key = elem(a, 0) |> Atom.to_string
      Enum.member?(map.attributes, key)
    end

    # add map id
    attr = Enum.into(attr, %{})
    attr = %{ attributes: attr }
    attr = Map.put(attr, :map_id, map.id)

    # insert
    case Repo.insert changeset(%Models.MapItem{}, attr) do
      { :ok, map_item } -> map_item
      { :error, changeset } ->
        raise GraphQL.CustomError,
          message: KeyMaps.Utils.get_error_from_changeset(changeset),
          status: 422
    end
  end


  defp __raise_map_error do
    raise GraphQL.CustomError, message: "Could not find map", status: 422
  end

end
