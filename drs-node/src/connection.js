import UUID from './uuid'
import { Error, Exception } from './error'
import Processor from './processor'

export default class Connection {
	constructor(raw, protocol) {
		this.cache = {}
		this._raw = raw
		this._protocol = protocol
		this._pending = {}
		this._processor = new Processor()
		this.on = this._processor.on.bind(this._processor)

		this._interval = setInterval(() => this._ping(), 1000)
		this._ping()
	}

	async _ping() {
		const result = await this.send({
			action: 'drs.ping'
		})
		this._time = {
			server: result,
			local: Date.now(),
		}
	}

	now() {
		const { server, local } = this._time
		return server + (Date.now() - local)
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

	read() {
		return new Promise(resolve => {
			this._raw.on('data', async data => {
				const cmd = this._protocol.decode(data)
				if (cmd.action === 'drs.response' || cmd.action === 'drs.error' || cmd.action === 'drs.exception') {
					const waiting = this._pending[cmd.key]
					if (waiting) {
						waiting(cmd)
						delete this._pending[cmd.key]
						return
					}
					return
				}
				try {
					const result = await this._processor.process(cmd)
					this._processor.respond(cmd, this, result)
				} catch (ex) {
					this._processor.respond(cmd, this, ex)
				}
			})

			this._raw.on('close', () => {
				clearInterval(this._interval)
				resolve()
			})
		})
	}

	close() {
		return this._raw.close()
	}
}
