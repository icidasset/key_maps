defmodule KeyMaps.Models.Map do
  use Ecto.Schema

  alias KeyMaps.{Repo, Models}

  import Ecto.Changeset
  import Ecto.Query, only: [from: 1, from: 2]


  schema "maps" do
    field :name, :string
    field :attributes, { :array, :string }

    belongs_to :user, Models.User
    has_many :map_items, Models.MapItem

    timestamps
  end


  def changeset(user, params) do
    user
    |> cast(params, ~w(name attributes user_id)a)
    |> validate_required(~w(name attributes user_id)a)
    |> unique_constraint(:name, name: :maps_name_user_id_index)
    |> validate_length(:attributes, min: 1)
    |> update_attributes(:attributes)
  end


  #
  # {field} Attributes
  #
  def update_attributes(changeset, field) do
    update_change changeset, field, fn(attributes) ->
      Enum.map(attributes, &Slugger.slugify_downcase/1)
    end
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
    args = %{ user_id: params.user_id } |> Map.merge(args)

    # insert
    case Repo.insert changeset(%Models.Map{}, args) do
      { :ok, map } -> map
      { :error, changeset } -> raise KeyMaps.Utils.get_error_from_changeset(changeset)
    end
  end


  def update(params, args, _) do
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

end
