import JSON from './protocol/json'
import UUID from './uuid'
import Error from './error'
import Connection from './connection'
import EventEmitter from 'events'

const ACTIONS = {
	response: 'drs.response',
	error: 'drs.error',
	exception: 'drs.exception',
}

function timeout(time) {
	return new Promise(resolve => {
		setTimeout(resolve, time)
	})
}

let count = 0

export default class Pipe {
	constructor() {
		count++
		this.protocol = JSON
		this.router = () => {
			return 'localhost'
		}
		this._handlers = {}
		this._connections = {}
		this._pending = {}
		this._queue = []
		this.closing = false
		this.events = new EventEmitter()
	}

	on(action, ...handlers) {
		this._handlers[action] = handlers
		if (this._on)
			this._on(action, handlers)
	}

	async _loop() {
		if (this._working)
			return
		const cmd = this._queue.shift()
		if (!cmd) {
			return
		}
		this._working = true
		while (true) {
			try {
				const conn = await this._route(cmd.action)
				await conn.send(cmd)
			} catch (ex) {
				if (this.closing)
					break
				console.log('out', ex)
				await timeout(1000)
				continue
			}
			break
		}
		this._working = false
		this._loop()
	}

	async send(cmd) {
		if (!cmd.key)
			cmd.key = UUID.ascending()
		while (true) {
			if (this.closing)
				throw new Error('forcing close')
			const prom = new Promise(resolve => {
				this._pending[cmd.key] = {
					resolve,
				}
			})
			this._queue.push(cmd)
			this._loop()
			const response = await prom
			if (response.action === ACTIONS.exception) {
				await timeout(1000)
				continue
			}
			if (response.action === ACTIONS.error)
				throw response.body
			return response.body
		}
	}

	async _route(action) {
		const host = this.router(action)
		return await this._dial(host)
	}

	async _dial(host) {
		let conn = this._connections[host]
		if (conn)
			return conn
		const rw = await this._connect(host)
		if (this.closing) {
			rw.close()
			throw new Error('closing')
		}
		conn = new Connection(rw, this.protocol)
		this._connections[host] = conn
		conn.raw.on('close', () => delete this._connections[host])
		this._handle(conn)
		this.events.emit('connect', conn, host)
		return conn
	}

	_handle(conn) {
		conn.raw.on('data', async function(data) {
			try {
				const cmd = this.protocol.decode(data)
				await this._process(conn, cmd)
			} catch (ex) {
				console.log('in', ex)
			}
		}.bind(this))
	}

	async _process(conn, cmd) {
		if (cmd.action === ACTIONS.response || cmd.action === ACTIONS.error || cmd.action === ACTIONS.exception) {
			const match = this._pending[cmd.key]
			if (!match)
				return
			match.resolve(cmd)
			delete this._pending[cmd.key]
			return
		}

		const handlers = this._handlers[cmd.action]
		if (!handlers)
			return
		const ctx = {}
		let result
		try {
			for (let h of handlers) {
				result = await h(cmd, conn, ctx)
			}
		} catch (ex) {
			const response = {
				key: cmd.key,
				action: ACTIONS.exception,
				body: {
					message: String(ex),
				}
			}
			if (ex instanceof Error) {
				response.action = ACTIONS.error
				response.body = ex
			}
			conn.send(response)
			return
		}
		conn.send({
			key: cmd.key,
			action: ACTIONS.response,
			body: result,
		})
	}

	close() {
		count--
		this.closing = true
		this._queue = []
		this.events.removeAllLiseners('connect')
		Object.values(this._pending).map(key => {
			this._pending[key].resolve({
				key,
				action: ACTIONS.error,
				body: {
					message: 'forcing close',
				},
			})
		})
		this._pending = {}
		Object.values(this._connections).map(conn => conn.raw.close())
		return
	}

}

setInterval(() => console.log('total drs: ' + count), 10000)
