defmodule KeyMaps.GraphQL.Plug do
  import Plug.Conn
  import KeyMaps.Utils

  alias Plug.Conn
  alias GraphQL.Plug.Endpoint

  @behaviour Plug


  def init(opts), do: Endpoint.init(opts)


  def call(%Conn{method: m} = conn, opts) when m in ["GET", "POST"] do
    handle_call(conn, Endpoint.extract_arguments(conn, opts))
  end


  def call(%Conn{method: _} = conn, _) do
    render_error(conn, 405, "GraphQL only supports GET and POST requests.")
  end


  def handle_call(conn, %{query: nil}), do: render_error(conn, 400, "Must provide query string")
  def handle_call(conn, args), do: execute(conn, args)


  defp execute(conn, args) do
    try do
      Plug.Conn.put_resp_content_type(conn, "application/json")

      res = GraphQL.execute(
        args.schema,
        args.query,
        args.root_value,
        args.variables,
        args.operation_name
      )

      case res do
        {:ok, data} ->
          case Poison.encode(data) do
            {:ok, json}      -> send_resp(conn, 200, json)
            {:error, errors} -> send_resp(conn, 500, errors)
          end
        {:error, errors} ->
          case Poison.encode(errors) do
            {:ok, json}      -> send_resp(conn, 400, json)
            {:error, errors} -> send_resp(conn, 400, errors)
          end
      end

    rescue
      err in GraphQL.CustomError ->
        msg = if Map.has_key?(err, :message),
          do: err.message,
          else: err.__struct__.message(err)

        render_error(conn, 400, msg)

    end
  end

end
