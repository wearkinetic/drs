var WS = require('ws')
var qs = require('querystring')
var EventEmitter = require('events')

function Websocket(query, secure) {
	this._query = query || {}
	this._secure = secure || false
}

Websocket.prototype.dial = function(host) {
	return new Promise(resolve => {
		var ws = new WS(`${this._secure ? 'wss' : 'ws'}://${host}/socket?` + qs.stringify(this._query))
		ws.once('error', () => resolve())
		ws.on('open', () => resolve(new Raw(ws)))
	})
}

function Raw(ws) {
	this._ws = ws
	this._ws.on('message', data => this.emit('data', data))
	this._ws.on('close', () => this.emit('close'))
}

Raw.prototype.send = function(data) {
	console.log(data)
	return new Promise((resolve, reject) => {
		this._ws.send(data, err => {
			if (err) {
				reject()
				return
			}
			resolve()
		})
	})
}

Raw.prototype.close = function() {
	this._ws.close()
}

Raw.prototype.__proto__ = EventEmitter.prototype

module.exports = Websocket
