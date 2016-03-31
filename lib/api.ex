defmodule KeyMaps.Api do
  use Plug.Router

  alias KeyMaps.{ Models }

  import Comeonin.Bcrypt
  import KeyMaps.Utils

  # middleware
  plug Plug.Parsers,
    parsers: [:graphql, :urlencoded, :multipart, :json],
    pass: ["*/*"],
    json_decoder: Poison

  plug Guardian.Plug.VerifyHeader
  plug Guardian.Plug.LoadResource

  # => endpoint
  plug :match
  plug :dispatch


  def start_link do
    { :ok, _ } = Plug.Adapters.Cowboy.http KeyMaps.Api, [], port: 8080
  end


  #
  # GraphQL API
  #
  defmodule ApiPipeline do

    def init(_) do
      opts = [
        schema: { KeyMaps.GraphQL.Schema, :schema },
        root_value: &KeyMaps.GraphQL.Session.root_value/1
      ]

      graphiql = GraphQL.Plug.GraphiQL.init(opts)
      endpoint = KeyMaps.GraphQL.Plug.init(opts)

      opts = Keyword.merge(graphiql, endpoint)
      opts = Enum.dedup(opts)
      opts
    end

    def call(conn, opts) do
      if GraphQL.Plug.GraphiQL.use_graphiql?(conn, opts),
        do: GraphQL.Plug.GraphiQL.call(conn, opts),
        else: KeyMaps.GraphQL.Plug.call(conn, opts)
    end

  end


  forward "/api",
    to: ApiPipeline


  #
  # Authentication
  #
  post "/sign-in" do
    email = conn.params["email"]
    password = conn.params["password"]
    user = Models.User.get_by_email(email)

    if user && checkpw(password, user.password_hash) do
      render_token(conn, user)
    else
      render_error(conn, 403, "Email and/or password were invalid")
    end
  end


  post "/sign-up" do
    if conn.params["email"] == nil ||
       conn.params["password"] == nil do
      render_error(conn, 400, "Need a email and a password")
    else
      create_user(conn)
    end
  end


  def create_user(conn) do
    attr = %{
      email: conn.params["email"],
      password: conn.params["password"]
    }

    case Models.User.create(attr) do
      { :ok, user } -> render_token(conn, user)
      { :error, changeset } -> render_error(conn, 400, get_error_from_changeset(changeset))
    end
  end


  def unauthenticated(conn, _) do
     render_error(conn, 403, "Forbidden")
  end


  #
  # Public
  #
  forward "/public",
    to: KeyMaps.Public.Plug

end
