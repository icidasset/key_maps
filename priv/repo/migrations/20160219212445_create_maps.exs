defmodule KeyMaps.Repo.Migrations.CreateMaps do
  use Ecto.Migration

  def change do
    create table(:maps) do
      add :key, :string
      add :name, :string

      timestamps
    end
  end
end
