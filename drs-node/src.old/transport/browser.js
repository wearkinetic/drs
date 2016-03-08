import qs from 'querystring'
import Pipe from '../pipe'

export default class WS extends Pipe {
	constructor(query = {}, proto = 'ws') {
		super()
		this._query = query
		this._proto = proto
	}

	_connect(host) {
		return new Promise((resolve, reject) => {
			const ws = new global.WebSocket(`${this._proto}://${host}/socket?` + qs.stringify(this._query))
			ws.onerror = reject
			ws.onopen = () => resolve(session(ws))
		})
	}
}

function session(ws) {
	return {
		send(data) {
			return new Promise(resolve => {
				ws.send(data)
				resolve()
			})
		},
		on(action, cb) {
			if (action === 'data')
				ws.onmessage = msg => cb(msg.data)
			if (action === 'close')
				ws.onclose = () => cb()
		},
		close() {
			ws.close()
		}
	}
}
