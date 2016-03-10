export default class Processor {
	constructor() {
		this._handlers = {}
		this.redirect = undefined
	}

	on(action, ...handlers) {
		this._handlers[action] = handlers
	}

	process(cmd, conn) {
		if (this.redirect)
			return this.redirect.process(cmd, conn)
		const handlers = this._handlers[cmd.action]
		if (!handlers)
			return undefined
		return this._trigger(cmd, conn, ...handlers)
	}

	respond(cmd, conn, body) {
		const response = {
			key: cmd.key,
			body,
			action: 'drs.response'
		}
		conn.fire(response)
	}

	async _trigger(cmd, conn, ...handlers) {
		const ctx = {}
		let result
		for (let h of handlers) {
			result = await h(cmd, conn, ctx)
		}
		return result
	}
}
