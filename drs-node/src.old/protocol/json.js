export default {
	encode(body) {
		return JSON.stringify(body)
	},
	decode(data) {
		return JSON.parse(data)
	}
}
