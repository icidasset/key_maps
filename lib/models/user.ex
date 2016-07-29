defmodule KeyMaps.Models.User do
  use Ecto.Schema

  alias KeyMaps.{Repo, Models}

  import Comeonin.Bcrypt
  import Ecto.Changeset


  schema "users" do
    field :email, :string
    field :password, :string, virtual: true
    field :password_hash, :string
    field :username, :string

    has_many :maps, Models.Map

    timestamps
  end


  def changeset(user, params) do
    user
    |> cast(params, ~w(email password_hash username)a)
    |> validate_required(~w(email password_hash username)a)
    |> validate_format(:email, ~r/@/)
    |> validate_length(:password, min: 5)
    |> validate_length(:username, min: 2)
    |> unique_constraint(:email)
    |> unique_constraint(:username)
  end


  #
  # Queries
  #
  def get_by_email(email) do
    Repo.get_by(Models.User, email: email)
  end


  def get_by_username(username) do
    Repo.get_by(Models.User, username: username)
  end


  def create(attr) do
    attr = Map.put(attr, :password_hash, hashpwsalt(attr.password))

    # insert
    if is_nil(System.get_env("ENABLE_SIGN_UP")) === false,
      do: Repo.insert changeset(%Models.User{}, attr)
  end

end
