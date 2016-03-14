import WS from './transport/websocket'

async function start() {
	const ws = new WS({
		// token: 'eW2q1S7noJzFwfLwapnQpY9bZP3B4ELGlFuwDZ3f',
		token: 'U7Vwcc2kA3XX4H9LLR5mxkhWxZY60RwKIaPFW96P',
		device: String(Math.random()),
	})
	ws.on('mutation', async cmd => {
		console.log(cmd)
	})
	ws.router = () => 'drs.virginia.inboxtheapp.com'
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
