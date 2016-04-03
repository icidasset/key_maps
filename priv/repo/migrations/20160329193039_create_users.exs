defmodule KeyMaps.Repo.Migrations.CreateUsers do
  use Ecto.Migration

  def change do
    execute "CREATE EXTENSION IF NOT EXISTS citext"

    create table(:users) do
      add :email, :citext
      add :password_hash, :string
      add :username, :citext

      timestamps
    end

    create unique_index :users, [:email]
    create unique_index :users, [:username]
  end
end
