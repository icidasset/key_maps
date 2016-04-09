defmodule KeyMaps.Models.MapItem do
  use Ecto.Schema

  alias KeyMaps.{Repo, Models}

  import Ecto.Changeset
  import Ecto.Query, only: [from: 1, from: 2]


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


  def changeset(user, params) do
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
  def all(params, %{ map: map }, _) do
    map = Models.Map.get(params, %{ name: map }, nil)

    if map do
      Repo.all(from m in Models.MapItem, where: m.map_id == ^map.id)
    else
      do_raise_map_error()
    end
  end


  def get(params, %{ id: id }, _) do
    query = from i in Models.MapItem,
      join: m in assoc(i, :map),
      where: i.id == ^id,
      where: m.user_id == ^params.user_id

    Repo.one(query)
  end


  def create(params, %{ map: map }, internal) do
    map = Models.Map.get(params, %{ name: map }, nil)
    other_args = KeyMaps.Utils.extract_other_arguments(internal)

    if map,
      do: do_create(map, other_args),
    else: do_raise_map_error()
  end


  def create(map, args) do
    do_create(map, args)
  end


  def delete(params, args, _) do
    map_item = Models.MapItem.get(params, args, nil)

    if map_item do
      case Repo.delete map_item do
        { :ok, struct } -> struct
        { :error, changeset } ->
          raise GraphQL.CustomError,
            message: KeyMaps.Utils.get_error_from_changeset(changeset),
            status: 422
      end

    else
      do_raise_map_error()

    end
  end


  #
  # Private
  #
  defp do_create(map, args) do
    args = Enum.filter args, fn(a) ->
      key = elem(a, 0) |> Atom.to_string
      Enum.member?(map.attributes, key)
    end

    # add map id
    args = Enum.into(args, %{})
    args = %{ attributes: args }
    args = Map.put(args, :map_id, map.id)

    # insert
    case Repo.insert changeset(%Models.MapItem{}, args) do
      { :ok, map_item } -> map_item
      { :error, changeset } ->
        raise GraphQL.CustomError,
          message: KeyMaps.Utils.get_error_from_changeset(changeset),
          status: 422
    end
  end


  defp do_raise_map_error do
    raise GraphQL.CustomError, message: "Could not find map", status: 404
  end

end
