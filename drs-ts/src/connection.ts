import { Raw } from './'
import Command from './command'
import Processor from './processor'

function sleep(ms: Number): Promise<void> {
	return new Promise<void>(resolve => {
		setTimeout(resolve, ms)
	})
}

abstract class Connection extends Processor {
	private _raw: Raw
	constructor() {
		super()
	}

	public async dial(host: string, reconnect: boolean): Promise<void> {
		while (true) {
			try {
				const raw = await this.connect(host)
				await this.handle(raw)
			} catch (ex) {
				console.log(ex)
			}
			if (!reconnect)
				return
			console.log('Reconnecting')
			await sleep(1000)
		}
	}

	protected abstract connect(host: string): Promise<Raw>

	private async handle(raw: Raw) {
		this._raw = raw
		return new Promise(resolve => {
			this._raw.onData = data => {
				const command: Command = JSON.parse(data)
				this.process(command)
			}
			this._raw.onClose = () => {
				resolve()
			}
		})
	}
}

export default Connection
