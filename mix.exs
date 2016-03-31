defmodule KeyMaps.Mixfile do
  use Mix.Project

  def project do
    [app: :key_maps,
     version: "0.0.1",
     elixir: "~> 1.2",
     build_embedded: Mix.env == :prod,
     start_permanent: Mix.env == :prod,
     deps: deps]
  end

  # Configuration for the OTP application
  #
  # Type "mix help compile.app" for more information
  def application do
    [
      applications: [
        :comeonin,
        :cowboy,
        :ecto,
        :ectograph,
        :guardian,
        :logger,
        :plug,
        :plug_graphql,
        :postgrex,
      ],
      mod: {
        KeyMaps,
        []
      },
    ]
  end

  # Dependencies can be Hex packages:
  #
  #   {:mydep, "~> 0.3.0"}
  #
  # Or git/path repositories:
  #
  #   {:mydep, git: "https://github.com/elixir-lang/mydep.git", tag: "0.1.0"}
  #
  # Type "mix help deps" for more examples and options
  defp deps do
    [
      { :corsica, "~> 0.4.1" },
      { :comeonin, "~> 2.3.0" },
      { :cowboy, "~> 1.0.4" },
      { :ecto, "~> 1.1.5" },
      { :ectograph, "~> 0.0.3" },
      { :graphql, "~> 0.2.0" },
      { :guardian, "~> 0.10.1" },
      { :plug, "~> 1.1.2" },
      { :plug_graphql, "~> 0.2.0" },
      { :postgrex, "~> 0.11.1" },
    ]
  end
end
