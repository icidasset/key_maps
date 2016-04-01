defmodule KeyMaps.Models.Map do
  use Ecto.Schema

  alias KeyMaps.{Repo, Models}

  import Ecto.Changeset
  import Ecto.Query


  schema "maps" do
    field :name, :string
    field :attributes, { :array, :string }

    belongs_to :user, Models.User
    has_many :map_items, Models.MapItem

    timestamps
  end


  def changeset(user, params \\ :empty) do
    user
    |> cast(params, ~w(name attributes user_id), ~w())
    |> unique_constraint(:name, name: :maps_name_user_id_index)
  end


  #
  # Queries
  #
  def all(params, _, _) do
    Repo.all(from m in Models.Map, where: m.user_id == ^params.user_id)
  end


  def get(params, %{ name: name }, _) do
    Repo.get_by(Models.Map, name: name, user_id: params.user_id)
  end


  def create(params, attr, _) do
    attr = %{ user_id: params.user_id } |> Map.merge(attr)

    case Repo.insert changeset(%Models.Map{}, attr) do
      { :ok, map } -> map
      { :error, changeset } ->
        raise GraphQL.CustomError,
          message: KeyMaps.Utils.get_error_from_changeset(changeset),
          status: 422
    end
  end

end
