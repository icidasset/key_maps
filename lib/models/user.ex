defmodule KeyMaps.Models.User do
  use Ecto.Schema
  alias KeyMaps.{Repo, Models}
  import Ecto.Changeset


  schema "users" do
    field :email, :string

    has_many :maps, Models.Map

    timestamps
  end


  def changeset(user, params) do
    user
    |> cast(params, ~w(email)a)
    |> validate_required(~w(email)a)
    |> validate_format(:email, ~r/@/)
    |> unique_constraint(:email)
  end


  #
  # Queries
  #
  def get_by_id(id) do
    Repo.get_by(Models.User, id: id)
  end


  def get_by_email(email) do
    Repo.get_by(Models.User, email: email)
  end


  def create(attr) do
    Repo.insert changeset(%Models.User{}, attr)
  end

end
