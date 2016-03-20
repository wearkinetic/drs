import Connection from './connection'
import Command from './command'

declare type CommandHandler = (Connection, Command, Object) => Object

class Processor {
	private _handlers: Map<string, Array<CommandHandler>>
	private _pending: Map<string, (cmd: Command) => void>
	public redirect: Processor

	constructor() {
		this._pending = new Map<string, (cmd: Command) => void>()
		this._handlers = new Map<string, Array<CommandHandler>>()
	}

	public async on(action: string, ...handlers: Array<CommandHandler>) {
		this._handlers[action] = handlers
	}

	public async wait(cmd: Command, cb: () => void): Promise<Object> {
		const result  = <Command>(await new Promise(async (resolve) => {
			await cb()
			this._pending[cmd.key] = resolve
		}))
		if (result.action !== 'drs.response')
			throw result.body
		return result.body
	}

	protected async process(cmd: Command, conn: Connection) {
		if (this.redirect)
			return this.redirect.process(cmd, conn)
		if (cmd.action === 'drs.response' || cmd.action === 'drs.error' || cmd.action === 'drs.exception') {
			const match = this._pending[cmd.key]
			if (match)
				match(cmd)
			return
		}

		const handlers = this._handlers[cmd.action]
		if (!handlers)
			return
		let context = {}
		const output: Command = {
			key: cmd.key,
			action: 'drs.response',
			body: {},
		}
		try {
			for (let h of handlers) {
				cmd.body = await h(cmd, conn, context)
			}
		} catch (ex) {
			cmd.action = 'drs.error'
			cmd.body = ex
		}
		conn.fire(cmd)
	}

	protected async clear() {
		for (let key in this._pending) {
			this._pending[key]({
				action: 'drs.exception',
				body: {
					message: 'Connection closed'
				}
			})
		}
	}
}

export default Processor
