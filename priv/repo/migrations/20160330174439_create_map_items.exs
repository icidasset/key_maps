defmodule KeyMaps.Repo.Migrations.CreateMapItems do
  use Ecto.Migration

  def change do
    create table(:map_items) do
      add :attributes, :map
      add :map_id, :integer

      timestamps
    end
  end
end
