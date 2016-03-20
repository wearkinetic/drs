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
	private _closed: boolean
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
			if (!reconnect || this._closed)
				return
			console.log('Reconnecting')
			await sleep(1000)
		}
	}

	public async request(cmd: Command): Promise<Object> {
		if (!cmd.key)
			cmd.key = String(Math.random())
		const result = await this.wait(cmd, () => {
			this.fire(cmd)
		})
		return result
	}

	public async fire(cmd: Command) {
		if (!cmd.key)
			cmd.key = String(Math.random())
		while(true) {
			try {
				this._raw.send(JSON.stringify(cmd))
				return
			} catch (ex) {
			}
			await sleep(1000)
		}
	}

	protected abstract connect(host: string): Promise<Raw>

	private async handle(raw: Raw) {
		this._raw = raw
		return new Promise(resolve => {
			this._raw.onData = data => {
				const command: Command = JSON.parse(data)
				this.process(command, this)
			}
			this._raw.onClose = () => {
				try {
					this.clear()
				} catch (ex) {
					console.log(ex)
				}
				resolve()
			}
		})
	}

	public close(): void {
		this._closed = true
		this._raw.close()
	}
}

export default Connection
