defmodule KeyMaps do
  use Application

  def start(_type, _args) do
    import Supervisor.Spec, warn: false

    children = [
      worker(KeyMaps.Repo, []),
      worker(KeyMaps.Router, [])
    ]

    Supervisor.start_link(
      children,
      strategy: :one_for_one,
      name: KeyMaps.Supervisor
    )
  end

end
