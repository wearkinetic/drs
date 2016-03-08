import UUID from './uuid'
import { Error, Exception } from './error'

export class Connection {
	constructor(raw, protocol) {
		this.cache = {}
		this._raw = raw
		this._protocol = protocol
		this._pending = {}
	}

	async send(cmd) {
		if (!cmd.key)
			cmd.key = UUID.ascending()
		const result = await new Promise(resolve => {
			this._pending[cmd.key] = resolve
			this.fire(cmd)
		})
		if (result.action === 'drs.error')
			throw new Error(result.body)
		if (result.action === 'drs.exception')
			throw new Exception(result.body)
		return result.body
	}

	fire(cmd) {
		if (!cmd.key)
			cmd.key = UUID.ascending()
		this._raw.send(this._protocol.encode(cmd))
	}

	static async dial(transport, protocol, host) {
		const raw = await transport.dial(host)
		return new Connection(raw, protocol)
	}
}
