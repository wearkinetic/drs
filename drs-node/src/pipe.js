import JSON from './protocol/json'
import UUID from './uuid'
import Error from './error'
import Connection from './connection'

const ACTIONS = {
	response: 'drs.response',
	error: 'drs.error',
	exception: 'drs.exception',
}

export default class Pipe {
	constructor(transport) {
		this.protocol = JSON
		this.router = action => {
			return 'localhost'
		}

		this._handlers = {}
		this._connections = {}
		this._pending = {}
	}

	on(action, ...handlers) {
		this._handlers[action] = handlers
		if (this._on)
			this._on(action, handlers)
	}

	async send(cmd) {
		if (!cmd.key)
			cmd.key = UUID.ascending()
		while (true) {
			let conn
			try {
				conn = await this._route(cmd.action)
			} catch(ex) {
				continue
			}
			const prom = new Promise((resolve, reject) => {
				this._pending[cmd.key] = {
					resolve,
					reject,
				}
			})
			conn.encode(cmd)
			const response = await prom
			if (response.action === ACTIONS.error)
				throw response.body
			if (response.action === ACTIONS.exception)
				continue
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
		conn.raw.on('close', () => delete(this._connections[host]))
		this._handle(conn)
		return conn
	}

	_handle(conn) {
		conn.raw.on('data', async function(data) {
			try {
				const cmd = this.protocol.decode(data)
				await this._process(conn, cmd)
			} catch(ex) {
				console.log(ex)
			}
		}.bind(this))
	}

	async _process(conn, cmd) {
		if (cmd.action === ACTIONS.response || cmd.action === ACTIONS.error || cmd.action === ACTIONS.exception) {
			const match = this._pending[cmd.key]
			match.resolve(cmd)
		}

		const handlers = this._handlers[cmd.action]
		const ctx = {}
		let result = undefined
		try {
			for (let h of handlers) {
				result = await h(cmd, conn, ctx)
			}
		} catch(ex) {
			const response = {
				key: cmd.key,
				action: ACTIONS.exception,
				body: {
					nessage: String(ex),
				}
			}
			if (ex instanceof Error) {
				response.Action = ACTIONS.error
				response.Body = ex
			}
			conn.encode(response)
			return
		}
		conn.encode({
			key: cmd.key,
			action: ACTIONS.response,
			body: result,
		})
	}

}
