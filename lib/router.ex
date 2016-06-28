defmodule KeyMaps.Router do
  use Plug.Router

  alias KeyMaps.{Models}

  import Comeonin.Bcrypt
  import KeyMaps.Utils

  # middleware
  plug Plug.Parsers,
    parsers: [:graphql, :urlencoded, :multipart, :json],
    pass: ["*/*"],
    json_decoder: Poison

  # => endpoint
  plug :match
  plug :dispatch


  def start_link do
    { :ok, _ } = Plug.Adapters.Cowboy.http KeyMaps.Router, [], port: 8080
  end


  #
  # Authentication
  #
  post "/sign-in" do
    accessor = if conn.params["email"] && String.length(conn.params["email"]) > 0,
      do: "email",
    else: "username"

    accessor_value = conn.params[accessor]
    password = conn.params["password"]

    user = if accessor === "email",
      do: Models.User.get_by_email(accessor_value),
    else: Models.User.get_by_username(accessor_value)

    if user && checkpw(password, user.password_hash) do
      render_token(conn, 200, user)
    else
      accessor_label = String.capitalize(accessor)
      render_error(conn, 403, "#{accessor_label} and/or password were invalid")
    end
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

end
