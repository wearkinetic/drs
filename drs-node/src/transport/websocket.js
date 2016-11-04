import WS from 'uws'
import qs from 'querystring'
import EventEmitter from 'events'

export default class Websocket {
	constructor(query = {}, secure = false) {
		this._query = query
		this._secure = secure
	}

	dial(host) {
		return new Promise((resolve, reject) => {
			const ws = new WS(`${this._secure ? 'wss' : 'ws'}://${host}/socket?` + qs.stringify(this._query))
			ws.once('error', reject)
			ws.on('open', () => resolve(new Raw(ws)))
		})
	}
}

class Raw extends EventEmitter {
	constructor(ws) {
		super()
		this._ws = ws
		this._ws.on('error', e => this.emit('error', e))
		this._ws.on('message', async data => this.emit('data', data))
		this._ws.on('close', () => this.emit('close'))
	}
	send(data) {
		this._ws.send(data)
	}
	close() {
		this._ws.close()
	}
}
