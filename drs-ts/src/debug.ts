import Websocket from './transports/websocket'
import Command from './command'

async function start(): Promise<void> {
	const ws = new Websocket()
	ws.dial('localhost:12000', true)
	for (let i = 0; i < 1000; i++) {
		try {
			const result = await ws.request({
				key: undefined,
				action: 'drs.ping',
				body: {},
			})
			console.log(result)
		} catch (ex) {
			console.log(ex)
		}
	}
	ws.close()
}

start()
