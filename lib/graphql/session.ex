defmodule KeyMaps.GraphQL.Session do

  def root_value(conn) do
    user = Guardian.Plug.current_resource(conn)

    if user,
      do: %{ user_id: user.id },
      else: %{}
  end

end
