import WS from './transport/websocket'
import JSON from './protocol/json'
import Connection from './connection'

const transport = new WS({
	token: 'ZGhD0ptkn0XTk7DbphrvwXPwuOPI2pEB5qIBWqjq',
	device: String(Math.random()),
})

async function start() {
	try {
		const conn = await Connection.dial(transport, JSON, 'localhost:12000')
		conn.on('test', () => {
			console.log('ok')
		})
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
