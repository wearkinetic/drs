import Connection from './connection';

declare type CommandHandler = (Connection, Command, Object) => Object

class Command {
	key: string
	action: string
	body: Object
}

class Processor {
	private _hooks: Map<string, Array<CommandHandler>>
	constructor() {
	}

	public async on(action: string, ...handlers: Array<CommandHandler>) {
		this._hooks[action]
	}

	protected async process(cmd: Command) {
	}
}

export default Processor
