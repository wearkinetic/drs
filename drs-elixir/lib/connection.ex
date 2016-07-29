defmodule DRS.Connection do
	use GenServer

	def start_link(url) do
		GenServer.start(__MODULE__, url)
	end

	def init(url) do
		{:ok, socket} = DRS.Socket.start_link(url, self())
		{:ok, %{
			socket: socket,
			pending: %{},
		}}
	end

	def request(process, action, body) do
		{:ok, result} = GenServer.call(process, {:request, %{action: action, body: body}}, :infinity)
		result
	end

	def handle_call({:request, cmd}, caller, state) do
		key = UUID.uuid4()
		cmd = Map.put(cmd, "key", key)
		pending = Map.put(state.pending, key, caller)
		data = Poison.encode!(cmd)
		:websocket_client.cast(state.socket, {:text, data})
		{:noreply, %{state | pending: pending}}
	end

	def process(connection, data) do
		GenServer.cast(connection, {:process, data})
	end

	def handle_cast({:process, data}, state) do
		Task.start_link(fn ->
			cmd = %{
				"action" => action,
				"body" => body,
				"key" => key,
			} = data
			# |> :zlib.gunzip
			|> Poison.decode!(as: %{})
			case state.pending do
				%{^key => process} ->
					GenServer.reply(process, {:ok, body})
				_ ->
					IO.inspect(cmd)
			end
		end)
		{:noreply, state}
	end
end

defmodule DRS.Socket do
	@behaviour :websocket_client_handler

	def start_link(url, drs) do
		:websocket_client.start_link(url, __MODULE__, drs)
	end

	def init(drs, conn) do
		{:ok, drs}
	end

	def websocket_handle({:text, message}, _conn, drs) do
		DRS.Connection.process(drs, message)
		{:ok, drs}
	end


end
