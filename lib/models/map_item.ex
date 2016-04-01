defmodule KeyMaps.Models.MapItem do
  use Ecto.Schema

  alias KeyMaps.{Repo, Models}

  import Ecto.Changeset
  import Ecto.Query


  schema "map_items" do
    field :attributes, :map

    field :map, :string
    field :map_id, :integer

    timestamps
  end


  def changeset(user, params \\ :empty) do
    user
    |> cast(params, ~w(attributes map_id), ~w())
    |> validate_attributes(:attributes)
  end


  #
  # {field} Attributes
  #
  def validate_attributes(changeset, field, _ \\ []) do
    validate_change changeset, field, fn(_, attributes) ->
      field = :attributes

      cond do
        count_attributes(attributes) === 0 ->
          [{ field, "must not be empty" }]
        find_non_nil_attribute(attributes) !== nil ->
          [{ field, "must have atleast one attribute that's not nil" }]
        true ->
          []
      end
    end
  end


  def count_attributes(attributes) do
    length(Map.keys(attributes))
  end


  def find_non_nil_attribute(attributes) do
    Enum.find(Map.values(attributes), fn(m) -> m !== nil end)
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


  def create(params, attr, _) do
    map = Models.Map.get(params, %{ name: attr.map }, nil)

    if map,
      do: __create(attr, map.id),
    else: __raise_map_error()
  end


  #
  # Private
  #
  defp __create(attr, map_id) do
    attr = %{
      map_id: map_id,
      attributes: if Map.has_key?(attr, :attributes)
        do Poison.decode!(attr.attributes)
        else nil
      end
    }

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
