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

export default class Pipe {
	constructor() {
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
				if(this.closing)
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
				response.Action = ACTIONS.error
				response.Body = ex
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
		this.closing = true
		Object.keys(this._connections).map(key => {
			this._connections[key].raw.close()
		})
		Object.keys(this._pending).forEach(key => {
			this._pending[key].resolve({
				key,
				action: ACTIONS.error,
				body: {
					message: 'forcing close',
				},
			})
		})
		this._queue = []
		return 
	}

}
