defmodule KeyMaps.Repo.Migrations.AddTypesToMaps do
  use Ecto.Migration

  def change do
    alter table(:maps) do
      add :types, :map
    end
  end
end
