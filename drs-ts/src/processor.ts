import Connection from './connection'
import Command from './command'

declare type CommandHandler = (Connection, Command, Object) => Object


class Processor {
	private _hooks: Map<string, Array<CommandHandler>>
	constructor() {
	}

	public async on(action: string, ...handlers: Array<CommandHandler>) {
		this._hooks[action] = handlers
	}

	protected async process(cmd: Command) {
		console.log(cmd)
	}
}

export default Processor
