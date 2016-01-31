import WS from 'ws'
import qs from 'querystring'
import Pipe from '../pipe'

export default class Websocket extends Pipe {
	constructor(query = {}) {
		super()
		this._query = query
	}

	_connect(host) {
		const ws = new WS(`ws://${host}:12000/socket?` + qs.stringify(this._query))
		return new Promise(resolve => {
			ws.on('open', () => {
				resolve({
					send(data) {
						ws.send(data)
					},
					on(action, cb) {
						if (action === 'data')
							ws.on('message', data => cb(data))
						if (action === 'close')
							ws.on('close', () => cb())
					}
				})
			})
		})
	}
}
