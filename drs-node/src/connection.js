export default class Connection {
	constructor(rw, protocol) {
		this.raw = rw
		this._protocol = protocol
		this._cache = {}
	}

	set(key, value) {
		this._cache[key] = value
	}

	get(key) {
		return this._cache[key]
	}

	async send(cmd) {
		const data = await this._protocol.encode(cmd)
		return this.raw.send(data)
	}
}
