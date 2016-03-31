defmodule KeyMaps.Public.Plug do
  import Plug.Conn
  import KeyMaps.Utils

  alias Plug.Conn
  alias GraphQL.Plug.Endpoint

  @behaviour Plug

  def init(_) do
    []
  end


  def call(%Conn{method: m} = conn, _) when m in ["GET"] do
    # IO.inspect(conn)

    render_data(conn, 200, %{})
  end

end
