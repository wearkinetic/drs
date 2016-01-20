var pipe = new Pipe()

// Server
pipe.on('echo', (conn, cmd) => {
	conn.send({
		action: 'echo',
		body: cmd.body,
	})
})
pipe.listen()

// Client
pipe.connect(action => 'localhost:3002')
pipe.send({
	action: 'echo',
	body: 'hello',
})
