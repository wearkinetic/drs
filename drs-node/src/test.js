import WS from './transport/websocket'
import JSON from './protocol/json'
import Connection from './connection'

const transport = new WS({
	token: 'djkhaled',
	device: String(Math.random()),
})

async function start() {
	try {
		const conn = await Connection.dial(transport, JSON, 'localhost:12000')
		conn.read()
		const result = await conn.send({
			action: 'delta.mutation',
			body: {
				op: {
					'node.test': {
						$merge: {
							nice: Date.now(),
						}
					}
				}
			}
		})
		console.log(result)
	} catch (ex) {
		console.log(ex)
		console.log(ex.stack)
	}
}

start()
