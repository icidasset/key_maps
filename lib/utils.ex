defmodule KeyMaps.Utils do

  def render_data(conn, status, map) do
    render_json(conn, status, %{ data: map })
  end


  def render_error(conn, status, message) do
    render_json(conn, status, %{ errors: [%{ message: message }] })
  end


  def render_token(conn, status, user) do
    conn = Guardian.Plug.api_sign_in(conn, user)
    token = Guardian.Plug.current_token(conn)

    data = %{
      token: token,
    }

    render_data(conn, status, data)
  end


  def get_error_from_changeset(changeset) do
    cond do
      # errors
      length(changeset.errors) > 0 ->
        field = Keyword.keys(changeset.errors) |> List.first |> Atom.to_string
        field = String.replace(field, "_", " ")
        fieldError = Keyword.values(changeset.errors) |> List.first

        if is_tuple(fieldError) do
          fieldError = String.replace(
            elem(fieldError, 0),
            "%{count}",
            to_string(elem(fieldError, 1)[:count])
          )
        end

        String.capitalize(field) <> " " <> fieldError

      # constraints
      length(changeset.constraints) > 0 ->
        changeset.constraints[0].message

      # fallback
      true ->
        model = changeset.model.__struct__ |> Module.split |> List.last
        action = to_string(changeset.action)

        "Could not " <> action <> " " <> model
    end
  end


  def extract_other_arguments(internal) do
    inn = internal.field_asts |> List.first

    other_attr = if Map.has_key?(inn, :arguments) do
      Enum.reduce inn.arguments, %{}, fn(arg, acc) ->
        k = String.to_atom(arg.name.value)
        v = arg.value.value
        Map.put(acc, k, v)
      end
    end

    other_attr
  end


  #
  # Private
  #
  defp render_json(conn, status, map) do
    content = Poison.encode!(map)

    Plug.Conn.put_resp_content_type(conn, "application/json")
    Plug.Conn.send_resp(conn, status, content)
  end

end
