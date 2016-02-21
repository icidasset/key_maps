defmodule KeyMaps.Models.Map do
  use Ecto.Schema

  alias KeyMaps.{Repo, Models}

  schema "maps" do
    field :key, :string
    field :name, :string

    timestamps
  end


  def all(params, _, _) do
    # IO.inspect(params.authToken)
    Models.Map |> Repo.all
  end


  def get(_, %{ id: id }, _) do
    Models.Map |> Repo.get_by(id: id)
  end


  def create(_, args, _) do
    case Repo.insert %Models.Map{ name: args.name } do
      { :ok, model } -> IO.inspect(model)
      { :error, changeset } -> IO.inspect(changeset)
    end
  end

end
