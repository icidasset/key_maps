defmodule KeyMaps.Auth do
  import KeyMaps.Utils

  def get_auth0_domain do
    System.get_env("AUTH0_DOMAIN")
    |> String.trim_trailing("/")
  end


  def generate_token(user, auth0_id_token) do
    claims = %{ auth0_id_token: auth0_id_token, user: %{ id: user.id ,email: user.email }}
    Guardian.encode_and_sign(user, :long_lived_token, claims)
  end


  #
  # Start
  #
  def start(email, origin) do
    client_id = System.get_env("AUTH0_CLIENT_ID")
    auth0_start(email, client_id, origin)
  end


  #
  # Exchange
  #
  def exchange(auth0_id_token) do
    a = if is_nil(auth0_id_token), do: { :error, "Invalid Auth0 token" }

    if a,
      do: a,
    else: start_exchange_flow(auth0_id_token)
  end


  defp start_exchange_flow(auth0_id_token) do
    auth0_tokeninfo(auth0_id_token)
    |> create_user
  end


  defp create_user(results) do
    case results do
      { :ok, info } ->
        email = info["email"]
        user = KeyMaps.Models.User.get_by_email(email)
        disallow_signup = (System.get_env("ENABLE_SIGN_UP") != "1")

        cond do
          user -> { :ok, user }
          disallow_signup -> { :error, "Sign-up is current disabled" }

          true ->
            case KeyMaps.Models.User.create(%{ email: email }) do
              { :ok, user } -> { :ok, user }
              { :error, changeset } -> { :error, get_error_from_changeset(changeset) }
            end
        end

      _ -> results
    end
  end


  #
  # Validate
  #
  def validate_token(token) do
    case Guardian.decode_and_verify(token) do
      { :ok, claims } ->
        auth0_id_token = claims["auth0_id_token"]

        case auth0_tokeninfo(auth0_id_token) do
          { :ok, _ } -> { :ok }
          { :error, _ } -> { :error }
        end

      { :error, _ } -> { :error }
    end
  end


  #
  # Auth0 requests
  #
  defp auth0_start(email, client_id, origin) do
    url = "https://" <> get_auth0_domain <> "/passwordless/start"
    body = %{
      client_id: client_id,
      connection: "email",
      email: email,
      send: "link",
      authParams: %{ redirect_uri: origin }
    }
    body = Poison.encode!(body)
    headers = [{ "Content-Type", "application/json" }]

    case HTTPoison.post(url, body, headers) do
      { :ok, response } ->
        x = response.status_code
        cond do
          x < 300 -> { :ok }
          x > 300 -> { :error, "Auth0 authentication failed" }
        end
      { :error, _ } ->
        { :error, "Could not authenticate via Auth0" }
    end
  end


  defp auth0_tokeninfo(auth0_id_token) do
    url = "https://" <> get_auth0_domain <> "/tokeninfo"
    body = Poison.encode!(%{ id_token: auth0_id_token })
    headers = [{ "Content-Type", "application/json" }]

    case HTTPoison.post(url, body, headers) do
      { :ok, response } ->
        x = response.status_code
        cond do
          x < 300 -> { :ok, Poison.decode!(response.body) }
          x > 300 -> { :error, "Invalid Auth0 token" }
        end
      { :error, _ } ->
        { :error, "Could not validate Auth0 token" }
    end
  end

end
