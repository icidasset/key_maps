defmodule KeyMaps.Repo.Migrations.AddSettingsToMaps do
  use Ecto.Migration

  def change do
    alter table(:maps) do
      add :settings, :map
    end
  end
end
