defmodule KeyMaps do
  use Application

  def start(_type, _args) do
    import Supervisor.Spec

    children = [
      worker(KeyMaps.Repo, []),
      worker(KeyMaps.Api, [])
    ]

    Supervisor.start_link(children, strategy: :one_for_one)
  end

end
