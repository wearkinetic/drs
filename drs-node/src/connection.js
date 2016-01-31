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

	encode(cmd) {
		const data = this._protocol.encode(cmd)
		this.raw.send(data)
	}
}
