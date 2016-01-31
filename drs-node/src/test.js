import Websocket from 'ws'
import WS from './transport/websocket'

async function start() {
	const ws = new WS()
	ws.router = () => 'delta.virginia.inbox.svc.tutum.io'
	setInterval(async () => {
		const result = await ws.send({
			action: 'drs.ping',
		})
		console.log(result)
	}, 1000)
}

start()
