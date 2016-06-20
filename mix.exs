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
        :slugger,
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
      { :corsica, "~> 0.4.2" },
      { :comeonin, "~> 2.4.0" },
      { :cowboy, "~> 1.0.4" },
      { :ecto, "~> 2.0.0-rc.6" },
      { :ectograph, "~> 0.1.1" },
      { :graphql, "~> 0.3.1" },
      { :guardian, "~> 0.12.0" },
      { :plug, "~> 1.1.6" },
      { :plug_graphql, "~> 0.3.1" },
      { :postgrex, "~> 0.11.2" },
      { :slugger, "~> 0.1.0" }
    ]
  end
end
