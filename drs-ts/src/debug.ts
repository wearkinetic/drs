import Websocket from './transports/websocket'

async function start(): Promise<void> {
	const ws = new Websocket()
	console.log(ws)
	await ws.dial('localhost:12000', true)
}

start()
