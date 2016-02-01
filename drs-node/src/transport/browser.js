import qs from 'querystring'
import Pipe from '../pipe'

export default class WS extends Pipe {
	constructor(query = {}) {
		super()
		this._query = query
	}

	_connect(host) {
		return new Promise((resolve, reject) => {
			const ws = new global.WebSocket(`ws://${host}/socket?` + qs.stringify(this._query))
			ws.onerror = reject
			ws.onopen = () => resolve(session(ws))
		})
	}
}

function session(ws) {
	return {
		send(data) {
			return new Promise((resolve, reject) => {
				ws.send(data)
				resolve()
			})
		},
		on(action, cb) {
			if (action === 'data')
				ws.onmessage = data => cb(data)
			if (action === 'close')
				ws.onclose = data => cb()
		},
		close() {
			ws.close()
		}
	}
}
