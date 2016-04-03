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


  def get(params, %{ name: name }, _) do
    Repo.get_by(Models.Map, name: name, user_id: params.user_id)
  end


  def create(params, args, _) do
    args = %{ user_id: params.user_id } |> Map.merge(args)

    # insert
    case Repo.insert changeset(%Models.Map{}, args) do
      { :ok, map } -> map
      { :error, changeset } ->
        raise GraphQL.CustomError,
          message: KeyMaps.Utils.get_error_from_changeset(changeset),
          status: 422
    end
  end

end
