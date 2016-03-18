/// <reference path="../../typings/ws/ws.d.ts"/>
import WebSocket = require('ws')
import { Raw, Connection } from '../';

export default class Websocket extends Connection  {
	public query = {}
	public secure = false
	protected async connect(host: string): Promise<Raw> {
		const params = Object.keys(this.query).map(key => {
			return {
				key,
				value: this.query[key]
			}
		}).reduce((collect, obj) => {
			return collect += obj.key + '=' + obj.value + '&'
		}, '')
		const url = `${this.secure ? 'wss' : 'ws'}://${host}/socket?` + params
		return new Promise<Raw>((resolve, reject) => {
			const ws = new WebSocket(url, err => {
				if (err)
					reject(err)
			})
			ws.on('error', reject)
			ws.on('open', () => {
				const result = new WebsocketRaw(ws)
				resolve(result)
			})
		})
	}
}

class WebsocketRaw extends Raw {
	_ws
	constructor(ws) {
		super()
		this._ws = ws
		ws.on('message', data => {
			this.onData(data)
		})

		ws.on('close', () => {
			this.onClose()
		})
	}
	async close() {
		this._ws.close()
	}
}
