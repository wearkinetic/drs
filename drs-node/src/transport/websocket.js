import WS from 'ws'
import qs from 'querystring'
import Pipe from '../pipe'

export default class Websocket extends Pipe {
	constructor(query = {}, proto = 'ws') {
		super()
		this._query = query
		this._proto = proto
	}

	_connect(host) {
		return new Promise((resolve, reject) => {
			const ws = new WS(`${this._proto}://${host}/socket?` + qs.stringify(this._query))
			ws.once('error', reject)
			ws.on('open', () => resolve(session(ws)))
		})
	}
}

function session(ws) {
	return {
		send(data) {
			return new Promise((resolve, reject) => {
				ws.send(data, err => {
					if (err)
						reject(err)
					resolve()
				})
			})
		},
		on(action, cb) {
			if (action === 'data')
				return ws.on('message', cb)
			ws.on('close', cb)
		},
		close() {
			ws.close(1000)
		}
	}
}
