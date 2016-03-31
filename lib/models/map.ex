defmodule KeyMaps.Models.Map do
  use Ecto.Schema

  alias KeyMaps.{Repo, Models}

  import Ecto.Changeset
  import Ecto.Query


  schema "maps" do
    field :name, :string
    field :attributes, { :array, :string }
    field :user_id, :integer

    timestamps
  end


  def changeset(user, params \\ :empty) do
    user
    |> cast(params, ~w(name attributes user_id), ~w())
    |> unique_constraint(:name)
  end


  #
  # Queries
  #
  def all(params, _, _) do
    Repo.all(from m in Models.Map, where: m.user_id == ^params.user_id)
  end


  def get(_, %{ id: id }, _),   do: Models.Map |> Repo.get_by(id: id)
  def create(params, attr, _),  do: insert(attr, params.user_id)


  #
  # Private
  #
  defp insert(attr, user_id) do
    attr = %{ user_id: user_id } |> Map.merge(attr)

    case Repo.insert changeset(%Models.Map{}, attr) do
      { :ok, map } -> map
      { :error, changeset } ->
        raise GraphQL.CustomError,
          message: KeyMaps.Utils.get_error_from_changeset(changeset),
          status: 422
    end
  end

end
