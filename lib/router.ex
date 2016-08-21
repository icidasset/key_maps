defmodule KeyMaps.Router do
  use Plug.Router
  alias KeyMaps.{Auth}
  import KeyMaps.Utils
  require Logger

  # middleware
  plug Plug.Parsers,
    parsers: [:graphql, :urlencoded, :multipart, :json],
    pass: ["*/*"],
    json_decoder: Poison

  plug Corsica,
    allow_headers: ["accept", "authorization", "content-type", "origin"],
    origins: [
      ~r{^https?://localhost:\d+$},
      ~r{^https?://keymaps.surge.sh$},
    ]

  # => endpoint
  plug :match
  plug :dispatch


  def start_link do
    port = String.to_integer(System.get_env("PORT")) || 4000
    if Mix.env == :dev, do: Logger.info "Running `Key Maps` on port " <> Integer.to_string(port)
    { :ok, _ } = Plug.Adapters.Cowboy.http KeyMaps.Router, [], port: port
  end


  #
  # Authentication
  #
  post "/auth/start" do
    email = conn.params["email"]
    origin = get_req_header(conn, "origin") |> List.first

    case Auth.start(email, origin) do
      { :ok } -> render_empty(conn, 200)
      { :error, reason } -> render_error(conn, 422, reason)
    end

    render_empty(conn, 200)
  end


  post "/auth/exchange" do
    auth0_id_token = conn.params["auth0_id_token"]

    case Auth.exchange(auth0_id_token) do
      { :ok, user } -> render_token(conn, 200, user, auth0_id_token)
      { :error, reason } -> render_error(conn, 422, reason)
    end
  end


  get "/auth/validate" do
    token = conn.params["token"]

    case Auth.validate_token(token) do
      { :ok } -> render_empty(conn, 202)
      { :error }-> render_empty(conn, 403)
    end
  end


  def unauthenticated(conn, _) do
    render_error(conn, 403, "Forbidden")
  end


  #
  # Private API (GraphQL)
  #
  defmodule ApiPipeline do
    use Plug.Builder

    plug Guardian.Plug.VerifyHeader
    plug Guardian.Plug.LoadResource
    plug Guardian.Plug.EnsureAuthenticated, handler: KeyMaps.Router

    plug GraphQL.Plug,
      schema: { KeyMaps.GraphQL.Schema, :schema },
      root_value: &KeyMaps.GraphQL.Session.root_value/1
  end


  forward "/api",
    to: ApiPipeline


  #
  # Public API
  #
  forward "/public",
    to: KeyMaps.Public.Plug


  #
  # 404
  #
  match _, do: send_resp(conn, 404, "404")

end
