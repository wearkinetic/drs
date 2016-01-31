import Websocket from 'ws'
import WS from './transport/websocket'

async function start() {
	const ws = new WS()
	setInterval(async () => {
		const result = await ws.send({
			action: 'ping',
		})
		console.log(result)
	}, 1000)
}

start()
