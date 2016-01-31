import WS from './transport/websocket'

async function start() {
	const ws = new WS()
	const result = await ws.send({
		action: 'ping',
	})
	console.log(result)
}

start()
