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
  def all(params, args, _) do
    map = do_get_map(params, args)

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


  def create(params, args, internal) do
    map = do_get_map(params, args)
    other_args = KeyMaps.Utils.extract_other_arguments(internal)

    if map,
      do: do_create(map, other_args),
    else: do_raise_map_error()
  end


  def create(map, args) do
    do_create(map, args)
  end


  def update(params, args, internal) do
    map_item = get(params, args, nil)
    other_args = KeyMaps.Utils.extract_other_arguments(internal)

    if map_item,
      do: do_update(map_item, params, other_args),
    else: raise "Could not find map item"
  end


  def delete(params, args, _) do
    map_item = Models.MapItem.get(params, args, nil)

    if map_item do
      case Repo.delete map_item do
        { :ok, struct } -> struct
        { :error, changeset } -> raise KeyMaps.Utils.get_error_from_changeset(changeset)
      end

    else
      do_raise_map_error()

    end
  end


  #
  # Private
  #
  defp do_create(map, args) do
    args = args
      |> do_filter_attributes(map)
      |> do_transform_attributes

    # add map id
    args = Map.put(args, :map_id, map.id)

    # insert
    case Repo.insert changeset(%Models.MapItem{}, args) do
      { :ok, map_item } -> map_item
      { :error, changeset } -> raise KeyMaps.Utils.get_error_from_changeset(changeset)
    end
  end


  defp do_update(map_item, params, args) do
    map = do_get_map(params, %{ map_id: map_item.map_id })

    args = args
      |> do_filter_attributes(map)
      |> do_transform_attributes

    args = put_in(
      args,
      [:attributes],
      Map.merge(map_item.attributes, args.attributes)
    )

    # update
    case Repo.update changeset(map_item, args) do
      { :ok, map_item } -> map_item
      { :error, changeset } -> raise KeyMaps.Utils.get_error_from_changeset(changeset)
    end
  end


  defp do_get_map(params, args) do
    args = cond do
      Map.has_key?(args, :map)    -> %{ name: args[:map] }
      Map.has_key?(args, :map_id) -> %{ id: args[:map_id] }
      true                        -> raise "Cannot select map items, missing map identifier"
    end

    Models.Map.get(params, args, nil)
  end


  defp do_raise_map_error do
    raise "Could not find map"
  end


  defp do_filter_attributes(args, map) do
    Enum.filter args, fn(a) ->
      key = elem(a, 0) |> Atom.to_string
      Enum.member?(map.attributes, key)
    end
  end


  defp do_transform_attributes(args) do
    %{ attributes: Enum.into(args, %{}) }
  end

end
