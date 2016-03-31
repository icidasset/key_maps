defmodule KeyMaps.Repo.Migrations.CreateMaps do
  use Ecto.Migration

  def change do
    create table(:maps) do
      add :name, :string
      add :attributes, { :array, :string }
      add :user_id, references(:users)

      timestamps
    end

    create unique_index :maps, [:name]
  end
end
