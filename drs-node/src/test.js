import WS from './transport/websocket'

async function start() {
	const ws = new WS({
		token: 'eW2q1S7noJzFwfLwapnQpY9bZP3B4ELGlFuwDZ3f',
		device: String(Math.random()),
	})
	ws.router = () => 'localhost:12000'
	setInterval(async () => {
		const result = await ws.send({
			action: 'mutation',
			body: {
				op: {
					'echo.foo': {
						$merge: {
							message: 'hello',
						}
					}
				}
			}
		})
		console.log(result)
	}, 1000)
}

start()
