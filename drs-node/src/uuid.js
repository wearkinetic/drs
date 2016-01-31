function build(seed) {
	let result = seed.toString(16)
	result = new Array(17 - result.length).join('0') + result

	for (var i = 0; i < 16; i++) {
		result += Math.floor(Math.random() * 10).toString(16)
	}
	return result
}

const MAX = Math.pow(2, 42)

class UUID {
	ascending() {
		return build(new Date().getTime())
	}

	descending() {
		return build(MAX - new Date().getTime())
	}
}

export default new UUID()
