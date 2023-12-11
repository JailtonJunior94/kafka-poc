defmodule MyProducer.MixProject do
  use Mix.Project

  def project do
    [
      app: :my_producer,
      version: "0.1.0",
      elixir: "~> 1.15",
      start_permanent: Mix.env() == :prod,
      deps: deps()
    ]
  end

  # Run "mix help compile.app" to learn about applications.
  def application do
    [
      extra_applications: [:logger, :kaffe],
      mod: {MyProducer.Application, []}
    ]
  end

  # Run "mix help deps" to learn about dependencies.
  defp deps do
    [
      # {:dep_from_hexpm, "~> 0.3.0"},
      # {:dep_from_git, git: "https://github.com/elixir-lang/my_dep.git", tag: "0.1.0"}
      {:kaffe, "~> 1.24"},
      {:protox, "~> 1.6"},
      {:jason, "~> 1.2"},
      {:confluent_schema_registry, "~> 0.1.1"}
    ]
  end
end
