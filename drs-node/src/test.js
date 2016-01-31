import WS from './transport/websocket'

async function start() {
	const ws = new WS({
		token: 'eW2q1S7noJzFwfLwapnQpY9bZP3B4ELGlFuwDZ3f',
		device: String(Math.random()),
	})
	ws.router = () => 'drs.virginia.inboxtheapp.com'
	ws.on('mutation', cmd => console.log(cmd))
	setInterval(async () => {
		const result = await ws.send({
			action: 'drs.ping',
		})
		console.log(result)
	}, 1000)
}

start()
