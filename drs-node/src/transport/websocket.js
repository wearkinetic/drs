import WS from 'ws'
import qs from 'querystring'
import Pipe from '../pipe'

export default class Websocket extends Pipe {
	constructor(query = {}) {
		super()
		this._query = query
	}

	_connect(host) {
		return new Promise((resolve, reject) => {
			const ws = new WS(`ws://${host}:12000/socket?` + qs.stringify(this._query))
			ws.once('error', reject)
			ws.on('open', () => resolve(session(ws)))
		})
	}
}

function session(ws) {
	return {
		send(data) {
			return new Promise((resolve, reject) => {
				ws.send(data, error => {
					if (error)
						reject(error)
					resolve()
				})
			})
		},
		on(action, cb) {
			if (action === 'data')
				ws.on('message', data => cb(data))
			if (action === 'close')
				ws.on('close', () => cb())
		}
	}
}
