import WS from './transport/websocket'

async function start() {
	const ws = new WS({
		// token: 'eW2q1S7noJzFwfLwapnQpY9bZP3B4ELGlFuwDZ3f',
		token: 'djkhaled',
		device: String(Math.random()),
	})
	ws.on('jarvis.event', async cmd => {
		console.log(cmd)
	})
	ws.router = () => 'localhost:12000'
	console.log('Connected')
	await ws.send({
		action: 'jarvis.listen',
		body: {
			kind: 'convo.hello',
		},
	})
	await ws.send({
		action: 'jarvis.event',
		body: {
			kind: 'convo.hello',
			context: {
				sender: 'test'
			}
		}
	})
}

start()
