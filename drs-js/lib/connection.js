'use strict'
var WS = require('./transports/websocket')
var Processor = require('./processor')

function Connection() {
	this._closed = false
	this._processor = new Processor()
}


Connection.prototype.dial = function (transport, host, reconnect) {
	return new Promise(resolve => {
		var loop = () => {
			if (this._closed)
				return
			transport.dial(host).then(raw => {
				if (!raw) {
					setTimeout(loop, 1000)
					return
				}
				this.raw = raw
				this._handle().then(() => {
					if (reconnect) {
						setTimeout(loop, 1000)
						return
					}
				})
			})
		}
		loop()
	})
}

Connection.prototype._handle = function (cb) {
	return new Promise(resolve => {
		this.raw.on('close', () => resolve())
		this.raw.on('data', data => {
			console.log(data)
		})
	})
}

Connection.prototype.close = function() {
	this._closed = true
	this.raw.close()
}

Connection.prototype.fire = function(cmd) {
	if (cmd.key == '')
		cmd.key = String(Math.random())
	return new Promise((resolve, reject) => {
		var loop = () => {
			if (this._closed) {
				reject()
				return
			}
			try {
				console.log('Sending')
				this.raw.send(JSON.stringify(cmd)).catch(() => setTimeout(loop, 1000))
				resolve()
			} catch (ex) {
				setTimeout(loop, 1000)
			}
		}
		loop()
	})
}

module.exports = Connection

var transport = new WS({
	token: 'ZGhD0ptkn0XTk7DbphrvwXPwuOPI2pEB5qIBWqjq'
})
var conn = new Connection()
conn.dial(transport, 'localhost:12000', true).catch(console.log)
conn.fire({
	action: 'delta.sync',
	body: {
		offset: '',
	},
})
