defmodule KeyMaps.Router do
  use Plug.Router

  alias KeyMaps.{Models}

  import Comeonin.Bcrypt
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
    if Mix.env == :dev, do: Logger.info "Running `Key Maps` on port 4000"
    port = System.get_env("PORT") || 4000
    { :ok, _ } = Plug.Adapters.Cowboy.http KeyMaps.Router, [], port: port
  end


  #
  # Authentication
  #
  post "/sign-in" do
    login = conn.params["login"]
    password = conn.params["password"]

    user = Models.User.get_by_email(login) ||
           Models.User.get_by_username(login)

    if user && checkpw(password, user.password_hash),
      do: render_token(conn, 200, user),
    else: render_error(conn, 403, "The login and/or password were invalid")
  end


  post "/sign-up" do
    attr = %{
      email: conn.params["email"],
      password: conn.params["password"],
      username: conn.params["username"]
    }

    case Models.User.create(attr) do
      { :ok, user } -> render_token(conn, 201, user)
      { :error, changeset } -> render_error(conn, 400, get_error_from_changeset(changeset))
      nil -> render_error(conn, 403, "Sign-up is currently disabled")
    end
  end


  get "/validate-token" do
    token = conn.params["token"]

    case Guardian.decode_and_verify(token) do
      { :ok, _ } -> render_empty(conn, 202)
      { :error, _ } -> render_empty(conn, 403)
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
