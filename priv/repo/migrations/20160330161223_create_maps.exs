defmodule KeyMaps.Repo.Migrations.CreateMaps do
  use Ecto.Migration

  def change do
    execute "CREATE EXTENSION IF NOT EXISTS citext"

    create table(:maps) do
      add :name, :citext
      add :attributes, { :array, :string }
      add :user_id, references(:users)

      timestamps
    end

    create unique_index :maps, [:name, :user_id]
    # name: maps_name_user_id_index
  end
end
