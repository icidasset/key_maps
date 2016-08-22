defmodule KeyMaps.Repo.Migrations.CreateMapItems do
  use Ecto.Migration

  def change do
    execute "CREATE EXTENSION IF NOT EXISTS hstore"

    create table(:map_items) do
      add :attributes, :map
      add :map_id, references(:maps)

      timestamps
    end
  end
end
