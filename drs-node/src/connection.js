import UUID from './uuid'
import { Error, Exception } from './error'
import Processor from './processor'

function timeout(ms) {
	return new Promise(resolve => {
		setTimeout(resolve, ms)
	})
}

export default class Connection {
	constructor(protocol) {
		this.cache = {}
		this._protocol = protocol
		this._pending = {}
		this._processor = new Processor()
		this.on = this._processor.on.bind(this._processor)
		this._interval = setInterval(() => this._ping(), 60 * 60 * 1000)
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
		const result = await new Promise(async resolve => {
			this._pending[cmd.key] = resolve
			await this.fire(cmd)
		})
		const wrapped = { ...result, key: cmd.key } 
		if (result.action === 'drs.error')
			throw new Error(wrapped)
		if (result.action === 'drs.exception')
			throw new Exception(wrapped)
		return wrapped
	}

	async fire(cmd) {
		if (!cmd.key)
			cmd.key = UUID.ascending()
		while (!this._closed) {
			try {
				this._raw.send(this._protocol.encode(cmd))
				break
			} catch (ex) {
				//
			}
			await timeout(1000)
		}
	}

	static async dial(transport, protocol, host) {
		const raw = await transport.dial(host)
		return new Connection(raw, protocol)
	}

	dial(transport, host) {
		if (this._closed)
			return
		return transport.dial(host)
			.then(async raw => {
				this._raw = raw
				this.read().then(async () => {
					await timeout(1000)
					await this.dial(transport, host)
				})
			})
			.catch(async ex => {
				console.log(ex)
				await timeout(1000)
				await this.dial(transport, host)
			})
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
			this._raw.on('close', resolve)
		})
	}

	close() {
		this._closed = true
		clearInterval(this._interval)
		Object.keys(this._pending).map(key => {
			this._pending[key].resolve({
				key,
				action: 'drs.exception',
				body: {
					message: 'Closing connection'
				}
			})
			delete this._pending[key]
		})
		if (this._raw) {
			this._raw.close()
		}
		return true
	}
}
