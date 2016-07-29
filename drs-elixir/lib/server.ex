defmodule Server do

  def start(_type, _args) do

    dispatch = :cowboy_router.compile([

      { :_,

        [
          {"/", Server.Default, []},
          {"/socket", Server.Socket, []},
      ]}
    ])
    { :ok, _ } = :cowboy.start_http(:http,
                                    100,
                                   [{:port, 12000}],
                                   [{ :env, [{:dispatch, dispatch}]}]
                                   )

    #device = File.open!("dbg.log", [:write])
    #Application.put_env(:dbg, :device, device)
    #IEx.configure([colors: [enabled: false]])
    #Dbg.reset
    #CowboyElixirExample.Supervisor.start_link([])
  end
end

defmodule Server.Default do

	def init(_type, req, []) do
		{:ok, req, :no_state}
	end

	def handle(request, state) do
		{ :ok, reply } = :cowboy_req.reply(
			200,
			[ {"content-type", "text/html"} ],
			"ok",
			request
			)
		{:ok, reply, state}
	end

	def terminate(reason, request, state) do
	:ok
	end

end

defmodule Server.Socket do
	@behaviour :cowboy_websocket_handler

	def init({_tcp, _http}, _req, _opts) do
		{:upgrade, :protocol, :cowboy_websocket}
	end

	def websocket_init(_TransportName, req, _opts) do
		IO.puts("New Connection")
		{:ok, req, :undefined_state }
	end

	# Required callback.  Put any essential clean-up here.
	def websocket_terminate(_reason, _req, _state) do
		IO.puts("Connection closed")
		:ok
	end

	def websocket_handle({:text, content}, req, state) do
		content
		# |> :zlib.gunzip
		|> Poison.decode!(as: %{})
		|> Server.Processor.process
		{:ok, req, state}
	end

	def websocket_info({:response, cmd}, req, state) do
		data = cmd
		|> Poison.encode!
		# |> :zlib.gzip
		{:reply, {:text, data}, req, state}
	end

	def websocket_info(_info, req, state) do
		{:ok, req, state}
	end
end

defmodule Server.Processor do
	def process(cmd) do
		caller = self()
		Task.start_link(fn ->
			handle(caller, cmd)
		end)
	end

	defp handle(caller, cmd) do
		%{"action" => action, "body" => body, "key" => key} = cmd
		body = on({action, body})
		result = %{
			key: key,
			action: "drs.response",
			body: body,
		}
		send(caller, {:response, result})
	end

	defp on({"raw.upload", body}) do
		%{"data" => data, "device" => device, "created" => created } = body
		data = data |> Base.decode64!
		path = "/raw/#{device}/#{created}"
		IO.puts("Uploading raw")
		Core.S3.put_object!("kinetic-device-data", path, data)
		IO.puts("Done uploading")
		true
	end

end
