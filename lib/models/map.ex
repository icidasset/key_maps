defmodule KeyMaps.Models.Map do
  use Ecto.Schema

  alias KeyMaps.{Repo, Models}

  import Ecto.Changeset
  import Ecto.Query, only: [from: 2]


  schema "maps" do
    field :name, :string
    field :attributes, { :array, :string }
    field :types, :map, default: %{}
    field :settings, :map, default: %{}

    belongs_to :user, Models.User
    has_many :map_items, Models.MapItem, on_delete: :delete_all

    timestamps
  end


  def changeset(user, params) do
    user
    |> cast(params, ~w(name attributes types settings user_id)a)
    |> validate_required(~w(name attributes user_id)a)
    |> unique_constraint(:name, name: :maps_name_user_id_index)
    |> validate_length(:attributes, min: 1)
    |> validate_change(:attributes, &validate_attributes/2)
  end


  def validate_attributes(:attributes, attributes) do
    errors = Enum.map attributes, fn(attr) ->
      if attr && Regex.match?(~r/^\w+$/, attr),
        do: nil,
      else: [attributes: "cannot contain invalid keys" <>
                         " (only alphanumeric characters and underscores)"]
    end

    errors
      |> Enum.reject(&(is_nil(&1)))
      |> List.flatten
  end


  #
  # Queries
  #
  def all(params, _, _) do
    Repo.all(from m in Models.Map, where: m.user_id == ^params.user_id)
  end


  def get(params, args, _) do
    k = cond do
      Map.has_key?(args, :name)   -> [name: args[:name]]
      Map.has_key?(args, :id)     -> [id: args[:id]]
      true                        -> raise "Cannot select map, missing parameters"
    end

    Repo.get_by(Models.Map, [user_id: params.user_id] ++ k)
  end


  def create(params, args, _) do
    args = clean_args(args)
    args = %{ user_id: params.user_id } |> Map.merge(args)

    # insert
    case Repo.insert changeset(%Models.Map{}, args) do
      { :ok, map } -> map
      { :error, changeset } -> raise KeyMaps.Utils.get_error_from_changeset(changeset)
    end
  end


  def update(params, args, _) do
    args = clean_args(args)
    map_args = Map.take(args, [:id])
    map = Models.Map.get(params, map_args, nil)

    if map do
      case Repo.update changeset(map, args) do
        { :ok, map } -> map
        { :error, changeset } -> raise KeyMaps.Utils.get_error_from_changeset(changeset)
      end

    else
      raise "Could not find map"

    end
  end


  def delete(params, args, _) do
    map = Models.Map.get(params, args, nil)

    if map do
      case Repo.delete map do
        { :ok, struct } -> struct
        { :error, changeset } -> raise KeyMaps.Utils.get_error_from_changeset(changeset)
      end

    else
      raise "Could not find map"

    end
  end


  #
  # Associations
  #
  def load_map_items_into(map) do
    query = from(m in Models.MapItem, order_by: [desc: m.inserted_at])
    Repo.preload(map, map_items: query)
  end


  def load_map_item_into(map, item_id) do
    query = from(m in Models.MapItem, where: m.id == ^item_id)
    Repo.preload(map, map_items: query)
  end


  def clean_args(args) do
    if Map.has_key?(args, :attributes) do
      n = cond do
        length(args.attributes) === 0 -> false
        length(Enum.reject(args.attributes, &(is_nil(&1)))) === 0 -> false
        true -> true
      end

      if n == false,
        do: Map.delete(args, :attributes),
      else: args
    else
      args
    end
  end

end
